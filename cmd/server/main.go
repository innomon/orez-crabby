package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"orez-crabby/internal/db"
	"orez-crabby/internal/models"
	"orez-crabby/pkg/agent"
	"orez-crabby/pkg/config"
)

type Server struct {
	agent         *agent.Agent
	clients       map[chan models.Step]bool
	clientsMu     sync.Mutex
	configManager *config.ConfigManager
}

func NewServer() *Server {
	cfgMgr, _ := config.NewConfigManager("/tmp/config.json")
	cfgMgr.Load()
	cfg := cfgMgr.Get()
	
	// Default to Ollama for now
	provider := agent.NewOllamaProvider(cfg.Provider.BaseURL, cfg.Provider.Model)
	if cfg.Provider.Name == "" {
		provider = agent.NewOllamaProvider("http://localhost:11434", "llama3")
	}

	return &Server{
		agent:         agent.NewAgent(provider),
		clients:       make(map[chan models.Step]bool),
		configManager: cfgMgr,
	}
}

func (s *Server) Start(port string) error {
	// Initialize DB
	if _, err := db.InitDB(); err != nil {
		log.Printf("Warning: Failed to initialize DB: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/RunAgent", s.handleRunAgent)
	mux.HandleFunc("/api/events", s.handleEvents)
	mux.HandleFunc("/api/GetConfig", s.handleGetConfig)
	
	// Serve frontend from ./frontend/dist
	fs := http.FileServer(http.Dir("./frontend/dist"))
	mux.Handle("/", fs)

	log.Printf("Server listening on port %s", port)
	return http.ListenAndServe(":"+port, mux)
}

func (s *Server) handleRunAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"sessionID"`
		Input     string `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := s.agent.Run(ctx, req.SessionID, req.Input, func(step models.Step) {
		s.broadcastStep(step)
	})

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "Success"})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan models.Step, 10)
	s.clientsMu.Lock()
	s.clients[ch] = true
	s.clientsMu.Unlock()

	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, ch)
		s.clientsMu.Unlock()
		close(ch)
	}()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case step := <-ch:
			data, _ := json.Marshal(step)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}

func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.configManager.Get())
}

func (s *Server) broadcastStep(step models.Step) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()
	for ch := range s.clients {
		ch <- step
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := NewServer()
	if err := server.Start(port); err != nil {
		log.Fatal(err)
	}
}
