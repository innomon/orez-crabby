export type StepKind = 'thought' | 'tool_call' | 'response' | 'error';

export interface Step {
    id: string;
    session_id: string;
    kind: StepKind;
    content: string;
    tool_name?: string;
    tool_input?: string;
    created_at: string;
}

export type StepGroupMode = 'exploration' | 'execution' | 'thought' | 'response';

export interface StepGroup {
    id: string;
    mode: StepGroupMode;
    steps: Step[];
}

const EXPLORATION_TOOLS = new Set(['read_file', 'glob', 'grep', 'ls', 'list_files']);

export function clusterSteps(steps: Step[]): StepGroup[] {
    const groups: StepGroup[] = [];
    let currentGroup: StepGroup | null = null;

    steps.forEach((step) => {
        let mode: StepGroupMode = 'thought';
        if (step.kind === 'tool_call') {
            mode = EXPLORATION_TOOLS.has(step.tool_name || '') ? 'exploration' : 'execution';
        } else if (step.kind === 'response') {
            mode = 'response';
        }

        if (currentGroup && currentGroup.mode === mode && mode === 'exploration') {
            currentGroup.steps.push(step);
        } else {
            currentGroup = {
                id: step.id || Math.random().toString(36).substr(2, 9),
                mode,
                steps: [step],
            };
            groups.push(currentGroup);
        }
    });

    return groups;
}
