import React, {useState} from 'react';
import {Terminal, ChevronRight, ChevronDown, CheckCircle2, AlertCircle, Loader2} from 'lucide-react';
import {Step} from '../../types/agent';

interface StepCardProps {
    step: Step;
}

export const StepCard: React.FC<StepCardProps> = ({step}) => {
    const [isExpanded, setIsExpanded] = useState(false);

    if (step.kind === 'thought') {
        return (
            <div className="flex gap-4 group">
                <div className="w-8 h-8 rounded-full bg-blue-500/10 flex items-center justify-center border border-blue-500/30 mt-1 shrink-0">
                    <div className="w-1.5 h-1.5 rounded-full bg-blue-400" />
                </div>
                <div className="flex-1 space-y-1">
                    <div className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">Agent Thought</div>
                    <div className="text-zinc-300 text-sm leading-relaxed">{step.content}</div>
                </div>
            </div>
        );
    }

    if (step.kind === 'tool_call') {
        const isWaiting = step.content === 'WAITING_FOR_APPROVAL';
        const isExecuting = step.content === 'EXECUTING';
        const isSuccess = !isWaiting && !isExecuting && !step.content.startsWith('Error:');

        return (
            <div className="flex gap-4 group">
                <div className="w-8 h-8 rounded-full bg-zinc-800 flex items-center justify-center border border-zinc-700 mt-1 shrink-0">
                    <Terminal size={14} className="text-zinc-400" />
                </div>
                <div className="flex-1 space-y-2 bg-zinc-900/40 border border-zinc-800 rounded-lg overflow-hidden">
                    <button 
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="w-full flex items-center justify-between p-3 hover:bg-zinc-800/50 transition-colors"
                    >
                        <div className="flex items-center gap-3">
                            <span className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">Tool: {step.tool_name}</span>
                            <code className="text-[11px] text-zinc-400 bg-zinc-800 px-1.5 py-0.5 rounded">{step.tool_input}</code>
                        </div>
                        <div className="flex items-center gap-3">
                            {isWaiting && (
                                <span className="text-[10px] bg-amber-500/10 text-amber-500 px-1.5 py-0.5 rounded border border-amber-500/20 flex items-center gap-1">
                                    Pending Approval
                                </span>
                            )}
                            {isExecuting && (
                                <span className="text-[10px] bg-blue-500/10 text-blue-500 px-1.5 py-0.5 rounded border border-blue-500/20 flex items-center gap-1">
                                    <Loader2 size={10} className="animate-spin" /> Executing
                                </span>
                            )}
                            {isSuccess && (
                                <CheckCircle2 size={14} className="text-emerald-500" />
                            )}
                            {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                        </div>
                    </button>
                    
                    {isExpanded && (
                        <div className="p-4 pt-0 border-t border-zinc-800 bg-zinc-950/50 font-mono text-[11px] text-zinc-400 overflow-x-auto whitespace-pre">
                            {step.content}
                        </div>
                    )}
                </div>
            </div>
        );
    }

    if (step.kind === 'response') {
        return (
            <div className="flex gap-4 group">
                 <div className="w-8 h-8 rounded-full bg-emerald-500/10 flex items-center justify-center border border-emerald-500/30 mt-1 shrink-0">
                    <div className="w-1.5 h-1.5 rounded-full bg-emerald-400" />
                </div>
                <div className="flex-1 space-y-1">
                    <div className="text-[10px] font-bold text-zinc-500 uppercase tracking-widest">Assistant</div>
                    <div className="text-zinc-100 text-[13px] leading-relaxed prose prose-invert max-w-none">
                        {step.content}
                    </div>
                </div>
            </div>
        );
    }

    return null;
};
