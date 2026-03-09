import React, { useState, useEffect } from 'react';

const bootMessages = [
  { time: '0.000000', msg: 'Initializing Book2Shelf Terminal v2.0...' },
  { time: '0.102453', msg: 'Loading system configuration... OK' },
  { time: '0.234891', msg: 'Mounting filesystem: /library... OK' },
  { time: '0.456123', msg: 'Connecting to database... OK' },
  { time: '0.678234', msg: 'Loading book catalog... OK' },
  { time: '0.891456', msg: 'Starting web interface...' },
  { time: '1.234567', msg: 'System ready. Welcome, user.' },
];

function BootSequence({ onComplete }) {
  const [visibleLines, setVisibleLines] = useState([]);
  const [showLogo, setShowLogo] = useState(false);
  const [fadeOut, setFadeOut] = useState(false);
  const [phase, setPhase] = useState('lines'); // 'lines' | 'logo' | 'done'

  useEffect(() => {
    let isCancelled = false;
    let lineIndex = 0;
    let intervalId;
    let timeoutIds = [];

    const cleanup = () => {
      isCancelled = true;
      if (intervalId) clearInterval(intervalId);
      timeoutIds.forEach(id => clearTimeout(id));
    };

    // Show boot lines one by one
    intervalId = setInterval(() => {
      if (isCancelled) return;
      
      const nextLine = bootMessages[lineIndex];
      if (lineIndex < bootMessages.length && nextLine) {
        setVisibleLines(prev => [...prev, nextLine]);
        lineIndex++;
      } else {
        clearInterval(intervalId);
        intervalId = null;
        
        // Show logo after lines complete
        const t1 = setTimeout(() => {
          if (!isCancelled) {
            setShowLogo(true);
            setPhase('logo');
          }
        }, 400);
        timeoutIds.push(t1);
        
        // Start fade out
        const t2 = setTimeout(() => {
          if (!isCancelled) {
            setFadeOut(true);
          }
        }, 1800);
        timeoutIds.push(t2);
        
        // Complete boot sequence
        const t3 = setTimeout(() => {
          if (!isCancelled) {
            setPhase('done');
            sessionStorage.setItem('bootComplete', 'true');
            if (onComplete) onComplete();
          }
        }, 2300);
        timeoutIds.push(t3);
      }
    }, 150);

    return cleanup;
  }, []); // Empty deps - run only once

  // Don't render if complete
  if (phase === 'done') {
    return null;
  }

  return (
    <div className={`boot-sequence ${fadeOut ? 'fade-out' : ''}`}>
      <div className="boot-container">
        {showLogo && (
          <pre className="boot-logo">
{`
  ____              _    ____  ____  _          _  __ 
 | __ )  ___   ___ | | _|___ \\/ ___|| |__   ___| |/ _|
 |  _ \\ / _ \\ / _ \\| |/ / __) \\___ \\| '_ \\ / _ \\ | |_ 
 | |_) | (_) | (_) |   < / __/ ___) | | | |  __/ |  _|
 |____/ \\___/ \\___/|_|\\_\\_____|____/|_| |_|\\___|_|_|  
                                                      
`}
          </pre>
        )}
        <div className="boot-lines">
          {visibleLines
            .filter(Boolean)
            .map((line, i) => (
              <div key={i} className="boot-line">
                <span className="boot-time">[{line.time}]</span>
                <span className="boot-msg">{line.msg}</span>
                {line.msg.includes('OK') && <span className="boot-ok"> ✓</span>}
              </div>
            ))}
        </div>
        {showLogo && (
          <div className="boot-cursor">
            <span className="cursor-prompt">root@book2shelf:~$</span>
            <span className="cursor-blink">_</span>
          </div>
        )}
      </div>
    </div>
  );
}

export default BootSequence;
