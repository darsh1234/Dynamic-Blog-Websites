import { useMemo, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function ResetPasswordPage() {
  const auth = useAuth();
  const location = useLocation();
  const token = useMemo(() => new URLSearchParams(location.search).get('token') || '', [location.search]);

  const [form, setForm] = useState({ password: '', confirmPassword: '' });
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const onSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setMessage('');

    if (!token) {
      setError('Missing reset token in URL');
      return;
    }

    if (form.password !== form.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setLoading(true);
    try {
      const payload = await auth.confirmPasswordReset(token, form.password);
      setMessage(payload?.message || 'Password updated successfully');
    } catch (err) {
      setError(err.message || 'Could not reset password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="card auth-card">
      <h1>Set a New Password</h1>
      <p>Use the token in your URL to complete password reset.</p>
      <form onSubmit={onSubmit}>
        <label>
          New Password
          <input
            type="password"
            minLength={8}
            value={form.password}
            onChange={(event) => setForm((curr) => ({ ...curr, password: event.target.value }))}
            required
          />
        </label>
        <label>
          Confirm New Password
          <input
            type="password"
            minLength={8}
            value={form.confirmPassword}
            onChange={(event) => setForm((curr) => ({ ...curr, confirmPassword: event.target.value }))}
            required
          />
        </label>
        {error && <p className="error-text">{error}</p>}
        {message && <p className="success-text">{message}</p>}
        <button type="submit" disabled={loading}>
          {loading ? 'Updating...' : 'Update Password'}
        </button>
      </form>
      <div className="auth-footer">
        <Link to="/login">Back to login</Link>
      </div>
    </section>
  );
}
