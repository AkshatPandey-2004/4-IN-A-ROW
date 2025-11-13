import React, { useEffect, useState } from 'react';
import './Leaderboard.css';

const Leaderboard = ({ compact = false, maxItems = 7 }) => {
  const [leaderboard, setLeaderboard] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchLeaderboard();
  }, []);

  const fetchLeaderboard = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/leaderboard');
      const data = await response.json();
      // Limit to maxItems
      setLeaderboard((data || []).slice(0, maxItems));
    } catch (error) {
      console.error('Error fetching leaderboard:', error);
    } finally {
      setLoading(false);
    }
  };

  const getWinRate = (wins, totalGames) => {
    if (totalGames === 0) return '0%';
    return ((wins / totalGames) * 100).toFixed(1) + '%';
  };

  if (loading) {
    return (
      <div className={`leaderboard ${compact ? 'compact' : ''}`}>
        <h2>🏆 LEADERBOARD</h2>
        <div className="loading-spinner">
          <div className="spinner-small"></div>
          <p>Loading...</p>
        </div>
      </div>
    );
  }

  if (leaderboard.length === 0) {
    return (
      <div className={`leaderboard ${compact ? 'compact' : ''}`}>
        <h2>🏆 LEADERBOARD</h2>
        <p className="empty-message">No games played yet!</p>
      </div>
    );
  }

  return (
    <div className={`leaderboard ${compact ? 'compact' : ''}`}>
      <h2>🏆 LEADERBOARD</h2>
      <div className="table-container">
        <table>
          <thead>
            <tr>
              <th>Rank</th>
              <th>Player</th>
              <th>Wins</th>
              <th>Losses</th>
              <th>Draws</th>
              {!compact && <th>Total</th>}
              {!compact && <th>Win %</th>}
            </tr>
          </thead>
          <tbody>
            {leaderboard.map((entry, index) => (
              <tr key={entry.username} className={`row-${index < 3 ? 'medal' : 'normal'}`}>
                <td className="rank">#{index + 1}</td>
                <td className="username">{entry.username}</td>
                <td className="wins">{entry.wins}</td>
                <td className="losses">{entry.losses}</td>
                <td className="draws">{entry.draws}</td>
                {!compact && <td className="total">{entry.total_games}</td>}
                {!compact && <td className="win-rate">{getWinRate(entry.wins, entry.total_games)}</td>}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Leaderboard;