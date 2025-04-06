import React from "react";
import { Header } from "./components/header/header";
import { UserProvider } from "./contexts/UserContext";
import './App.css';

const App: React.FC = () => {
  return (
      <UserProvider>
        <div className="App" style={{ padding: "0 140px" }}>
          <Header />
        </div>
      </UserProvider>
  );
};

export default App;
