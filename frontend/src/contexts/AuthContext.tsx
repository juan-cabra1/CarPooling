import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import type { User } from '@/types';
import { usersService } from '@/services/api';
import type { LoginCredentials, RegisterData } from '@/services/api';
import { authUtils } from '@/services/auth/auth.utils';

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => Promise<void>;
  updateUserData: (user: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
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

  const login = async (credentials: LoginCredentials) => {
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

  const register = async (data: RegisterData) => {
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

  const updateUserData = (updatedUser: User) => {
    setUser(updatedUser);
    authUtils.setUser(updatedUser);
  };

  const value: AuthContextType = {
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
