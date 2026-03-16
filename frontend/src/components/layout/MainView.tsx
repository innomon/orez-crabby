import React, {useState, useEffect, useRef} from 'react';
import {Send, Terminal} from 'lucide-react';
import {EventsOn} from '../../../wailsjs/runtime';
import {RunAgent} from '../../../wailsjs/go/main/App';
import {Step} from '../../types/agent';
import {Timeline} from '../agent/Timeline';

export const MainView: React.FC = () => {
  const [input, setInput] = useState('');
  const [steps, setSteps] = useState<Step[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Listen for agent steps from Go
    const unsubscribe = EventsOn('agent:step', (step: Step) => {
      setSteps(prev => {
        // If it's a tool update, replace the existing step
        const existingIdx = prev.findIndex(s => s.id === step.id && step.id !== "");
        if (existingIdx !== -1) {
           const next = [...prev];
           next[existingIdx] = step;
           return next;
        }
        return [...prev, step];
      });
    });

    return () => unsubscribe();
  }, []);

  useEffect(() => {
    // Auto scroll to bottom
    if (scrollRef.current) {
        scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [steps]);

  const handleSubmit = async () => {
    if (!input.trim() || isStreaming) return;
    
    setIsStreaming(true);
    setInput('');
    
    // For now, we use a hardcoded session ID
    await RunAgent("default-session", input);
    setIsStreaming(false);
  };

  return (
    <div className="flex-1 flex flex-col bg-zinc-950 text-zinc-100 h-screen">
      {/* Header */}
      <div className="h-14 border-b border-zinc-900 flex items-center px-6 justify-between">
        <div className="flex items-center gap-4">
          <span className="font-semibold text-zinc-200">Session: Main</span>
          <span className="text-[10px] bg-zinc-900 text-zinc-500 px-2 py-1 rounded-full border border-zinc-800 font-bold tracking-widest uppercase">
            Local Engine: Llama3
          </span>
        </div>
      </div>

      {/* Timeline Area */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto">
        <Timeline steps={steps} />
      </div>

      {/* Composer */}
      <div className="p-6 border-t border-zinc-900 bg-zinc-950/50 backdrop-blur">
        <div className="max-w-4xl mx-auto relative">
          <textarea 
            rows={1}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleSubmit();
                }
            }}
            placeholder="Describe what you want to achieve..."
            className="w-full bg-zinc-900 border border-zinc-800 rounded-xl py-4 px-5 pr-12 focus:outline-none focus:ring-1 focus:ring-zinc-600 transition-all resize-none text-zinc-200 text-sm"
          />
          <button 
            onClick={handleSubmit}
            disabled={isStreaming}
            className="absolute right-3 bottom-3 p-1.5 bg-zinc-100 text-zinc-950 rounded-lg hover:bg-zinc-200 transition-colors disabled:opacity-50"
          >
            <Send size={18} />
          </button>
        </div>
        <div className="max-w-4xl mx-auto mt-3 flex items-center justify-between text-[10px] text-zinc-600 px-1 font-bold uppercase tracking-wider">
          <div className="flex gap-4">
            <span>Enter to send</span>
            <span>Shift + Enter for new line</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className={`w-1.5 h-1.5 rounded-full ${isStreaming ? 'bg-blue-500 animate-pulse' : 'bg-green-500'}`} />
            {isStreaming ? 'Agent is working...' : 'Ready'}
          </div>
        </div>
      </div>
    </div>
  );
};
