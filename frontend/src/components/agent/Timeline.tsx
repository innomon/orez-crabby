import React, {useMemo} from 'react';
import {Step, clusterSteps} from '../../types/agent';
import {StepCard} from './StepCard';

interface TimelineProps {
    steps: Step[];
}

export const Timeline: React.FC<TimelineProps> = ({steps}) => {
    const groups = useMemo(() => clusterSteps(steps), [steps]);

    return (
        <div className="flex-1 overflow-y-auto p-6 space-y-8 max-w-4xl mx-auto w-full">
            {groups.map((group) => (
                <div key={group.id} className="space-y-4">
                    {group.mode === 'exploration' ? (
                        <div className="border-l-2 border-zinc-800 ml-4 pl-8 space-y-4">
                            <div className="text-[10px] font-bold text-zinc-600 uppercase tracking-widest -ml-4 bg-zinc-950 inline-block px-2">
                                Exploration Sequence ({group.steps.length} steps)
                            </div>
                            {group.steps.map((step) => (
                                <StepCard key={step.id} step={step} />
                            ))}
                        </div>
                    ) : (
                        group.steps.map((step) => (
                            <StepCard key={step.id} step={step} />
                        ))
                    )}
                </div>
            ))}
        </div>
    );
};
