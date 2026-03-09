import { useEffect } from 'react';
import { Routes, Route, Link, useLocation, useNavigate } from 'react-router-dom';
import ThemeToggle from '../components/ThemeToggle';
import DashboardHome from './admin/DashboardHome';
import LibraryManager from './admin/LibraryManager';
import BookForm from './admin/BookForm';

function AdminDashboard() {
  const navigate = useNavigate();
  const location = useLocation();

  const navItems = [
    { to: '/book2shadmin/dashboard', label: 'Dashboard' },
    { to: '/book2shadmin/dashboard/library', label: 'Library' },
    { to: '/book2shadmin/dashboard/books/new', label: 'Add Book' },
  ];

  useEffect(() => {
    if (!localStorage.getItem('adminToken')) {
      navigate('/book2shadmin');
    }
  }, [navigate]);

  return (
    <div className="admin-layout">
      <header className="admin-topbar">
        <div className="admin-topbar__left">
          <h1 className="admin-brand">book2sh_admin</h1>
          <nav className="admin-nav">
            {navItems.map((item) => {
              const active = location.pathname === item.to || location.pathname.startsWith(`${item.to}/`);
              return (
                <Link key={item.to} to={item.to} className={active ? 'nav-link active' : 'nav-link'}>
                  {active ? '> ' : ''}{item.label}
                </Link>
              );
            })}
          </nav>
        </div>
        <ThemeToggle />
      </header>

      <main className="admin-content">
        <Routes>
          <Route index element={<DashboardHome />} />
          <Route path="library" element={<LibraryManager />} />
          <Route path="books/new" element={<BookForm />} />
          <Route path="books/edit/:id" element={<BookForm />} />
        </Routes>
      </main>
    </div>
  );
}

export default AdminDashboard;
