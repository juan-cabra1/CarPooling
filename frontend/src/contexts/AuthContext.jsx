import React, { createContext, useContext, useState, useEffect } from 'react';
import { usersService } from '@/services/api';
import { authUtils } from '@/services/auth/auth.utils';

const AuthContext = createContext(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  // Initialize auth state from localStorage
  useEffect(() => {
    const initAuth = async () => {
      try {
        const { token: storedToken, user: storedUser } = authUtils.initAuth();

        if (storedToken && storedUser) {
          setToken(storedToken);
          setUser(storedUser);

          // Verify token is still valid by fetching current user
          try {
            const response = await usersService.getCurrentUser();
            setUser(response.data);
            authUtils.setUser(response.data);
          } catch (error) {
            // Token is invalid, clear auth
            authUtils.clearAuth();
            setToken(null);
            setUser(null);
          }
        }
      } catch (error) {
        console.error('Failed to initialize auth:', error);
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  const login = async (credentials) => {
    try {
      const response = await usersService.login(credentials);
      const { user: userData, token: authToken } = response.data;

      authUtils.setToken(authToken);
      authUtils.setUser(userData);

      setToken(authToken);
      setUser(userData);
    } catch (error) {
      throw error;
    }
  };

  const register = async (data) => {
    try {
      // Primero registrar al usuario
      await usersService.register(data);

      // Luego hacer login automático con las credenciales
      const loginResponse = await usersService.login({
        email: data.email,
        password: data.password,
      });

      const { user: userData, token: authToken } = loginResponse.data;

      authUtils.setToken(authToken);
      authUtils.setUser(userData);

      setToken(authToken);
      setUser(userData);
    } catch (error) {
      throw error;
    }
  };

  const logout = async () => {
    try {
      await usersService.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      authUtils.clearAuth();
      setToken(null);
      setUser(null);
    }
  };

  const updateUserData = (updatedUser) => {
    setUser(updatedUser);
    authUtils.setUser(updatedUser);
  };

  const value = {
    user,
    token,
    isAuthenticated: !!token && !!user,
    isLoading,
    login,
    register,
    logout,
    updateUserData,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
