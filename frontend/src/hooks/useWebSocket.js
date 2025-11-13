import { useEffect, useRef, useState } from 'react';

const useWebSocket = (username) => {
  const [messages, setMessages] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const ws = useRef(null);

  useEffect(() => {
    // Clear messages when username changes or is cleared
    setMessages([]);
    
    if (!username) {
      // Close existing connection when logging out
      if (ws.current) {
        ws.current.close();
        ws.current = null;
      }
      setIsConnected(false);
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?username=${encodeURIComponent(username)}`;
    
    // Create new WebSocket connection
    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log('WebSocket Connected');
      setIsConnected(true);
    };

    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages((prev) => [...prev, message]);
    };

    ws.current.onerror = (error) => {
      console.error('WebSocket Error:', error);
    };

    ws.current.onclose = () => {
      console.log('WebSocket Disconnected');
      setIsConnected(false);
    };

    // Cleanup function
    return () => {
      if (ws.current) {
        ws.current.close();
        ws.current = null;
      }
      setMessages([]);
      setIsConnected(false);
    };
  }, [username]);

  const sendMessage = (message) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected');
    }
  };

  return { messages, isConnected, sendMessage };
};

export default useWebSocket;