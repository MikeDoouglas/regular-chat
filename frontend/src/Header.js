import React from 'react';
import './Header.css';

const Header = () => {
  return (
    <div className="header">
      <img src="logo.png" alt="Logo" className="header-logo" />
      <div className="header-texts">
        <h1 className="header-title">Regular Chat</h1>
        <p className="header-slogan">Just a regular chat!</p>
      </div>
    </div>
  );
};

export default Header;
