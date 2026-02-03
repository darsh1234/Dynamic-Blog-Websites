import { Link, NavLink } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const APP_NAME = import.meta.env.VITE_APP_NAME || 'Go Gin Blog Platform';

export default function AppShell({ children }) {
  const auth = useAuth();
  const isAdmin = auth.user?.role === 'admin';

  return (
    <div className="app-root">
      <header className="topbar">
        <Link className="brand" to="/posts">
          {APP_NAME}
        </Link>
        <nav className="nav-links">
          {auth.isAuthenticated ? (
            <>
              <NavLink to="/posts">Posts</NavLink>
              {isAdmin && <NavLink to="/admin/users">Admin</NavLink>}
            </>
          ) : (
            <>
              <NavLink to="/login">Login</NavLink>
              <NavLink to="/register">Register</NavLink>
            </>
          )}
        </nav>
        <div className="auth-chip">
          {auth.isAuthenticated ? (
            <>
              <span>{auth.user.email}</span>
              <em>{auth.user.role}</em>
              <button type="button" onClick={auth.logout}>
                Logout
              </button>
            </>
          ) : (
            <span>Guest</span>
          )}
        </div>
      </header>
      <main className="content">{children}</main>
    </div>
  );
}
