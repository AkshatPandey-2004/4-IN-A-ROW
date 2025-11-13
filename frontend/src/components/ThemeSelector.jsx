import React, { useState } from 'react';
import { useTheme } from '../contexts/ThemeContext';
import './ThemeSelector.css';

const ThemeSelector = () => {
  const { currentTheme, changeTheme, themes } = useTheme();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="theme-selector">
      <button 
        className="theme-toggle-btn" 
        onClick={() => setIsOpen(!isOpen)}
        title="Change Theme"
      >
        🎨
      </button>
      
      {isOpen && (
        <div className="theme-dropdown">
          <div className="theme-dropdown-header">
            <span>Choose Theme</span>
            <button onClick={() => setIsOpen(false)} className="close-btn">×</button>
          </div>
          <div className="theme-options">
            {Object.entries(themes).map(([key, theme]) => (
              <button
                key={key}
                className={`theme-option ${currentTheme === key ? 'active' : ''}`}
                onClick={() => {
                  changeTheme(key);
                  setIsOpen(false);
                }}
              >
                <div 
                  className="theme-preview" 
                  style={{ background: theme.primary }}
                ></div>
                <span>{theme.name}</span>
                {currentTheme === key && <span className="check-mark">✓</span>}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default ThemeSelector;