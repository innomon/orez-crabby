import {Sidebar} from './components/layout/Sidebar';
import {MainView} from './components/layout/MainView';

function App() {
    return (
        <div id="App" className="flex h-screen w-screen overflow-hidden">
            <Sidebar />
            <MainView />
        </div>
    )
}

export default App
