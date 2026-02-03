import { useCallback, useEffect, useMemo, useState } from 'react';
import { createApiClient } from '../lib/client';
import { useAuth } from '../context/AuthContext';

const ROLE_OPTIONS = ['reader', 'author', 'admin'];

export default function AdminUsersPage() {
  const auth = useAuth();
  const client = useMemo(() => createApiClient(auth), [auth]);

  const [users, setUsers] = useState([]);
  const [meta, setMeta] = useState({ page: 1, limit: 10, total_pages: 1, total: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadUsers = useCallback(async (page = 1, limit = 10) => {
    setLoading(true);
    setError('');
    try {
      const payload = await client.listUsers(page, limit);
      setUsers(payload?.data || []);
      setMeta(payload?.meta || { page, limit, total_pages: 1, total: 0 });
    } catch (err) {
      setError(err.message || 'Failed to load users');
    } finally {
      setLoading(false);
    }
  }, [client]);

  useEffect(() => {
    loadUsers(1, 10);
  }, [loadUsers]);

  const onChangeRole = async (userId, role) => {
    setError('');
    try {
      await client.updateUserRole(userId, role);
      await loadUsers(meta.page, meta.limit);
    } catch (err) {
      setError(err.message || 'Failed to update user role');
    }
  };

  return (
    <section className="stack">
      <div className="section-title">
        <h1>Admin: Users</h1>
        <p>Manage user roles for reader, author, and admin access.</p>
      </div>

      <article className="card">
        {error && <p className="error-text">{error}</p>}
        {loading ? (
          <p>Loading users...</p>
        ) : (
          <>
            <table className="table">
              <thead>
                <tr>
                  <th>Email</th>
                  <th>User ID</th>
                  <th>Role</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id}>
                    <td>{user.email}</td>
                    <td className="mono">{user.id}</td>
                    <td>
                      <select value={user.role} onChange={(event) => onChangeRole(user.id, event.target.value)}>
                        {ROLE_OPTIONS.map((role) => (
                          <option key={role} value={role}>
                            {role}
                          </option>
                        ))}
                      </select>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div className="row-actions">
              <button type="button" className="ghost" onClick={() => loadUsers(meta.page - 1, meta.limit)} disabled={meta.page <= 1}>
                Previous
              </button>
              <span>
                Page {meta.page} of {Math.max(meta.total_pages || 1, 1)}
              </span>
              <button
                type="button"
                className="ghost"
                onClick={() => loadUsers(meta.page + 1, meta.limit)}
                disabled={meta.page >= meta.total_pages}
              >
                Next
              </button>
            </div>
          </>
        )}
      </article>
    </section>
  );
}
