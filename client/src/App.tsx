import React from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <div style={{backgroundColor: "green", padding: "30px", borderRadius: "16px"}}>
          <h1>
            <b><a href="https://rutube.ru/video/c6cc4d620b1d4338901770a44b3e82f4/" style={{color: "white"}}>Click me!</a></b>
          </h1>
        </div>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
