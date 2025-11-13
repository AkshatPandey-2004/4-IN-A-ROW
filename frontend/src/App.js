import React, { useState, useEffect } from 'react';
import useWebSocket from './hooks/useWebSocket';
import GameBoard from './components/GameBoard';
import Leaderboard from './components/Leaderboard';
import ThemeSelector from './components/ThemeSelector';
import { useTheme } from './contexts/ThemeContext';
import './App.css';

function App() {
  const [username, setUsername] = useState('');
  const [inputUsername, setInputUsername] = useState('');
  const [gameState, setGameState] = useState(null);
  const [status, setStatus] = useState('idle');
  const [message, setMessage] = useState('');
  const [refreshLeaderboard, setRefreshLeaderboard] = useState(0);
  const [matchmakingTime, setMatchmakingTime] = useState(0);

  const { messages, isConnected, sendMessage } = useWebSocket(username);
  const { theme } = useTheme();

  // Matchmaking timer
  useEffect(() => {
    let interval;
    if (status === 'searching') {
      setMatchmakingTime(0);
      interval = setInterval(() => {
        setMatchmakingTime(prev => prev + 1);
      }, 1000);
    } else {
      setMatchmakingTime(0);
    }
    return () => clearInterval(interval);
  }, [status]);

  useEffect(() => {
    document.documentElement.style.setProperty('--primary', theme.primary);
    document.documentElement.style.setProperty('--secondary', theme.secondary);
    document.documentElement.style.setProperty('--accent', theme.accent);
    document.documentElement.style.setProperty('--background', theme.background);
    document.documentElement.style.setProperty('--text', theme.text);
    document.documentElement.style.setProperty('--text-secondary', theme.textSecondary);
    document.documentElement.style.setProperty('--board-bg', theme.board);
    document.documentElement.style.setProperty('--cell-bg', theme.cell);
    document.documentElement.style.setProperty('--piece1', theme.piece1);
    document.documentElement.style.setProperty('--piece2', theme.piece2);
    document.documentElement.style.setProperty('--card-bg', theme.cardBg);
    document.documentElement.style.setProperty('--shadow', theme.shadow);
    document.documentElement.style.setProperty('--glow', theme.glow);
  }, [theme]);

  useEffect(() => {
    if (messages.length === 0) return;

    const lastMessage = messages[messages.length - 1];

    switch (lastMessage.type) {
      case 'game_start':
        setGameState(lastMessage.game);
        setStatus('playing');
        setMessage(
          lastMessage.game.is_bot
            ? ' Playing against Bot'
            : ' Match Found!'
        );
        break;

      case 'move_made':
        setGameState(lastMessage.game);
        break;

      case 'game_end':
        setGameState(lastMessage.game);
        setStatus('finished');
        setRefreshLeaderboard(prev => prev + 1);
        
        if (lastMessage.result === 'draw') {
          setMessage(" It's a Draw!");
        } else if (lastMessage.winner) {
          const isWinner = lastMessage.winner.username === username;
          setMessage(isWinner ? ' Victory!' : ' Try Again!');
        }
        break;

      case 'error':
        setMessage(`❌ ${lastMessage.message}`);
        break;

      default:
        break;
    }
  }, [messages, username]);

  const handleLogin = (e) => {
    e.preventDefault();
    if (inputUsername.trim()) {
      // Clear any previous state before logging in
      setGameState(null);
      setStatus('idle');
      setMessage('');
      setMatchmakingTime(0);
      setUsername(inputUsername.trim());
    }
  };

  const handleFindMatch = () => {
    setStatus('searching');
    setMessage('  Searching...');
    sendMessage({ type: 'find_match' });
  };

  const handleColumnClick = (col) => {
    if (status !== 'playing') return;
    if (gameState.current_turn !== (username === gameState.player1.username ? 1 : 2)) {
      return;
    }
    sendMessage({
      type: 'make_move',
      column: col,
    });
  };

  const handlePlayAgain = () => {
    setGameState(null);
    setStatus('idle');
    setMessage('');
  };

  const handleLogout = () => {
    // Clear all game state
    setGameState(null);
    setStatus('idle');
    setMessage('');
    setMatchmakingTime(0);
    
    // Clear username to trigger WebSocket disconnect
    setUsername('');
    setInputUsername('');
    
    // Refresh leaderboard
    setRefreshLeaderboard(prev => prev + 1);
  };

  const formatTime = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (!username) {
    return (
      <div className="App">
        <ThemeSelector />
        <div className="login-page-wrapper">
          <div className="login-page">
            <div className="login-container">
              <div className="neon-title">
                <h1>🎮 4 IN A ROW</h1>
                <p className="subtitle">Connect Four • Neon Edition</p>
              </div>
              <form onSubmit={handleLogin}>
                <input
                  type="text"
                  placeholder="Enter your username"
                  value={inputUsername}
                  onChange={(e) => setInputUsername(e.target.value)}
                  className="username-input"
                  maxLength={20}
                  autoFocus
                />
                <button type="submit" className="btn btn-primary btn-glow">
                  START GAME
                </button>
              </form>
            </div>

            {/* How to Play Section */}
            <div className="how-to-play-login">
              <h2>📖 HOW TO PLAY</h2>
              <div className="rules-grid">
                <div className="rule-item">
                  <div className="rule-icon"></div>
                  <div className="rule-text">
                    <h3>Drop Your Piece</h3>
                    <p>Click any cell to drop your piece in that column</p>
                  </div>
                </div>
                <div className="rule-item">
                  <div className="rule-icon">🎮</div>
                  <div className="rule-text">
                    <h3>Connect 4</h3>
                    <p>Line up 4 pieces in a row to win the game</p>
                  </div>
                </div>
                <div className="rule-item">
                  <div className="rule-icon"></div>
                  <div className="rule-text">
                    <h3>Any Direction</h3>
                    <p>Win with horizontal, vertical, or diagonal lines</p>
                  </div>
                </div>
                <div className="rule-item">
                  <div className="rule-icon"></div>
                  <div className="rule-text">
                    <h3>Multiplayer</h3>
                    <p>Play against other players or challenge the Bot</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Leaderboard on Login Page */}
          <div className="login-leaderboard">
            <Leaderboard key={refreshLeaderboard} maxItems={7} />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      <div className="game-page">
        {/* Header */}
        <header className="game-header">
          <div className="header-left">
            <h1 className="neon-text">4 IN A ROW</h1>
            <div className="player-info">
              <span className="player-label">Player:</span>
              <span className="player-name">{username}</span>
            </div>
          </div>
          <div className="header-right">
            <div className={`status-badge ${isConnected ? 'connected' : 'disconnected'}`}>
              {isConnected ? '🟢 ONLINE' : '🔴 OFFLINE'}
            </div>
            <button onClick={handleLogout} className="btn btn-logout">
               LOGOUT
            </button>
            <ThemeSelector />
          </div>
        </header>

        {/* Main Content - Always 3 Columns */}
        <div className="game-content">
          {/* Left Panel - Always Show Leaderboard */}
          <div className="left-panel">
            <Leaderboard key={refreshLeaderboard} maxItems={7} compact={status === 'playing'} />
          </div>

          {/* Center Panel - Game */}
          <div className="center-panel">
            {status === 'idle' && (
              <div className="game-idle">
                <button onClick={handleFindMatch} className="btn btn-primary btn-glow btn-huge">
                  FIND MATCH
                </button>
                <p className="hint-text">Ready to play? Click to start!</p>
              </div>
            )}

            {status === 'searching' && (
              <div className="game-searching">
                <div className="neon-spinner"></div>
                <p className="searching-text">{message}</p>
                <p className="searching-time">Time: {formatTime(matchmakingTime)}</p>
                <p className="searching-hint">Finding opponent...</p>
              </div>
            )}

            {(status === 'playing' || status === 'finished') && gameState && (
              <div className="game-active">
                <div className={`game-status ${status === 'finished' ? 'finished' : ''}`}>
                  <p className="status-message">{message}</p>
                  {status === 'playing' && (
                    <div className="turn-display">
                      <div className={`turn-indicator ${gameState.current_turn === 1 ? 'player1' : 'player2'}`}>
                        <div className="pulse-dot"></div>
                        <span>
                          {gameState.current_turn === gameState.player1.piece
                            ? `${gameState.player1.username}'s Turn`
                            : `${gameState.player2.username}'s Turn`}
                        </span>
                      </div>
                    </div>
                  )}
                </div>

                <GameBoard
                  board={gameState.board}
                  onColumnClick={handleColumnClick}
                  currentTurn={gameState.current_turn}
                  myPiece={username === gameState.player1.username ? 1 : 2}
                  gameStatus={status}
                />

                {status === 'finished' && (
                  <div className="game-end-actions">
                    <button onClick={handlePlayAgain} className="btn btn-primary btn-glow btn-large">
                       PLAY AGAIN
                    </button>
                    <button onClick={handleLogout} className="btn btn-secondary btn-large">
                       BACK TO HOME
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Right Panel - Game Info */}
          <div className="right-panel">
            <div className="info-card">
              <h3>GAME INFO</h3>
              {gameState && status === 'playing' && (
                <div className="players-info">
                  <div className="player-card player1">
                    <div className="player-avatar">🔴</div>
                    <div>
                      <p className="player-title">Player 1</p>
                      <p className="player-username">{gameState.player1.username}</p>
                    </div>
                  </div>
                  <div className="vs-divider">VS</div>
                  <div className="player-card player2">
                    <div className="player-avatar">🟡</div>
                    <div>
                      <p className="player-title">Player 2</p>
                      <p className="player-username">
                        {gameState.player2 ? gameState.player2.username : 'Waiting...'}
                      </p>
                    </div>
                  </div>
                </div>
              )}
              
              {(status === 'idle' || status === 'searching') && (
                <div className="game-rules">
                  <h4>HOW TO PLAY</h4>
                  <ul>
                    <li> Click any cell to drop your piece</li>
                    <li> Connect 4 in a row to win</li>
                    <li> Horizontal, vertical, or diagonal</li>
                    <li> Play vs Bot or other players</li>
                  </ul>
                </div>
              )}

              {status === 'finished' && (
                <div className="game-stats">
                  <h4> GAME OVER</h4>
                  <p className="stats-text">Check the leaderboard for updated rankings!</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;