import React, {useState, useEffect} from 'react';
import {Plus, Settings, MessageSquare, Folder, FolderPlus, FileText} from 'lucide-react';
import {SelectWorkspace, GetWorkspaceFiles} from '../../wailsjs/go/main/App';
import {main} from '../../wailsjs/go/models';

export const Sidebar: React.FC = () => {
  const [workspace, setWorkspace] = useState<main.Workspace | null>(null);
  const [files, setFiles] = useState<any[]>([]);

  const handleSelectWorkspace = async () => {
    try {
      const ws = await SelectWorkspace();
      setWorkspace(ws);
      const fs = await GetWorkspaceFiles(ws.path);
      setFiles(fs);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <div className="w-64 bg-zinc-900 border-r border-zinc-800 flex flex-col text-zinc-300">
      <div className="p-4 border-b border-zinc-800 flex items-center justify-between">
        <h1 className="font-bold text-zinc-100 flex items-center gap-2">
          <Folder size={18} className="text-blue-500" />
          OpenWork-Go
        </h1>
      </div>

      <div className="p-4 space-y-3">
        {!workspace ? (
          <button 
            onClick={handleSelectWorkspace}
            className="w-full flex items-center justify-center gap-2 bg-zinc-100 text-zinc-900 py-2 px-4 rounded-md font-medium hover:bg-zinc-200 transition-colors text-sm"
          >
            <FolderPlus size={16} />
            Select Workspace
          </button>
        ) : (
          <button className="w-full flex items-center justify-center gap-2 bg-zinc-100 text-zinc-900 py-2 px-4 rounded-md font-medium hover:bg-zinc-200 transition-colors text-sm">
            <Plus size={16} />
            New Session
          </button>
        )}
      </div>

      <div className="flex-1 overflow-y-auto px-2 space-y-4">
        {workspace && (
           <div>
             <div className="text-[10px] font-bold text-zinc-600 uppercase tracking-widest px-3 py-1 flex items-center justify-between">
               Files
               <span className="text-zinc-700 lowercase font-normal">{workspace.name}</span>
             </div>
             <div className="mt-1 space-y-0.5">
               {files.map(file => (
                 <div key={file.path} className="flex items-center gap-2 px-3 py-1.5 hover:bg-zinc-800 rounded-md cursor-pointer group">
                   {file.is_dir ? <Folder size={14} className="text-zinc-500" /> : <FileText size={14} className="text-zinc-500" />}
                   <span className="text-xs truncate text-zinc-400 group-hover:text-zinc-200">{file.name}</span>
                 </div>
               ))}
             </div>
           </div>
        )}

        <div>
            <div className="text-[10px] font-bold text-zinc-600 uppercase tracking-widest px-3 py-1">
            Recent Sessions
            </div>
            <div className="mt-1 space-y-0.5">
            <button className="w-full text-left px-3 py-1.5 rounded-md hover:bg-zinc-800 flex items-center gap-3 group">
                <MessageSquare size={14} className="text-zinc-600 group-hover:text-zinc-400" />
                <span className="text-xs truncate">Setup Wails Project</span>
            </button>
            <button className="w-full text-left px-3 py-1.5 rounded-md hover:bg-zinc-800 flex items-center gap-3 group">
                <MessageSquare size={14} className="text-zinc-600 group-hover:text-zinc-400" />
                <span className="text-xs truncate">Debug SQLite logic</span>
            </button>
            </div>
        </div>
      </div>

      <div className="p-4 border-t border-zinc-800">
        <button className="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-zinc-800 text-zinc-400 hover:text-zinc-200 transition-colors text-sm">
          <Settings size={18} />
          Settings
        </button>
      </div>
    </div>
  );
};
