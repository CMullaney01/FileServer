import React from "react";
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import DesignPage from "./pages/DesignPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage  />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/background" element={<DesignPage />} />
      </Routes>
    </Router>
  );
}

export default App;
