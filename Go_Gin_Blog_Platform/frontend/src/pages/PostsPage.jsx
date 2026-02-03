import { useCallback, useEffect, useMemo, useState } from 'react';
import { createApiClient } from '../lib/client';
import { useAuth } from '../context/AuthContext';

const DEFAULT_DRAFT = {
  title: '',
  content: '',
  status: 'published'
};

export default function PostsPage() {
  const auth = useAuth();
  const client = useMemo(() => createApiClient(auth), [auth]);

  const [posts, setPosts] = useState([]);
  const [meta, setMeta] = useState({ page: 1, limit: 10, total: 0, total_pages: 1 });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [draft, setDraft] = useState(DEFAULT_DRAFT);
  const [editingPostId, setEditingPostId] = useState('');
  const [editingDraft, setEditingDraft] = useState(DEFAULT_DRAFT);

  const canWrite = auth.user?.role === 'author' || auth.user?.role === 'admin';

  const loadPosts = useCallback(async (page = 1, limit = 10) => {
    setLoading(true);
    setError('');
    try {
      const payload = await client.listPosts(page, limit);
      setPosts(payload?.data || []);
      setMeta(
        payload?.meta || {
          page,
          limit,
          total: 0,
          total_pages: 1
        }
      );
    } catch (err) {
      setError(err.message || 'Failed to load posts');
    } finally {
      setLoading(false);
    }
  }, [client]);

  useEffect(() => {
    loadPosts(1, 10);
  }, [loadPosts]);

  const onCreatePost = async (event) => {
    event.preventDefault();
    setCreating(true);
    setError('');

    try {
      await client.createPost(draft);
      setDraft(DEFAULT_DRAFT);
      await loadPosts(1, meta.limit);
    } catch (err) {
      setError(err.message || 'Failed to create post');
    } finally {
      setCreating(false);
    }
  };

  const startEditing = (post) => {
    setEditingPostId(post.id);
    setEditingDraft({
      title: post.title,
      content: post.content,
      status: post.status
    });
  };

  const cancelEditing = () => {
    setEditingPostId('');
    setEditingDraft(DEFAULT_DRAFT);
  };

  const onUpdatePost = async (event) => {
    event.preventDefault();
    if (!editingPostId) {
      return;
    }

    setError('');
    try {
      await client.updatePost(editingPostId, editingDraft);
      cancelEditing();
      await loadPosts(meta.page, meta.limit);
    } catch (err) {
      setError(err.message || 'Failed to update post');
    }
  };

  const onDeletePost = async (postId) => {
    if (!window.confirm('Delete this post?')) {
      return;
    }

    setError('');
    try {
      await client.deletePost(postId);
      await loadPosts(meta.page, meta.limit);
    } catch (err) {
      setError(err.message || 'Failed to delete post');
    }
  };

  const nextPage = () => {
    if (meta.page < meta.total_pages) {
      loadPosts(meta.page + 1, meta.limit);
    }
  };

  const prevPage = () => {
    if (meta.page > 1) {
      loadPosts(meta.page - 1, meta.limit);
    }
  };

  return (
    <section className="stack">
      <div className="section-title">
        <h1>Posts</h1>
        <p>
          Browse posts with pagination. Signed-in <strong>{auth.user.role}</strong> users can act on permitted
          actions.
        </p>
      </div>

      {canWrite && (
        <article className="card">
          <h2>Create New Post</h2>
          <form className="stack-sm" onSubmit={onCreatePost}>
            <label>
              Title
              <input
                type="text"
                value={draft.title}
                onChange={(event) => setDraft((curr) => ({ ...curr, title: event.target.value }))}
                required
              />
            </label>
            <label>
              Content
              <textarea
                rows={5}
                value={draft.content}
                onChange={(event) => setDraft((curr) => ({ ...curr, content: event.target.value }))}
                required
              />
            </label>
            <label>
              Status
              <select
                value={draft.status}
                onChange={(event) => setDraft((curr) => ({ ...curr, status: event.target.value }))}
              >
                <option value="published">published</option>
                <option value="draft">draft</option>
              </select>
            </label>
            <button type="submit" disabled={creating}>
              {creating ? 'Creating...' : 'Create Post'}
            </button>
          </form>
        </article>
      )}

      {error && <p className="error-text">{error}</p>}

      <article className="card">
        <div className="row-between">
          <h2>Latest Posts</h2>
          <span>
            Page {meta.page} of {Math.max(meta.total_pages || 1, 1)}
          </span>
        </div>

        {loading ? (
          <p>Loading posts...</p>
        ) : posts.length === 0 ? (
          <p>No posts found.</p>
        ) : (
          <div className="post-list">
            {posts.map((post) => {
              const isOwnerOrAdmin = auth.user?.id === post.author_id || auth.user?.role === 'admin';
              const isEditing = editingPostId === post.id;

              return (
                <article key={post.id} className="post-card">
                  {isEditing ? (
                    <form className="stack-sm" onSubmit={onUpdatePost}>
                      <label>
                        Title
                        <input
                          type="text"
                          value={editingDraft.title}
                          onChange={(event) =>
                            setEditingDraft((curr) => ({ ...curr, title: event.target.value }))
                          }
                          required
                        />
                      </label>
                      <label>
                        Content
                        <textarea
                          rows={4}
                          value={editingDraft.content}
                          onChange={(event) =>
                            setEditingDraft((curr) => ({ ...curr, content: event.target.value }))
                          }
                          required
                        />
                      </label>
                      <label>
                        Status
                        <select
                          value={editingDraft.status}
                          onChange={(event) =>
                            setEditingDraft((curr) => ({ ...curr, status: event.target.value }))
                          }
                        >
                          <option value="published">published</option>
                          <option value="draft">draft</option>
                        </select>
                      </label>
                      <div className="row-actions">
                        <button type="submit">Save</button>
                        <button type="button" className="ghost" onClick={cancelEditing}>
                          Cancel
                        </button>
                      </div>
                    </form>
                  ) : (
                    <>
                      <h3>{post.title}</h3>
                      <p className="post-meta">
                        Author {post.author_id} · {post.status}
                      </p>
                      <p>{post.content}</p>
                      {isOwnerOrAdmin && canWrite && (
                        <div className="row-actions">
                          <button type="button" className="ghost" onClick={() => startEditing(post)}>
                            Edit
                          </button>
                          <button type="button" className="danger" onClick={() => onDeletePost(post.id)}>
                            Delete
                          </button>
                        </div>
                      )}
                    </>
                  )}
                </article>
              );
            })}
          </div>
        )}

        <div className="row-actions">
          <button type="button" className="ghost" onClick={prevPage} disabled={meta.page <= 1 || loading}>
            Previous
          </button>
          <button
            type="button"
            className="ghost"
            onClick={nextPage}
            disabled={meta.page >= meta.total_pages || loading}
          >
            Next
          </button>
        </div>
      </article>
    </section>
  );
}
