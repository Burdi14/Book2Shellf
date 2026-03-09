import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import ThemeToggle from './ThemeToggle';

const ASCII_LOGO = `
 ____              _    ____  ____  _          _  __ 
| __ )  ___   ___ | | _|___ \\/ ___|| |__   ___| |/ _|
|  _ \\ / _ \\ / _ \\| |/ / __) \\___ \\| '_ \\ / _ \\ | |_ 
| |_) | (_) | (_) |   < / __/ ___) | | | |  __/ |  _|
|____/ \\___/ \\___/|_|\\_\\_____|____/|_| |_|\\___|_|_|  
`;

function Header() {
  const location = useLocation();
  
  return (
    <header className="terminal-window">
      <div className="terminal-titlebar">
        <div className="terminal-buttons">
          <span className="terminal-btn close"></span>
          <span className="terminal-btn minimize"></span>
          <span className="terminal-btn maximize"></span>
        </div>
        <div className="terminal-title">reader@book2shelf:~/library</div>
        <div className="terminal-spacer"></div>
      </div>
      <div className="terminal-body">
        <div className="container">
          <div className="header-content">
            <Link to="/" className="logo-section">
              <pre className="ascii-logo">{ASCII_LOGO}</pre>
              <div className="logo-subtitle">&gt;&gt; Burdi Library</div>
            </Link>
            <div className="header-actions">
              <nav className="nav terminal-nav">
                <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
                  <span className="cmd-prefix">/</span>library
                </Link>
                <a href="https://me.burdi.ru" target="_blank" rel="noopener noreferrer">
                  <span className="cmd-prefix">/</span>about
                </a>
                <a href="https://github.com/burdi14" target="_blank" rel="noopener noreferrer">
                  <span className="cmd-prefix">/</span>github
                </a>
              </nav>
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}

export default Header;
