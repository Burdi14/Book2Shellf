import React from 'react';
import { PdfIcon, DownloadIcon, SizeIcon, CalendarIcon, FilesIcon } from './Icons';

function BookModal({ book, onClose }) {
  const formatSize = (bytes) => {
    if (!bytes) return '0 MB';
    const mb = bytes / (1024 * 1024);
    return mb.toFixed(2) + ' MB';
  };

  const formatDate = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const slugify = (title) => title.toLowerCase().replace(/\s+/g, '_');

  const truncatedSlug = (title) => {
    const base = slugify(title);
    const limit = 30;
    if (base.length <= limit) return base;
    return base.slice(0, limit) + '…';
  };

  const handleDownload = () => {
    // Force a real download to avoid in-browser zoomed previews on some mobile devices
    const link = document.createElement('a');
    link.href = `/api/books/${book.id}/download`;
    link.setAttribute('download', book.file_name || 'book');
    link.rel = 'noreferrer';
    document.body.appendChild(link);
    link.click();
    link.remove();
  };

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div className="modal-overlay" onClick={handleOverlayClick}>
      <div className="modal terminal-modal">
          <div className="modal-titlebar">
          <div className="terminal-buttons">
            <span className="terminal-btn close" onClick={onClose}></span>
            <span className="terminal-btn minimize"></span>
            <span className="terminal-btn maximize"></span>
          </div>
            <div className="terminal-title">cat ./library/{truncatedSlug(book.title)}.info</div>
          <div className="terminal-spacer"></div>
        </div>
        
        <div className="modal-content">
          <div className="modal-cover">
            {book.cover_url ? (
              <img src={book.cover_url} alt={book.title} />
            ) : (
              <div className="modal-placeholder">
                <span style={{ fontSize: '64px', display: 'inline-flex' }}>
                  {book.file_name?.toLowerCase().endsWith('.pdf') ? <PdfIcon size={64} /> : <FilesIcon size={64} />}
                </span>
                <span style={{ color: 'var(--accent-primary)' }}>
                  {book.file_name ? book.file_name.slice(book.file_name.lastIndexOf('.')) : '.file'}
                </span>
              </div>
            )}
          </div>
          
          <div className="modal-info">
            <div className="modal-output-line">
              <span className="output-prefix">$</span>
              <span className="output-cmd">file {truncatedSlug(book.title)}{book.file_name ? book.file_name.slice(book.file_name.lastIndexOf('.')) : ''}</span>
            </div>
            
            <h2 className="modal-title">{book.title}</h2>
            {book.author && <p className="modal-author">&gt;&gt; Author: {book.author}</p>}
            
            {book.description && (
              <div className="modal-description-block">
                <div className="desc-header"># Description</div>
                <p className="modal-description">{book.description}</p>
              </div>
            )}
            
            <div className="modal-stats terminal-stats">
              <div className="modal-stat">
                <div className="modal-stat-label"><DownloadIcon size={14} /> downloads</div>
                <div className="modal-stat-value">{book.downloads || 0}</div>
              </div>
              <div className="modal-stat">
                <div className="modal-stat-label"><SizeIcon size={14} /> size</div>
                <div className="modal-stat-value">{formatSize(book.file_size)}</div>
              </div>
              <div className="modal-stat">
                <div className="modal-stat-label"><CalendarIcon size={14} /> year</div>
                <div className="modal-stat-value">{book.year}</div>
              </div>
            </div>
            
            {book.section_name && (
              <p className="modal-section">
                <span className="section-label">section:</span> 
                <span className="section-value">{book.section_name}</span>
              </p>
            )}
            
            <div className="modal-actions">
              <button 
                className="btn btn-primary download-btn terminal-btn-download"
                onClick={handleDownload}
                disabled={!book.file_url}
              >
                {book.file_url ? `$ wget book${book.file_name ? book.file_name.slice(book.file_name.lastIndexOf('.')) : ''}` : '# File not available'}
              </button>
              <button className="btn btn-outline modal-close-btn" onClick={onClose}>
                close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default BookModal;
