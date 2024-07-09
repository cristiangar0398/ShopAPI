import './styles/App.css';
import NavigationBar from './components/Navar'
import Login from './components/Validation'

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <NavigationBar/>
        <Login/>
      </header>
    </div>
  );
}

export default App;
