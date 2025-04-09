import React from "react";
import { Header } from "./components/header/header";
import { UserProvider } from "./contexts/UserContext";
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { MainPage } from "./pages/main/Main";
import { LogInPage } from "./pages/LogIn";
import { SignUpPage } from "./pages/SignUp";

const App: React.FC = () => {
  return (
    <Router>
      <UserProvider>
        <div className="App" style={{ padding: "0 140px" }}>
          <Header />
          <Routes>
            <Route path="/" element={<MainPage />} />
            <Route path="/login" element={<LogInPage />} />
            <Route path="/signup" element={<SignUpPage />} />
          </Routes>
        </div>
      </UserProvider>
        
    </Router>
      
  );
};

export default App;
