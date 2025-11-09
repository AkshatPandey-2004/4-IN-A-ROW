import React from 'react';
import './GameBoard.css';

const GameBoard = ({ board, onColumnClick, currentTurn, myPiece, gameStatus }) => {
  const handleColumnClick = (col) => {
    if (gameStatus !== 'playing') return;
    if (currentTurn !== myPiece) return;
    onColumnClick(col);
  };

  const getPieceColor = (piece) => {
    if (piece === 1) return 'red';
    if (piece === 2) return 'yellow';
    return 'empty';
  };

  return (
    <div className="game-board">
      {board && board.map((row, rowIndex) => (
        <div key={rowIndex} className="board-row">
          {row.map((cell, colIndex) => (
            <div
              key={`${rowIndex}-${colIndex}`}
              className={`board-cell ${getPieceColor(cell)}`}
              onClick={() => rowIndex === 0 && handleColumnClick(colIndex)}
              style={{ cursor: rowIndex === 0 && currentTurn === myPiece ? 'pointer' : 'default' }}
            >
              <div className="piece"></div>
            </div>
          ))}
        </div>
      ))}
    </div>
  );
};

export default GameBoard;