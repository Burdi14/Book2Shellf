import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  getAdminBooks,
  getAdminSections,
  deleteBook,
  deleteSection,
  createSection,
  updateSection,
} from '../../api';
import { BookIcon, FolderIcon, DownloadIcon, EditIcon, TrashIcon } from '../../components/Icons';
import { formatSize } from '../../utils/formatSize';

function LibraryManager() {
  const navigate = useNavigate();
  const [books, setBooks] = useState([]);
  const [sections, setSections] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('books');
  const [filterSection, setFilterSection] = useState('');
  const [showSectionForm, setShowSectionForm] = useState(false);
  const [sectionForm, setSectionForm] = useState({ name: '', description: '', hidden: false });
  const [editingSectionId, setEditingSectionId] = useState(null);
  const [copiedToken, setCopiedToken] = useState(null);

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

  const handleDeleteBook = async (id) => {
    if (!window.confirm('Delete this book?')) return;
    try {
      await deleteBook(id);
      loadData();
    } catch (error) {
      console.error('Failed to delete book:', error);
    }
  };

  const handleSectionSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingSectionId) {
        await updateSection(editingSectionId, sectionForm);
      } else {
        await createSection(sectionForm);
      }
      setSectionForm({ name: '', description: '', hidden: false });
      setShowSectionForm(false);
      setEditingSectionId(null);
      loadData();
    } catch (error) {
      alert('Failed to save section');
    }
  };

  const handleEditSection = (section) => {
    setSectionForm({ name: section.name, description: section.description || '', hidden: !!section.hidden });
    setEditingSectionId(section.id);
    setShowSectionForm(true);
  };

  const handleDeleteSection = async (id) => {
    if (!window.confirm('Delete this section? Books will become uncategorized.')) return;
    try {
      await deleteSection(id);
      loadData();
    } catch (error) {
      console.error('Failed to delete section:', error);
    }
  };

  const hiddenSectionIds = new Set(sections.filter((s) => s.hidden).map((s) => s.id));

  const copyShareLink = (token) => {
    const url = `${window.location.origin}/api/share/${token}`;
    navigator.clipboard.writeText(url);
    setCopiedToken(token);
    setTimeout(() => setCopiedToken(null), 2000);
  };

  const copyDownloadLink = (book) => {
    const isHidden = hiddenSectionIds.has(book.section_id);
    const url = isHidden && book.share_token
      ? `${window.location.origin}/api/share/${book.share_token}`
      : `${window.location.origin}/api/books/${book.id}/download`;
    navigator.clipboard.writeText(url);
    setCopiedToken(book.id);
    setTimeout(() => setCopiedToken(null), 2000);
  };

  const filteredBooks = filterSection ? books.filter((b) => b.section_id === filterSection) : books;

  if (loading) {
    return (
      <div className="loading">
        <div className="loading-spinner"></div>
      </div>
    );
  }

  return (
    <div>
      <div className="admin-header">
        <h1 className="admin-title">
          <span style={{ color: 'var(--red-primary)' }}>//</span> Library Manager
        </h1>
        <button className="btn btn-primary" onClick={() => navigate('/book2shadmin/dashboard/books/new')}>
          + Add Book
        </button>
      </div>

      <div className="tab-strip">
        <button onClick={() => setActiveTab('books')} className={activeTab === 'books' ? 'tab active' : 'tab'}>
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '6px' }}><BookIcon size={16} /> Books ({books.length})</span>
        </button>
        <button onClick={() => setActiveTab('sections')} className={activeTab === 'sections' ? 'tab tab-red active' : 'tab tab-red'}>
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: '6px' }}><FolderIcon size={16} /> Sections ({sections.length})</span>
        </button>
      </div>

      {activeTab === 'books' ? (
        <>
          <div style={{ marginBottom: '20px', display: 'flex', gap: '15px', alignItems: 'center' }}>
            <span style={{ color: 'var(--text-muted)', fontSize: '13px' }}>Filter:</span>
            <select
              value={filterSection}
              onChange={(e) => setFilterSection(e.target.value)}
              style={{
                padding: '8px 12px',
                background: 'var(--bg-secondary)',
                border: '1px solid var(--border-color)',
                color: 'var(--text-primary)',
                fontFamily: 'var(--font-mono)',
                fontSize: '13px',
              }}
            >
              <option value="">All Sections</option>
              {sections.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name} ({s.book_count || 0})
                </option>
              ))}
            </select>
          </div>

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
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredBooks.map((book) => (
                  <tr key={book.id}>
                    <td>
                      <span className="pill pill-red subtle">{book.section_name || 'Uncategorized'}</span>
                      {hiddenSectionIds.has(book.section_id) && (
                        <span className="pill" style={{ marginLeft: '4px', fontSize: '9px', background: 'var(--text-muted)', color: 'var(--bg-primary)' }}>H</span>
                      )}
                    </td>
                    <td style={{ color: 'var(--accent-primary)' }}>{book.title}</td>
                    <td>{book.author || '—'}</td>
                    <td>{book.year || '—'}</td>
                    <td>{formatSize(book.file_size)}</td>
                    <td>{book.downloads}</td>
                    <td className="table-actions">
                      <button
                        className="btn btn-icon"
                        onClick={() => copyDownloadLink(book)}
                        title="Copy download link"
                        style={copiedToken === book.id ? { color: 'var(--green, #4caf50)' } : {}}
                      >
                        {copiedToken === book.id ? '✓' : '🔗'}
                      </button>
                      <button
                        className="btn btn-icon"
                        onClick={() => navigate(`/book2shadmin/dashboard/books/edit/${book.id}`)}
                        title="Edit"
                      >
                        <EditIcon size={14} />
                      </button>
                      <button className="btn btn-icon" onClick={() => handleDeleteBook(book.id)} title="Delete">
                        <TrashIcon size={14} />
                      </button>
                    </td>
                  </tr>
                ))}
                {filteredBooks.length === 0 && (
                  <tr>
                    <td colSpan="7" style={{ textAlign: 'center', padding: '40px', color: 'var(--text-muted)' }}>
                      No books found
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </>
      ) : (
        <>
          <div style={{ marginBottom: '20px' }}>
            <button
              className="btn btn-outline"
              onClick={() => {
                setShowSectionForm(true);
                setEditingSectionId(null);
                setSectionForm({ name: '', description: '', hidden: false });
              }}
              style={{ borderColor: 'var(--red-primary)', color: 'var(--red-primary)' }}
            >
              + Add Section
            </button>
          </div>

          {showSectionForm && (
            <div
              style={{
                background: 'var(--bg-card)',
                border: '1px solid var(--red-primary)',
                padding: '20px',
                marginBottom: '25px',
              }}
            >
              <form onSubmit={handleSectionSubmit}>
                <div className="form-row">
                  <div className="form-group">
                    <label className="form-label">Section Name *</label>
                    <input
                      type="text"
                      value={sectionForm.name}
                      onChange={(e) => setSectionForm((prev) => ({ ...prev, name: e.target.value }))}
                      placeholder="e.g., Programming, Science Fiction"
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label className="form-label">Description</label>
                    <input
                      type="text"
                      value={sectionForm.description}
                      onChange={(e) => setSectionForm((prev) => ({ ...prev, description: e.target.value }))}
                      placeholder="Optional description"
                    />
                  </div>
                </div>
                <div style={{ marginBottom: '15px' }}>
                  <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer', color: 'var(--text-primary)', fontSize: '13px' }}>
                    <input
                      type="checkbox"
                      checked={sectionForm.hidden}
                      onChange={(e) => setSectionForm((prev) => ({ ...prev, hidden: e.target.checked }))}
                      style={{ accentColor: 'var(--red-primary)' }}
                    />
                    Hidden — books in this section won't appear on the public site
                  </label>
                </div>
                <div style={{ display: 'flex', gap: '10px' }}>
                  <button type="submit" className="btn btn-primary">
                    {editingSectionId ? 'Update Section' : 'Create Section'}
                  </button>
                  <button
                    type="button"
                    className="btn btn-outline"
                    onClick={() => {
                      setShowSectionForm(false);
                      setEditingSectionId(null);
                    }}
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          )}

          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Description</th>
                  <th>Books</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {sections.map((section) => (
                  <tr key={section.id} style={section.hidden ? { opacity: 0.7 } : {}}>
                    <td style={{ color: 'var(--red-primary)', fontWeight: '600' }}>
                      {section.name}
                      {section.hidden && (
                        <span className="pill" style={{ marginLeft: '8px', fontSize: '10px', background: 'var(--text-muted)', color: 'var(--bg-primary)' }}>
                          HIDDEN
                        </span>
                      )}
                    </td>
                    <td>{section.description || '—'}</td>
                    <td>
                      <span style={{ color: 'var(--accent-primary)', fontWeight: '600' }}>
                        {section.book_count || 0}
                      </span>
                    </td>
                    <td className="table-actions">
                      <button className="btn btn-icon" onClick={() => handleEditSection(section)} title="Edit">
                        <EditIcon size={14} />
                      </button>
                      <button className="btn btn-icon" onClick={() => handleDeleteSection(section.id)} title="Delete">
                        <TrashIcon size={14} />
                      </button>
                    </td>
                  </tr>
                ))}
                {sections.length === 0 && (
                  <tr>
                    <td colSpan="4" style={{ textAlign: 'center', padding: '40px', color: 'var(--text-muted)' }}>
                      No sections found. Create your first section!
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}

export default LibraryManager;
