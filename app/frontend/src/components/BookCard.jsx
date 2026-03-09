import React from 'react';
import { PdfIcon, DownloadIcon, FilesIcon } from './Icons';

function BookCard({ book, onClick }) {
  const formatSize = (bytes) => {
    if (!bytes) return '0.0 MB';
    const mb = bytes / (1024 * 1024);
    return mb.toFixed(1) + ' MB';
  };

  const formatFilename = (title) => {
    return title.toLowerCase().replace(/\s+/g, '_').substring(0, 20);
  };

  const getFileExt = () => {
    if (book.file_name) {
      const dot = book.file_name.lastIndexOf('.');
      if (dot !== -1) return book.file_name.slice(dot);
    }
    return '.pdf';
  };

  const ext = getFileExt();
  const isPdf = ext.toLowerCase() === '.pdf';

  return (
    <div className="book-card terminal-card" onClick={onClick}>
      <div className="book-cover">
        {book.cover_url ? (
          <img src={book.cover_url} alt={book.title} />
        ) : (
          <div className="book-cover-placeholder">
            <span className="placeholder-icon">{isPdf ? <PdfIcon size={42} /> : <FilesIcon size={42} />}</span>
            <span className="placeholder-ext">{ext}</span>
          </div>
        )}
      </div>
      <div className="book-info">
        <div className="book-filename">
          <span className="file-prefix">&gt;</span> {formatFilename(book.title)}{ext}
        </div>
        <h3 className="book-title">{book.title}</h3>
        {book.author && <p className="book-author">-- {book.author}</p>}
        <div className="book-meta">
          <span className="book-downloads">
            <DownloadIcon size={14} /> <span className="meta-label">downloads</span> {book.downloads || 0}
          </span>
          <span className="book-size">[{formatSize(book.file_size)}]</span>
        </div>
      </div>
    </div>
  );
}

export default BookCard;
