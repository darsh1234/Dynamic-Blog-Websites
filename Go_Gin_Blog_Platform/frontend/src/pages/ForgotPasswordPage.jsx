import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function ForgotPasswordPage() {
  const auth = useAuth();
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const onSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setMessage('');
    setLoading(true);

    try {
      const payload = await auth.requestPasswordReset(email);
      setMessage(payload?.message || 'If the email exists, a password reset link will be sent.');
    } catch (err) {
      setError(err.message || 'Could not request password reset');
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="card auth-card">
      <h1>Reset Password</h1>
      <p>Request a one-time password reset link.</p>
      <form onSubmit={onSubmit}>
        <label>
          Email
          <input
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
        </label>
        {error && <p className="error-text">{error}</p>}
        {message && <p className="success-text">{message}</p>}
        <button type="submit" disabled={loading}>
          {loading ? 'Submitting...' : 'Send Reset Link'}
        </button>
      </form>
      <div className="auth-footer">
        <Link to="/login">Back to login</Link>
      </div>
    </section>
  );
}
