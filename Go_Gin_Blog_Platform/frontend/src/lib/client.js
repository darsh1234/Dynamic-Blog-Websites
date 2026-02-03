import { parseError, requestWithAuth } from './api';

function queryString(params = {}) {
  const search = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === '') {
      return;
    }
    search.set(key, String(value));
  });
  const raw = search.toString();
  return raw ? `?${raw}` : '';
}

export function createApiClient(auth) {
  async function doAuthRequest(path, options = {}) {
    const { response, payload } = await requestWithAuth(path, {
      ...options,
      accessToken: auth.tokens?.accessToken,
      refreshToken: auth.tokens?.refreshToken,
      onTokenRefresh: auth.updateTokens,
      onUnauthorized: auth.logout
    });

    if (!response.ok) {
      throw new Error(parseError(payload, 'API request failed'));
    }

    return payload;
  }

  return {
    async listPosts(page = 1, limit = 10) {
      return doAuthRequest(`/posts${queryString({ page, limit })}`, {
        method: 'GET',
        retry: false
      });
    },

    async getPost(postId) {
      return doAuthRequest(`/posts/${postId}`, {
        method: 'GET',
        retry: false
      });
    },

    async createPost(input) {
      return doAuthRequest('/posts', {
        method: 'POST',
        body: input
      });
    },

    async updatePost(postId, input) {
      return doAuthRequest(`/posts/${postId}`, {
        method: 'PATCH',
        body: input
      });
    },

    async deletePost(postId) {
      return doAuthRequest(`/posts/${postId}`, {
        method: 'DELETE'
      });
    },

    async listUsers(page = 1, limit = 10) {
      return doAuthRequest(`/admin/users${queryString({ page, limit })}`, {
        method: 'GET'
      });
    },

    async updateUserRole(userId, role) {
      return doAuthRequest(`/admin/users/${userId}/role`, {
        method: 'PATCH',
        body: { role }
      });
    }
  };
}
