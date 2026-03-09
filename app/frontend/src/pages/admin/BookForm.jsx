import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  getAdminBooks,
  getAdminSections,
  createBook,
  updateBook,
  uploadBookFile,
  uploadCoverFile,
} from '../../api';
import { PdfIcon, ImageIcon } from '../../components/Icons';

function BookForm() {
  const navigate = useNavigate();
  const { id } = useParams();
  const isEditing = Boolean(id);

  const [sections, setSections] = useState([]);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({
    title: '',
    author: '',
    description: '',
    year: '',
    section_id: '',
    cover_url: '',
    file_url: '',
    file_name: '',
    file_size: 0,
  });
  const [autoCoverUrl, setAutoCoverUrl] = useState('');
  const [customCoverUrl, setCustomCoverUrl] = useState('');
  const [useCustomCover, setUseCustomCover] = useState(false);

  useEffect(() => {
    loadSections();
    if (isEditing && id) {
      loadBook(id);
    }
  }, [id, isEditing]);

  const loadSections = async () => {
    try {
      const res = await getAdminSections();
      setSections(res.data.data || []);
    } catch (error) {
      console.error('Failed to load sections:', error);
    }
  };

  const loadBook = async (bookId) => {
    try {
      const res = await getAdminBooks();
      const foundBook = (res.data.data || []).find((b) => b.id === bookId);
      if (!foundBook) return;
      setForm({
        title: foundBook.title || '',
        author: foundBook.author || '',
        description: foundBook.description || '',
        year: foundBook.year || '',
        section_id: foundBook.section_id || '',
        cover_url: foundBook.cover_url || '',
        file_url: foundBook.file_url || '',
        file_name: foundBook.file_name || '',
        file_size: foundBook.file_size || 0,
      });
      setCustomCoverUrl(foundBook.cover_url || '');
      setAutoCoverUrl('');
      setUseCustomCover(Boolean(foundBook.cover_url));
    } catch (error) {
      console.error('Failed to load book:', error);
    }
  };

  useEffect(() => {
    if (useCustomCover && customCoverUrl) {
      setForm((prev) => ({ ...prev, cover_url: customCoverUrl }));
    } else if (!useCustomCover && autoCoverUrl) {
      setForm((prev) => ({ ...prev, cover_url: autoCoverUrl }));
    }
  }, [useCustomCover, customCoverUrl, autoCoverUrl]);

  const handleUseCustomToggle = (checked) => {
    setUseCustomCover(checked);
    setForm((prev) => ({
      ...prev,
      cover_url: checked ? customCoverUrl || prev.cover_url : autoCoverUrl || prev.cover_url,
    }));
  };

  const handleFileUpload = async (e, type) => {
    const file = e.target.files[0];
    if (!file) return;
    try {
      setLoading(true);
      if (type === 'book') {
        const res = await uploadBookFile(file);
        const coverFromBook = res.data.data.cover?.url || res.data.data.cover_url || '';
        setForm((prev) => ({
          ...prev,
          file_url: res.data.data.url,
          file_name: res.data.data.original_name,
          file_size: res.data.data.size,
          cover_url: useCustomCover ? prev.cover_url : coverFromBook || prev.cover_url,
        }));
        if (coverFromBook) {
          setAutoCoverUrl(coverFromBook);
          if (!useCustomCover) {
            setForm((prev) => ({ ...prev, cover_url: coverFromBook }));
          }
        }
      } else {
        const res = await uploadCoverFile(file);
        const coverUrl = res.data.data.cover?.url || res.data.data.cover_url || res.data.data.url;
        setCustomCoverUrl(coverUrl);
        setForm((prev) => ({ ...prev, cover_url: useCustomCover ? coverUrl : prev.cover_url }));
      }
    } catch (error) {
      alert('Failed to upload file');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!form.title) {
      alert('Title is required');
      return;
    }

    const chosenCover = useCustomCover
      ? customCoverUrl || autoCoverUrl || form.cover_url
      : autoCoverUrl || form.cover_url || customCoverUrl;

    try {
      setLoading(true);
      const bookData = {
        ...form,
        cover_url: chosenCover || '',
        year: form.year ? parseInt(form.year, 10) : 0,
      };

      if (isEditing && id) {
        await updateBook(id, bookData);
      } else {
        await createBook(bookData);
      }
      navigate('/book2shadmin/dashboard/library');
    } catch (error) {
      alert('Failed to save book');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1 className="admin-title">
        <span style={{ color: 'var(--red-primary)' }}>//</span> {isEditing ? 'Edit Book' : 'Add New Book'}
      </h1>

      <form onSubmit={handleSubmit} style={{ maxWidth: '900px', marginTop: '30px' }}>
        <div className="form-row">
          <div className="form-group">
            <label className="form-label">Title *</label>
            <input
              type="text"
              value={form.title}
              onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))}
              placeholder="Book title"
              required
            />
          </div>
          <div className="form-group">
            <label className="form-label">Author</label>
            <input
              type="text"
              value={form.author}
              onChange={(e) => setForm((prev) => ({ ...prev, author: e.target.value }))}
              placeholder="Author name"
            />
          </div>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label className="form-label">Description</label>
            <textarea
              rows="3"
              value={form.description}
              onChange={(e) => setForm((prev) => ({ ...prev, description: e.target.value }))}
              placeholder="Short description"
            ></textarea>
          </div>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label className="form-label">Year</label>
            <input
              type="number"
              value={form.year}
              onChange={(e) => setForm((prev) => ({ ...prev, year: e.target.value }))}
              placeholder="YYYY"
            />
          </div>
          <div className="form-group">
            <label className="form-label">Section</label>
            <select
              value={form.section_id}
              onChange={(e) => setForm((prev) => ({ ...prev, section_id: e.target.value }))}
            >
              <option value="">Select section</option>
              {sections.map((section) => (
                <option key={section.id} value={section.id}>
                  {section.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label className="form-label">Book File *</label>
            <div className="file-upload">
              <label className="file-upload-label">
                <PdfIcon size={16} /> Upload File
                <input type="file" onChange={(e) => handleFileUpload(e, 'book')} />
              </label>
              {form.file_name && (
                <p className="file-upload-meta">
                  {form.file_name} ({Math.round((form.file_size || 0) / 1024)} KB)
                </p>
              )}
            </div>
          </div>

          <div className="form-group">
            <label className="form-label">Cover Image</label>
            <div className="file-upload">
              <label className="file-upload-label">
                <ImageIcon size={16} /> Upload Cover
                <input type="file" accept="image/*" onChange={(e) => handleFileUpload(e, 'cover')} />
              </label>
              <div className="toggle-row">
                <label className="toggle">
                  <input
                    type="checkbox"
                    checked={useCustomCover}
                    onChange={(e) => handleUseCustomToggle(e.target.checked)}
                  />
                  <span>Use custom cover</span>
                </label>
              </div>
              {(customCoverUrl || autoCoverUrl || form.cover_url) && (
                <div className="cover-preview">
                  <img src={customCoverUrl || autoCoverUrl || form.cover_url} alt="Cover preview" />
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="form-actions">
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {isEditing ? 'Update Book' : 'Create Book'}
          </button>
          <button type="button" className="btn btn-secondary" onClick={() => navigate(-1)} disabled={loading}>
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}

export default BookForm;
