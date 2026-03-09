import { Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import AdminLogin from './pages/AdminLogin';
import AdminDashboard from './pages/AdminDashboard';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/book2shadmin" element={<AdminLogin />} />
      <Route path="/book2shadmin/dashboard/*" element={<AdminDashboard />} />
    </Routes>
  );
}

export default App;
