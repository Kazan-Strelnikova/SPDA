import React from "react";
import { Header } from "./components/header/header";
import { UserProvider } from "./contexts/UserContext";
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { MainPage } from "./pages/main/Main";
import { LogInPage } from "./pages/LogIn";
import { SignUpPage } from "./pages/SignUp";
import { createTheme, ThemeProvider } from '@mui/material/styles';
import variables from './variables.module.scss';
import { CreateEventPage } from "./pages/create-event/CreateEvent";


const theme = createTheme({
  palette: {
    primary: {
      main: variables.primary,
    },
  },
  components: {
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            '&.Mui-focused fieldset': {
              borderColor: variables.primary,
            },
          },
          '& .MuiInputLabel-root.Mui-focused': {
            color: variables.primary,
          },
        },
      },
    },
  },
});

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <Router>
        <UserProvider>
          <div className="App" style={{ padding: "0 140px" }}>
            <Header />
            <Routes>
              <Route path="/" element={<MainPage />} />
              <Route path="/create" element={<CreateEventPage />} />
              <Route path="/login" element={<LogInPage />} />
              <Route path="/signup" element={<SignUpPage />} />
            </Routes>
          </div>
        </UserProvider>
          
      </Router>
      </ThemeProvider>
      
  );
};

export default App;
