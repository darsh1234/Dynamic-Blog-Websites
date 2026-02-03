import { createContext, useContext, useMemo, useState } from 'react';
import { api, parseError } from '../lib/api';

const STORAGE_KEY = 'go_gin_blog_auth';

function loadStoredAuth() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return { user: null, tokens: null };
    }
    const parsed = JSON.parse(raw);
    return {
      user: parsed.user || null,
      tokens: parsed.tokens || null
    };
  } catch {
    return { user: null, tokens: null };
  }
}

function persistAuth(user, tokens) {
  if (!user || !tokens) {
    localStorage.removeItem(STORAGE_KEY);
    return;
  }
  localStorage.setItem(
    STORAGE_KEY,
    JSON.stringify({
      user,
      tokens
    })
  );
}

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const stored = loadStoredAuth();
  const [user, setUser] = useState(stored.user);
  const [tokens, setTokens] = useState(stored.tokens);

  const setAuthState = (nextUser, nextTokens) => {
    setUser(nextUser);
    setTokens(nextTokens);
    persistAuth(nextUser, nextTokens);
  };

  const register = async (email, password) => {
    const { response, payload } = await api.register(email, password);
    if (!response.ok) {
      throw new Error(parseError(payload, 'Registration failed'));
    }

    const nextUser = payload?.user || null;
    const nextTokens = {
      accessToken: payload?.tokens?.access_token,
      refreshToken: payload?.tokens?.refresh_token,
      tokenType: payload?.tokens?.token_type || 'Bearer'
    };

    setAuthState(nextUser, nextTokens);
    return nextUser;
  };

  const login = async (email, password) => {
    const { response, payload } = await api.login(email, password);
    if (!response.ok) {
      throw new Error(parseError(payload, 'Login failed'));
    }

    const nextUser = payload?.user || null;
    const nextTokens = {
      accessToken: payload?.tokens?.access_token,
      refreshToken: payload?.tokens?.refresh_token,
      tokenType: payload?.tokens?.token_type || 'Bearer'
    };

    setAuthState(nextUser, nextTokens);
    return nextUser;
  };

  const logout = async () => {
    if (tokens?.refreshToken) {
      try {
        await api.logout(tokens.refreshToken);
      } catch {
        // Ignore network errors and clear local state anyway.
      }
    }
    setAuthState(null, null);
  };

  const updateTokens = (nextTokens) => {
    if (!user) {
      return;
    }
    setAuthState(user, nextTokens);
  };

  const value = useMemo(
    () => ({
      user,
      tokens,
      isAuthenticated: Boolean(user && tokens?.accessToken),
      register,
      login,
      logout,
      updateTokens,
      requestPasswordReset: async (email) => {
        const { response, payload } = await api.requestPasswordReset(email);
        if (!response.ok) {
          throw new Error(parseError(payload, 'Password reset request failed'));
        }
        return payload;
      },
      confirmPasswordReset: async (token, newPassword) => {
        const { response, payload } = await api.confirmPasswordReset(token, newPassword);
        if (!response.ok) {
          throw new Error(parseError(payload, 'Password reset failed'));
        }
        return payload;
      }
    }),
    [user, tokens]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider');
  }
  return context;
}
