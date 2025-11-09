import React, { useState, useEffect } from 'react';
import useWebSocket from './hooks/useWebSocket';
import GameBoard from './components/GameBoard';
import Leaderboard from './components/Leaderboard';
import './App.css';

function App() {
  const [username, setUsername] = useState('');
  const [inputUsername, setInputUsername] = useState('');
  const [gameState, setGameState] = useState(null);
  const [status, setStatus] = useState('idle'); // idle, searching, playing, finished
  const [message, setMessage] = useState('');
  const [showLeaderboard, setShowLeaderboard] = useState(true);

  const { messages, isConnected, sendMessage } = useWebSocket(username);

  useEffect(() => {
    if (messages.length === 0) return;

    const lastMessage = messages[messages.length - 1];
    console.log('Received:', lastMessage);

    switch (lastMessage.type) {
      case 'game_start':
        setGameState(lastMessage.game);
        setStatus('playing');
        setMessage(
          lastMessage.game.is_bot
            ? '🤖 Playing against Bot'
            : '👥 Opponent found! Game started!'
        );
        break;

      case 'move_made':
        setGameState(lastMessage.game);
        break;

      case 'game_end':
        setGameState(lastMessage.game);
        setStatus('finished');
        
        if (lastMessage.result === 'draw') {
          setMessage("🤝 It's a draw!");
        } else if (lastMessage.winner) {
          const isWinner = lastMessage.winner.username === username;
          setMessage(isWinner ? '🎉 You won!' : '😢 You lost!');
        }
        break;

      case 'error':
        setMessage(`❌ Error: ${lastMessage.message}`);
        break;

      default:
        break;
    }
  }, [messages, username]);

  const handleLogin = (e) => {
    e.preventDefault();
    if (inputUsername.trim()) {
      setUsername(inputUsername.trim());
    }
  };

  const handleFindMatch = () => {
    setStatus('searching');
    setMessage('🔍 Searching for opponent...');
    sendMessage({ type: 'find_match' });
  };

  const handleColumnClick = (col) => {
    if (status !== 'playing') return;
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

  if (!username) {
    return (
      <div className="App">
        <div className="login-container">
          <h1>🎮 4 in a Row</h1>
          <form onSubmit={handleLogin}>
            <input
              type="text"
              placeholder="Enter your username"
              value={inputUsername}
              onChange={(e) => setInputUsername(e.target.value)}
              className="username-input"
            />
            <button type="submit" className="btn btn-primary">
              Start Playing
            </button>
          </form>
        </div>
        <Leaderboard />
      </div>
    );
  }

  return (
    <div className="App">
      <header>
        <h1>🎮 4 in a Row</h1>
        <div className="user-info">
          <span>Player: <strong>{username}</strong></span>
          <span className={`connection-status ${isConnected ? 'connected' : 'disconnected'}`}>
            {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
          </span>
        </div>
      </header>

      <main>
        {status === 'idle' && (
          <div className="game-controls">
            <button onClick={handleFindMatch} className="btn btn-primary">
              Find Match
            </button>
            <button onClick={() => setShowLeaderboard(!showLeaderboard)} className="btn btn-secondary">
              {showLeaderboard ? 'Hide' : 'Show'} Leaderboard
            </button>
          </div>
        )}

        {status === 'searching' && (
          <div className="searching">
            <div className="spinner"></div>
            <p>{message}</p>
          </div>
        )}

        {(status === 'playing' || status === 'finished') && gameState && (
          <div className="game-container">
            <div className="game-info">
              <p>{message}</p>
              {status === 'playing' && (
                <p className="turn-indicator">
                  {gameState.current_turn === gameState.player1.piece
                    ? `${gameState.player1.username}'s turn (Red)`
                    : `${gameState.player2.username}'s turn (Yellow)`}
                </p>
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
              <button onClick={handlePlayAgain} className="btn btn-primary">
                Play Again
              </button>
            )}
          </div>
        )}

        {showLeaderboard && <Leaderboard />}
      </main>
    </div>
  );
}

export default App;