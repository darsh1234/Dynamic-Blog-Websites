const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

function buildUrl(path) {
  return `${API_BASE_URL}${path}`;
}

async function safeJson(response) {
  const text = await response.text();
  if (!text) {
    return null;
  }

  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

export function parseError(payload, fallbackMessage) {
  if (payload && payload.error && payload.error.message) {
    return payload.error.message;
  }
  return fallbackMessage;
}

async function request(path, options = {}) {
  const response = await fetch(buildUrl(path), {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {})
    },
    body: options.body ? JSON.stringify(options.body) : undefined
  });

  const payload = await safeJson(response);
  return { response, payload };
}

export async function requestWithAuth(path, options) {
  const {
    accessToken,
    refreshToken,
    onTokenRefresh,
    onUnauthorized,
    retry = true,
    ...requestOptions
  } = options || {};

  const headers = {
    ...(requestOptions.headers || {})
  };

  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  const first = await request(path, { ...requestOptions, headers });

  if (first.response.status !== 401 || !refreshToken || !retry) {
    return first;
  }

  const refreshResult = await request('/auth/refresh', {
    method: 'POST',
    body: { refresh_token: refreshToken }
  });

  if (!refreshResult.response.ok || !refreshResult.payload?.tokens) {
    onUnauthorized?.();
    return first;
  }

  const newTokens = {
    accessToken: refreshResult.payload.tokens.access_token,
    refreshToken: refreshResult.payload.tokens.refresh_token,
    tokenType: refreshResult.payload.tokens.token_type
  };

  onTokenRefresh?.(newTokens);

  const retryHeaders = {
    ...(requestOptions.headers || {}),
    Authorization: `Bearer ${newTokens.accessToken}`
  };

  return request(path, {
    ...requestOptions,
    headers: retryHeaders
  });
}

export const api = {
  register: (email, password) =>
    request('/auth/register', {
      method: 'POST',
      body: { email, password }
    }),

  login: (email, password) =>
    request('/auth/login', {
      method: 'POST',
      body: { email, password }
    }),

  logout: (refreshToken) =>
    request('/auth/logout', {
      method: 'POST',
      body: { refresh_token: refreshToken }
    }),

  requestPasswordReset: (email) =>
    request('/auth/password-reset/request', {
      method: 'POST',
      body: { email }
    }),

  confirmPasswordReset: (token, newPassword) =>
    request('/auth/password-reset/confirm', {
      method: 'POST',
      body: { token, new_password: newPassword }
    })
};

export { API_BASE_URL };
