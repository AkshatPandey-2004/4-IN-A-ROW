import React from 'react';
import { useTheme } from '../contexts/ThemeContext';
import './GameBoard.css';

const GameBoard = ({ board, onColumnClick, currentTurn, myPiece, gameStatus }) => {
  const { theme } = useTheme();
  
  const handleCellClick = (rowIndex, colIndex) => {
    if (gameStatus !== 'playing') return;
    if (currentTurn !== myPiece) return;
    
    // Find the lowest empty row in this column
    for (let row = board.length - 1; row >= 0; row--) {
      if (board[row][colIndex] === 0) {
        onColumnClick(colIndex);
        return;
      }
    }
  };

  const getPieceColor = (piece) => {
    if (piece === 1) return 'red';
    if (piece === 2) return 'yellow';
    return 'empty';
  };

  const isMyTurn = currentTurn === myPiece && gameStatus === 'playing';

  return (
    <div className={`game-board-container ${isMyTurn ? 'my-turn' : ''}`}>
      <div className="game-board">
        {board && board.map((row, rowIndex) => (
          <div key={rowIndex} className="board-row">
            {row.map((cell, colIndex) => (
              <div
                key={`${rowIndex}-${colIndex}`}
                className={`board-cell ${getPieceColor(cell)} ${isMyTurn ? 'clickable' : ''}`}
                onClick={() => handleCellClick(rowIndex, colIndex)}
              >
                <div className="piece"></div>
                {isMyTurn && cell === 0 && (
                  <div className="hover-indicator"></div>
                )}
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default GameBoard;