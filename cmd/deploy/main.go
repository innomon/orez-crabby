package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	godotenv.Load(".env.local")

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	region := os.Getenv("GOOGLE_CLOUD_REGION")
	serviceName := os.Getenv("SERVICE_NAME")
	bucketName := os.Getenv("GOOGLE_CLOUD_BUILD_BUCKET")

	if serviceName == "" {
		serviceName = "orez-crabby"
	}
	if region == "" {
		region = "us-central1"
	}
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}
	if bucketName == "" {
		bucketName = fmt.Sprintf("%s-build-source", projectID)
	}

	ctx := context.Background()

	imageName := fmt.Sprintf("%s-docker.pkg.dev/%s/cloud-run-source-deploy/%s:latest", region, projectID, serviceName)

	// 1. Prepare and Upload Source
	sourceGCSPath, err := uploadSource(ctx, projectID, bucketName)
	if err != nil {
		log.Fatalf("Failed to upload source: %v", err)
	}

	// 2. Build Image
	err = buildImage(ctx, projectID, sourceGCSPath, imageName)
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	// 3. Deploy to Cloud Run
	err = deployToCloudRun(ctx, projectID, region, serviceName, imageName)
	if err != nil {
		log.Fatalf("Deployment failed: %v", err)
	}

	fmt.Printf("\nSuccess! Service deployed to %s\n", serviceName)
}

func uploadSource(ctx context.Context, projectID, bucketName string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	if _, err := bucket.Attrs(ctx); err != nil {
		fmt.Printf("Creating bucket %s...\n", bucketName)
		if err := bucket.Create(ctx, projectID, nil); err != nil {
			return "", err
		}
	}

	objectName := fmt.Sprintf("source-%d.tar.gz", time.Now().Unix())
	wc := bucket.Object(objectName).NewWriter(ctx)

	gw := gzip.NewWriter(wc)
	tw := tar.NewWriter(gw)

	// Files to ignore
	ignore := map[string]bool{
		".git":         true,
		"node_modules": true,
		"dist":         true,
		"bin":          true,
		"build":        true,
	}

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ignore[info.Name()] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = path

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})

	tw.Close()
	gw.Close()
	wc.Close()

	return fmt.Sprintf("gs://%s/%s", bucketName, objectName), nil
}

func buildImage(ctx context.Context, projectID, sourceGCSPath, imageName string) error {
	client, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket, object := parseGCSPath(sourceGCSPath)

	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: projectID,
		Build: &cloudbuildpb.Build{
			Source: &cloudbuildpb.Source{
				Source: &cloudbuildpb.Source_StorageSource{
					StorageSource: &cloudbuildpb.StorageSource{
						Bucket: bucket,
						Object: object,
					},
				},
			},
			Steps: []*cloudbuildpb.BuildStep{
				{
					Name: "gcr.io/cloud-builders/docker",
					Args: []string{"build", "-t", imageName, "."},
				},
				{
					Name: "gcr.io/cloud-builders/docker",
					Args: []string{"push", imageName},
				},
			},
			Images: []string{imageName},
		},
	}

	fmt.Printf("Starting Cloud Build for %s...\n", imageName)
	op, err := client.CreateBuild(ctx, req)
	if err != nil {
		return err
	}

	_, err = op.Wait(ctx)
	return err
}

func deployToCloudRun(ctx context.Context, projectID, region, serviceName, imageName string) error {
	client, err := run.NewServicesClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	servicePath := fmt.Sprintf("%s/services/%s", parent, serviceName)

	service := &runpb.Service{
		Template: &runpb.RevisionTemplate{
			Containers: []*runpb.Container{
				{
					Image: imageName,
					Ports: []*runpb.ContainerPort{
						{ContainerPort: 8080},
					},
				},
			},
		},
	}

	fmt.Printf("Deploying to Cloud Run: %s\n", serviceName)
	
	// Check if service exists
	_, err = client.GetService(ctx, &runpb.GetServiceRequest{Name: servicePath})
	var op *run.UpdateServiceOperation
	if err != nil {
		fmt.Println("Service does not exist, creating...")
		createOp, err := client.CreateService(ctx, &runpb.CreateServiceRequest{
			Parent:    parent,
			ServiceId: serviceName,
			Service:   service,
		})
		if err != nil {
			return err
		}
		_, err = createOp.Wait(ctx)
		return err
	}

	fmt.Println("Service exists, updating...")
	service.Name = servicePath
	op, err = client.UpdateService(ctx, &runpb.UpdateServiceRequest{
		Service: service,
	})
	if err != nil {
		return err
	}
	_, err = op.Wait(ctx)
	return err
}

func parseGCSPath(path string) (string, string) {
	s := strings.TrimPrefix(path, "gs://")
	parts := strings.SplitN(s, "/", 2)
	return parts[0], parts[1]
}
