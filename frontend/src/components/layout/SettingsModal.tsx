import React, { useState, useEffect } from 'react';
import { X, Server, Brain, Plus, Trash2, CheckCircle2, AlertCircle, Loader2 } from 'lucide-react';
import { GetConfig, SetProvider, AddMcpServer, RemoveMcpServer, ListMcpServers } from '../../wailsjs/go/main/App';
import { config } from '../../wailsjs/go/models';

interface SettingsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const SettingsModal: React.FC<SettingsModalProps> = ({ isOpen, onClose }) => {
  const [activeTab, setActiveTab] = useState<'provider' | 'mcp'>('provider');
  const [appConfig, setAppConfig] = useState<config.Config | null>(null);
  const [loading, setLoading] = useState(false);
  
  // Provider state
  const [providerName, setProviderName] = useState('ollama');
  const [providerUrl, setProviderUrl] = useState('http://localhost:11434');
  const [providerModel, setProviderModel] = useState('llama3');

  // MCP state
  const [mcpName, setMcpName] = useState('');
  const [mcpType, setMcpType] = useState<'stdio' | 'sse'>('stdio');
  const [mcpCommand, setMcpCommand] = useState('');
  const [mcpArgs, setMcpArgs] = useState('');
  const [mcpUrl, setMcpUrl] = useState('');

  useEffect(() => {
    if (isOpen) {
      loadConfig();
    }
  }, [isOpen]);

  const loadConfig = async () => {
    try {
      const cfg = await GetConfig();
      setAppConfig(cfg);
      setProviderName(cfg.provider.name);
      setProviderUrl(cfg.provider.baseUrl);
      setProviderModel(cfg.provider.model);
    } catch (err) {
      console.error('Failed to load config:', err);
    }
  };

  const handleSaveProvider = async () => {
    setLoading(true);
    try {
      await SetProvider(providerName, providerUrl, providerModel);
      await loadConfig();
    } catch (err) {
      console.error('Failed to save provider:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddMcp = async () => {
    setLoading(true);
    try {
      const args = mcpArgs.split(',').map(a => a.trim()).filter(a => a !== '');
      await AddMcpServer({
        name: mcpName,
        type: mcpType,
        command: mcpCommand,
        args: args,
        url: mcpUrl,
      } as any);
      setMcpName('');
      setMcpCommand('');
      setMcpArgs('');
      setMcpUrl('');
      await loadConfig();
    } catch (err) {
      console.error('Failed to add MCP server:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleRemoveMcp = async (name: string) => {
    try {
      await RemoveMcpServer(name);
      await loadConfig();
    } catch (err) {
      console.error('Failed to remove MCP server:', err);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div className="bg-zinc-900 border border-zinc-800 rounded-xl shadow-2xl w-full max-w-2xl max-h-[80vh] flex flex-col overflow-hidden">
        {/* Header */}
        <div className="p-4 border-b border-zinc-800 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-zinc-100 flex items-center gap-2">
            <Server size={20} className="text-blue-500" />
            Application Settings
          </h2>
          <button onClick={onClose} className="text-zinc-500 hover:text-zinc-200 p-1">
            <X size={20} />
          </button>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-zinc-800 bg-zinc-950/50">
          <button
            onClick={() => setActiveTab('provider')}
            className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
              activeTab === 'provider' ? 'border-blue-500 text-blue-500 bg-zinc-900' : 'border-transparent text-zinc-500 hover:text-zinc-300'
            }`}
          >
            <Brain size={16} className="inline mr-2" />
            LLM Provider
          </button>
          <button
            onClick={() => setActiveTab('mcp')}
            className={`px-6 py-3 text-sm font-medium transition-colors border-b-2 ${
              activeTab === 'mcp' ? 'border-blue-500 text-blue-500 bg-zinc-900' : 'border-transparent text-zinc-500 hover:text-zinc-300'
            }`}
          >
            <Server size={16} className="inline mr-2" />
            MCP Servers
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {activeTab === 'provider' && (
            <div className="space-y-4 animate-in fade-in duration-200">
              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase tracking-wider mb-2">Provider Type</label>
                <select 
                  value={providerName}
                  onChange={(e) => setProviderName(e.target.value)}
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="ollama">Ollama (Local)</option>
                  <option value="openai">OpenAI (API)</option>
                  <option value="anthropic">Anthropic (API)</option>
                </select>
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase tracking-wider mb-2">Base URL</label>
                <input 
                  type="text"
                  value={providerUrl}
                  onChange={(e) => setProviderUrl(e.target.value)}
                  placeholder="http://localhost:11434"
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:ring-2 focus:ring-blue-500 placeholder:text-zinc-600"
                />
              </div>

              <div>
                <label className="block text-xs font-bold text-zinc-500 uppercase tracking-wider mb-2">Model Name</label>
                <input 
                  type="text"
                  value={providerModel}
                  onChange={(e) => setProviderModel(e.target.value)}
                  placeholder="llama3"
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-md px-3 py-2 text-sm text-zinc-100 focus:outline-none focus:ring-2 focus:ring-blue-500 placeholder:text-zinc-600"
                />
              </div>

              <div className="pt-4">
                <button 
                  onClick={handleSaveProvider}
                  disabled={loading}
                  className="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-blue-800 text-white font-medium py-2 px-4 rounded-md transition-colors flex items-center justify-center gap-2"
                >
                  {loading && <Loader2 size={16} className="animate-spin" />}
                  Save Configuration
                </button>
              </div>
            </div>
          )}

          {activeTab === 'mcp' && (
            <div className="space-y-6 animate-in fade-in duration-200">
              {/* Add New Server */}
              <div className="bg-zinc-800/50 border border-zinc-700 rounded-lg p-4 space-y-4">
                <h3 className="text-sm font-semibold text-zinc-100 flex items-center gap-2">
                  <Plus size={16} />
                  Add New MCP Server
                </h3>
                
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-[10px] font-bold text-zinc-500 uppercase mb-1">Server Name</label>
                    <input 
                      type="text"
                      value={mcpName}
                      onChange={(e) => setMcpName(e.target.value)}
                      placeholder="e.g. filesystem"
                      className="w-full bg-zinc-900 border border-zinc-700 rounded px-2 py-1.5 text-xs text-zinc-100 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label className="block text-[10px] font-bold text-zinc-500 uppercase mb-1">Transport Type</label>
                    <select 
                      value={mcpType}
                      onChange={(e) => setMcpType(e.target.value as any)}
                      className="w-full bg-zinc-900 border border-zinc-700 rounded px-2 py-1.5 text-xs text-zinc-100 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    >
                      <option value="stdio">stdio (Local)</option>
                      <option value="sse">sse (Remote)</option>
                    </select>
                  </div>
                </div>

                {mcpType === 'stdio' ? (
                  <div className="space-y-3">
                    <div>
                      <label className="block text-[10px] font-bold text-zinc-500 uppercase mb-1">Command</label>
                      <input 
                        type="text"
                        value={mcpCommand}
                        onChange={(e) => setMcpCommand(e.target.value)}
                        placeholder="npx"
                        className="w-full bg-zinc-900 border border-zinc-700 rounded px-2 py-1.5 text-xs text-zinc-100 focus:outline-none focus:ring-1 focus:ring-blue-500"
                      />
                    </div>
                    <div>
                      <label className="block text-[10px] font-bold text-zinc-500 uppercase mb-1">Arguments (comma separated)</label>
                      <input 
                        type="text"
                        value={mcpArgs}
                        onChange={(e) => setMcpArgs(e.target.value)}
                        placeholder="-y, @modelcontextprotocol/server-everything"
                        className="w-full bg-zinc-900 border border-zinc-700 rounded px-2 py-1.5 text-xs text-zinc-100 focus:outline-none focus:ring-1 focus:ring-blue-500"
                      />
                    </div>
                  </div>
                ) : (
                  <div>
                    <label className="block text-[10px] font-bold text-zinc-500 uppercase mb-1">SSE URL</label>
                    <input 
                      type="text"
                      value={mcpUrl}
                      onChange={(e) => setMcpUrl(e.target.value)}
                      placeholder="http://localhost:8080/sse"
                      className="w-full bg-zinc-900 border border-zinc-700 rounded px-2 py-1.5 text-xs text-zinc-100 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                  </div>
                )}

                <button 
                  onClick={handleAddMcp}
                  disabled={loading || !mcpName}
                  className="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-zinc-700 disabled:text-zinc-500 text-white font-medium py-1.5 rounded text-xs transition-colors flex items-center justify-center gap-2"
                >
                  Connect Server
                </button>
              </div>

              {/* Active Servers List */}
              <div className="space-y-3">
                <h3 className="text-xs font-bold text-zinc-500 uppercase tracking-widest">Connected Servers</h3>
                <div className="space-y-2">
                  {appConfig?.mcpServers?.map(server => (
                    <div key={server.name} className="flex items-center justify-between p-3 bg-zinc-800 rounded-md border border-zinc-700">
                      <div className="flex items-center gap-3">
                        <div className="w-2 h-2 rounded-full bg-green-500" />
                        <div>
                          <div className="text-sm font-medium text-zinc-100">{server.name}</div>
                          <div className="text-[10px] text-zinc-500 font-mono">
                            {server.type} {server.type === 'stdio' ? `(${server.command})` : `(${server.url})`}
                          </div>
                        </div>
                      </div>
                      <button 
                        onClick={() => handleRemoveMcp(server.name)}
                        className="text-zinc-500 hover:text-red-400 p-1.5 transition-colors"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  ))}
                  {(!appConfig?.mcpServers || appConfig.mcpServers.length === 0) && (
                    <div className="text-center py-6 border-2 border-dashed border-zinc-800 rounded-md">
                      <p className="text-xs text-zinc-500">No MCP servers connected.</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
