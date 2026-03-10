import { useTheme } from '../ThemeContext';

const SunIcon = () => (
  <svg
    width="18"
    height="18"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <circle cx="12" cy="12" r="5" stroke="currentColor" strokeWidth="1.8" />
    <path d="M12 2v2.5M12 19.5V22M4.5 12H2M22 12h-2.5M5.4 5.4 3.8 3.8M20.2 20.2l-1.6-1.6M5.4 18.6 3.8 20.2M20.2 3.8l-1.6 1.6" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
  </svg>
);

const MoonIcon = () => (
  <svg
    width="18"
    height="18"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <path
      d="M21 13.2A8.2 8.2 0 0 1 10.8 3a8.2 8.2 0 1 0 10.2 10.2Z"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export default function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const isDark = theme === 'dark';

  return (
    <button
      type="button"
      className="theme-btn"
      onClick={() => setTheme(isDark ? 'light' : 'dark')}
      title={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {isDark ? <SunIcon /> : <MoonIcon />}
    </button>
  );
}
