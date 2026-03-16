import { useState } from 'react';
import { Sidebar } from './components/layout/Sidebar';
import { MainView } from './components/layout/MainView';
import { StatusBar } from './components/layout/StatusBar';
import { SettingsModal } from './components/layout/SettingsModal';

function App() {
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  return (
    <div id="App" className="flex flex-col h-screen w-screen overflow-hidden bg-zinc-950 text-zinc-100">
      <div className="flex flex-1 overflow-hidden">
        <Sidebar onSettingsClick={() => setIsSettingsOpen(true)} />
        <MainView />
      </div>
      <StatusBar />
      <SettingsModal isOpen={isSettingsOpen} onClose={() => setIsSettingsOpen(false)} />
    </div>
  );
}

export default App;
