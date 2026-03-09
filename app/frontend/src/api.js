import axios from 'axios';

const API_BASE = '/api';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('adminToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Public API
export const getBooks = () => api.get('/books');
export const getBook = (id) => api.get(`/books/${id}`);
export const downloadBook = (id) => api.get(`/books/${id}/download`, { responseType: 'blob' });
export const getSections = () => api.get('/sections');
export const getBooksBySection = (sectionId) => api.get(`/sections/${sectionId}/books`);

// Auth API
export const login = (credentials) => api.post('/login', credentials);

// Admin API
export const getAdminBooks = () => api.get('/admin/books');
export const getAdminSections = () => api.get('/admin/sections');
export const createBook = (book) => api.post('/admin/books', book);
export const updateBook = (id, book) => api.put(`/admin/books/${id}`, book);
export const deleteBook = (id) => api.delete(`/admin/books/${id}`);

export const createSection = (section) => api.post('/admin/sections', section);
export const updateSection = (id, section) => api.put(`/admin/sections/${id}`, section);
export const deleteSection = (id) => api.delete(`/admin/sections/${id}`);

export const uploadBookFile = (file) => {
  const formData = new FormData();
  formData.append('file', file);
  return api.post('/admin/upload/book', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
};

export const uploadCoverFile = (file) => {
  const formData = new FormData();
  formData.append('file', file);
  return api.post('/admin/upload/cover', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
};

export const cropCover = (payload) => api.post('/admin/cover/crop', payload);

export default api;
