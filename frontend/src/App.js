import React, { useState, useEffect, useRef } from 'react';
import './App.css';
import './Header.js';
import useViewportHeight from './useViewportHeight.js';

import Header from './Header.js'

const App = () => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [userId, setUserId] = useState('');
  const [userName, setUserName] = useState('');
  const socket = useRef(null);
  const messagesEndRef = useRef(null);
  const dynamicHeight = useViewportHeight();

  const isValidUUID = (uuid) => {
    const regex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;
    return regex.test(uuid);
  }

  useEffect(() => {
    const wsUrl = process.env.REACT_APP_WS_URL || 'ws://localhost:8080/ws';
    socket.current = new WebSocket(wsUrl);

    socket.current.onopen = () => {
      console.log('WebSocket conectado');
    };

    socket.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (!message.text && !userId && !userName && isValidUUID(message.user_id)) {
        setUserId(message.user_id);
        setUserName(message.user_name);
      }

      if (message.text) {
        setMessages((prevMessages) => [...prevMessages, message]);
      }
    };

    socket.current.onclose = (event) => {
      console.log(`WebSocket desconectado: ${event.code} - ${event.reason}`);
    };

    socket.current.onerror = (error) => {
      console.error(`WebSocket erro: ${error.message}`);
    };

    return () => {
      if (socket.current) {
        socket.current.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  const sendMessage = () => {
    if (input.trim() && socket.current.readyState === WebSocket.OPEN) {
      const message = { user_name: userName, user_id: userId, text: input };
      socket.current.send(JSON.stringify(message));
      setMessages((prevMessages) => [...prevMessages, message]);
      setInput('');
    }
  };

  const handleKeyPress = (event) => {
    if (event.key === 'Enter') {
      event.preventDefault();
      sendMessage();
    }
  };

  return (
    <div className="container" style={{ height: `${dynamicHeight}px` }}>
      <Header />
      <div className="chat-container" style={{ height: `${dynamicHeight}px` }}>
        <div className="messages-container">
          {messages.map((msg, index) => (
            <div key={index} className={msg.user_id === userId ? "message your-message" : "message stranger-message"}>
              <strong>{msg.user_name}:</strong> {msg.text}
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
        <div className="input-container">
          <textarea
            type="text"
            value={input}
            rows={4}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyPress}
          />
          <button onClick={sendMessage}><strong>Send</strong></button>
        </div>
      </div>
    </div>
  );
};

export default App;
