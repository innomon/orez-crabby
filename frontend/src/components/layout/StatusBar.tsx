import React, { useState, useEffect } from 'react';
import { Server, Brain, Activity, Wifi, WifiOff } from 'lucide-react';
import { GetConfig } from '../../../wailsjs/go/main/App';
import { config } from '../../../wailsjs/go/models';

export const StatusBar: React.FC = () => {
  const [appConfig, setAppConfig] = useState<config.Config | null>(null);

  useEffect(() => {
    loadConfig();
    const interval = setInterval(loadConfig, 10000); // Refresh every 10s
    return () => clearInterval(interval);
  }, []);

  const loadConfig = async () => {
    try {
      const cfg = await GetConfig();
      setAppConfig(cfg);
    } catch (err) {
      console.error('Failed to load config in StatusBar:', err);
    }
  };

  const mcpCount = appConfig?.mcpServers?.length || 0;

  return (
    <div className="h-8 bg-zinc-950 border-t border-zinc-800 flex items-center justify-between px-4 text-[11px] text-zinc-500 overflow-hidden shrink-0">
      <div className="flex items-center gap-6">
        {/* Provider Status */}
        <div className="flex items-center gap-2">
          <Brain size={12} className={appConfig?.provider ? "text-blue-500" : "text-zinc-600"} />
          <span className="font-medium text-zinc-400">{appConfig?.provider?.name || 'No Provider'}</span>
          <span className="text-zinc-600 font-mono">{appConfig?.provider?.model}</span>
        </div>

        {/* MCP Status */}
        <div className="flex items-center gap-2 border-l border-zinc-800 pl-6">
          <Server size={12} className={mcpCount > 0 ? "text-green-500" : "text-zinc-600"} />
          <span className="font-medium text-zinc-400">{mcpCount} MCP Servers</span>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-1.5 bg-zinc-900 px-2 py-0.5 rounded border border-zinc-800">
          <Activity size={10} className="text-green-500 animate-pulse" />
          <span className="text-zinc-400">System Ready</span>
        </div>
        <div className="flex items-center gap-1">
          <Wifi size={12} className="text-blue-500" />
          <span className="text-zinc-500 uppercase tracking-tighter">Connected</span>
        </div>
      </div>
    </div>
  );
};
