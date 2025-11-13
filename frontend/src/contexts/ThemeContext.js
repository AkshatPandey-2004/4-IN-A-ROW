import React, { createContext, useState, useContext, useEffect } from 'react';

const ThemeContext = createContext();

export const themes = {
  neonPurple: {
    name: 'Neon Purple',
    primary: '#9D4EDD',
    secondary: '#7B2CBF',
    accent: '#C77DFF',
    background: 'linear-gradient(135deg, #10002b 0%, #240046 50%, #3c096c 100%)',
    text: '#E0AAFF',
    textSecondary: '#C77DFF',
    board: 'linear-gradient(135deg, #5a189a, #3c096c)',
    cell: 'rgba(157, 78, 221, 0.2)',
    piece1: 'linear-gradient(135deg, #ff006e, #fb5607)',
    piece2: 'linear-gradient(135deg, #00f5ff, #00b4d8)',
    cardBg: 'rgba(16, 0, 43, 0.8)',
    shadow: 'rgba(157, 78, 221, 0.5)',
    glow: '0 0 30px rgba(157, 78, 221, 0.6)',
  },
  cyberpunk: {
    name: 'Cyberpunk',
    primary: '#FF0080',
    secondary: '#00FFFF',
    accent: '#FFFF00',
    background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #16213e 100%)',
    text: '#00FFFF',
    textSecondary: '#FF0080',
    board: 'linear-gradient(135deg, #1a1a2e, #0f3460)',
    cell: 'rgba(255, 0, 128, 0.2)',
    piece1: 'linear-gradient(135deg, #FF0080, #FF0048)',
    piece2: 'linear-gradient(135deg, #00FFFF, #00D4FF)',
    cardBg: 'rgba(26, 26, 46, 0.9)',
    shadow: 'rgba(255, 0, 128, 0.5)',
    glow: '0 0 30px rgba(0, 255, 255, 0.6)',
  },
  light: {
    name: 'Light',
    primary: '#667eea',
    secondary: '#764ba2',
    accent: '#f093fb',
    background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
    text: '#2d3748',
    textSecondary: '#4a5568',
    board: 'linear-gradient(135deg, #667eea, #764ba2)',
    cell: 'rgba(255, 255, 255, 0.3)',
    piece1: 'linear-gradient(135deg, #ff4444, #cc0000)',
    piece2: 'linear-gradient(135deg, #ffdd44, #ffaa00)',
    cardBg: 'rgba(255, 255, 255, 0.9)',
    shadow: 'rgba(0, 0, 0, 0.1)',
    glow: '0 4px 20px rgba(102, 126, 234, 0.3)',
  },
  dark: {
    name: 'Dark',
    primary: '#667eea',
    secondary: '#764ba2',
    accent: '#a78bfa',
    background: 'linear-gradient(135deg, #1a202c 0%, #2d3748 100%)',
    text: '#e2e8f0',
    textSecondary: '#cbd5e0',
    board: 'linear-gradient(135deg, #4a5568, #2d3748)',
    cell: 'rgba(255, 255, 255, 0.1)',
    piece1: 'linear-gradient(135deg, #ff6b6b, #ee5a6f)',
    piece2: 'linear-gradient(135deg, #ffd93d, #fcbf49)',
    cardBg: 'rgba(45, 55, 72, 0.9)',
    shadow: 'rgba(0, 0, 0, 0.3)',
    glow: '0 4px 20px rgba(102, 126, 234, 0.4)',
  },
  ocean: {
    name: 'Ocean',
    primary: '#06b6d4',
    secondary: '#0e7490',
    accent: '#22d3ee',
    background: 'linear-gradient(135deg, #0c4a6e 0%, #075985 50%, #0e7490 100%)',
    text: '#e0f2fe',
    textSecondary: '#bae6fd',
    board: 'linear-gradient(135deg, #0369a1, #0c4a6e)',
    cell: 'rgba(6, 182, 212, 0.2)',
    piece1: 'linear-gradient(135deg, #f97316, #ea580c)',
    piece2: 'linear-gradient(135deg, #22d3ee, #06b6d4)',
    cardBg: 'rgba(12, 74, 110, 0.8)',
    shadow: 'rgba(6, 182, 212, 0.5)',
    glow: '0 0 30px rgba(34, 211, 238, 0.6)',
  },
};

export const ThemeProvider = ({ children }) => {
  const [currentTheme, setCurrentTheme] = useState('neonPurple');

  useEffect(() => {
    const savedTheme = localStorage.getItem('gameTheme');
    if (savedTheme && themes[savedTheme]) {
      setCurrentTheme(savedTheme);
    }
  }, []);

  const changeTheme = (themeName) => {
    setCurrentTheme(themeName);
    localStorage.setItem('gameTheme', themeName);
  };

  return (
    <ThemeContext.Provider value={{ theme: themes[currentTheme], currentTheme, changeTheme, themes }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within ThemeProvider');
  }
  return context;
};