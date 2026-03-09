import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getAdminBooks, getAdminSections } from '../../api';
import { BookIcon, FolderIcon, DownloadIcon, DriveIcon, FilesIcon } from '../../components/Icons';
import StatCard from '../../components/admin/StatCard';
import { formatSize } from '../../utils/formatSize';

function DashboardHome() {
  const [books, setBooks] = useState([]);
  const [sections, setSections] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [booksRes, sectionsRes] = await Promise.all([getAdminBooks(), getAdminSections()]);
      setBooks(booksRes.data.data || []);
      setSections(sectionsRes.data.data || []);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const downloads = books.reduce((sum, b) => sum + (b.downloads || 0), 0);
  const totalSize = books.reduce((sum, b) => sum + (b.file_size || 0), 0);

  if (loading) {
    return (
      <div className="loading">
        <div className="loading-spinner"></div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="admin-title">
        <span style={{ color: 'var(--red-primary)' }}>//</span> Dashboard
      </h1>

      <div className="stats-grid">
        <StatCard label="Total Books" value={books.length} icon={<BookIcon />} />
        <StatCard label="Sections" value={sections.length} icon={<FolderIcon />} />
        <StatCard label="Downloads" value={downloads} icon={<DownloadIcon />} />
        <StatCard label="Total Size" value={formatSize(totalSize)} icon={<DriveIcon />} />
      </div>

      <div style={{ marginTop: '40px' }}>
        <h2 className="section-title">
          <span style={{ display: 'inline-flex', color: 'var(--red-primary)' }}><FilesIcon size={16} /></span> Library Overview
        </h2>

        {books.length === 0 ? (
          <div className="card-empty">
            <p>$ ls ./library/</p>
            <p style={{ marginTop: '10px' }}>
              No books found. <Link to="/book2shadmin/dashboard/library" style={{ color: 'var(--accent-primary)' }}>Add your first book</Link>
            </p>
          </div>
        ) : (
          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Section</th>
                  <th>Title</th>
                  <th>Author</th>
                  <th>Year</th>
                  <th>Size</th>
                  <th style={{ display: 'flex', alignItems: 'center', gap: '6px' }}><DownloadIcon size={14} /> Downloads</th>
                </tr>
              </thead>
              <tbody>
                {books.slice(0, 10).map((book) => (
                  <tr key={book.id}>
                    <td>
                      <span className="pill pill-red">{book.section_name || 'Uncategorized'}</span>
                    </td>
                    <td style={{ color: 'var(--accent-primary)' }}>{book.title}</td>
                    <td>{book.author || '—'}</td>
                    <td>{book.year || '—'}</td>
                    <td>{formatSize(book.file_size)}</td>
                    <td>{book.downloads}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

export default DashboardHome;
