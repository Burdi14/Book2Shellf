import React from 'react';

function Footer() {
  return (
    <footer className="footer terminal-footer">
      <div className="container">
        <div className="footer-content">
          <div className="footer-left">
            <span className="prompt">reader@book2shelf:~$</span>
            <span className="footer-text">echo "Book2Shelf"</span>
          </div>
          <div className="footer-right">
            <span className="footer-link">&gt; <a href="https://me.burdi.ru" target="_blank" rel="noopener noreferrer">about</a></span>
            <span className="footer-link">&gt; <a href="https://github.com/burdi14" target="_blank" rel="noopener noreferrer">source</a></span>
          </div>
        </div>
        <div className="footer-output">
          Book2Shelf | Burdi Library | All books belong to their respective authors
        </div>
      </div>
    </footer>
  );
}

export default Footer;
