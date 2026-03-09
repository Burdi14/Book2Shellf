import React, { useState, useEffect } from 'react';
import { getBooks, getSections, getBooksBySection } from '../api';
import BookCard from '../components/BookCard';
import BookModal from '../components/BookModal';
import Header from '../components/Header';
import Footer from '../components/Footer';
import BootSequence from '../components/BootSequence';

function Home() {
  const [books, setBooks] = useState([]);
  const [sections, setSections] = useState([]);
  const [activeSection, setActiveSection] = useState(null);
  const [selectedBook, setSelectedBook] = useState(null);
  const [loading, setLoading] = useState(true);
  
  // Check sessionStorage on mount to determine if boot sequence should show
  const [showBootSequence, setShowBootSequence] = useState(
    () => !sessionStorage.getItem('bootComplete')
  );

  useEffect(() => {
    loadInitialData();
  }, []);

  useEffect(() => {
    if (activeSection) {
      loadBooksBySection(activeSection);
    } else {
      loadAllBooks();
    }
  }, [activeSection]);

  const loadInitialData = async () => {
    try {
      const [booksRes, sectionsRes] = await Promise.all([
        getBooks(),
        getSections()
      ]);
      setBooks(booksRes.data.data || []);
      setSections(sectionsRes.data.data || []);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadAllBooks = async () => {
    setLoading(true);
    try {
      const res = await getBooks();
      setBooks(res.data.data || []);
    } catch (error) {
      console.error('Failed to load books:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadBooksBySection = async (sectionId) => {
    setLoading(true);
    try {
      const res = await getBooksBySection(sectionId);
      setBooks(res.data.data || []);
    } catch (error) {
      console.error('Failed to load books:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSectionClick = (sectionId) => {
    setActiveSection(activeSection === sectionId ? null : sectionId);
  };

  // Show boot sequence if not completed yet
  if (showBootSequence) {
    return <BootSequence onComplete={() => setShowBootSequence(false)} />;
  }

  return (
    <div className="app terminal-app">
      <Header />
      
      <main className="main">
        <div className="container">
          <div className="layout">
            {/* Sidebar with sections */}
            <aside className="sidebar terminal-sidebar">
              <div className="sidebar-header">
                <span className="prompt-symbol">$</span>
                <span className="sidebar-title">ls ./sections/</span>
              </div>
              <ul className="section-list">
                <li 
                  className={`section-item ${!activeSection ? 'active' : ''}`}
                  onClick={() => setActiveSection(null)}
                >
                  <span className="section-prefix">&gt;</span>
                  <span>all_books</span>
                  <span className="section-count">
                    [{sections.reduce((sum, s) => sum + (s.book_count || 0), 0)}]
                  </span>
                </li>
                {sections.map((section) => (
                  <li 
                    key={section.id}
                    className={`section-item ${activeSection === section.id ? 'active' : ''}`}
                    onClick={() => handleSectionClick(section.id)}
                  >
                    <span className="section-prefix">&gt;</span>
                    <span>{section.name.toLowerCase().replace(/\s+/g, '_')}</span>
                    <span className="section-count">[{section.book_count || 0}]</span>
                  </li>
                ))}
              </ul>
            </aside>

            {/* Books grid */}
            <div className="content">
              <h2 className="section-title">
                <span className="cmd-output">&gt;&gt;</span>
                {activeSection 
                  ? `cat ${sections.find(s => s.id === activeSection)?.name.toLowerCase().replace(/\s+/g, '_') || 'books'}/*`
                  : 'cat ./library/*'
                }
              </h2>
              
              {loading ? (
                <div className="loading">
                  <div className="loading-text">Loading...</div>
                  <div className="loading-spinner"></div>
                </div>
              ) : books.length === 0 ? (
                <div className="empty-state terminal-empty">
                  <div className="empty-output">$ ls ./books/</div>
                  <p className="empty-msg">No files found in this directory</p>
                </div>
              ) : (
                <div className="book-grid">
                  {books.map((book) => (
                    <BookCard 
                      key={book.id} 
                      book={book} 
                      onClick={() => setSelectedBook(book)}
                    />
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </main>

      <Footer />

      {/* Book detail modal */}
      {selectedBook && (
        <BookModal 
          book={selectedBook} 
          onClose={() => setSelectedBook(null)}
        />
      )}
    </div>
  );
}

export default Home;
