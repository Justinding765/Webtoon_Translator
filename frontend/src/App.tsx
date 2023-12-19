// src/App.tsx
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HomePage from './pages/hompage';
import Translation from './pages/Translation';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/translation" element={<Translation />} />
        {/* Add other routes here */}
      </Routes>
    </Router>
  );
}

export default App;
