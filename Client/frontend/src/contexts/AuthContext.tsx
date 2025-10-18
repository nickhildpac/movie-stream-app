import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { jwtDecode } from 'jwt-decode';
import type { User, LoginInput, RegisterInput } from '../types';

interface AuthContextType {
  user: User | null;
  login: (input: LoginInput) => Promise<void>;
  register: (input: RegisterInput) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
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
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const refreshToken = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/refresh', {
          method: 'POST',
          credentials: 'include',
        });

        if (response.ok) {
          const data = await response.json();
          const decodedToken: { id: string; name: string; email: string } = jwtDecode(data.access_token);
          const user: User = {
            id: decodedToken.id,
            name: decodedToken.name,
            email: decodedToken.email,
          };
          setUser(user);
        } 
      } catch (error) {
        console.error('Could not refresh token:', error);
      } finally {
        setIsLoading(false);
      }
    };

    refreshToken();
  }, []);

  const login = async (input: LoginInput) => {
    const response = await fetch('http://localhost:8080/v1/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(input),
      credentials: 'include',
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Login failed');
    }

    const data = await response.json();
    const decodedToken: { id: string; name: string; email: string } = jwtDecode(data.access_token);
    const user: User = {
      id: decodedToken.id,
      name: decodedToken.name,
      email: decodedToken.email,
    };
    setUser(user);
  };

  const register = async (input: RegisterInput) => {
    if (input.password !== input.confirmPassword) {
      throw new Error('Passwords do not match');
    }

    const response = await fetch('http://localhost:8080/v1/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        first_name: input.first_name,
        last_name: input.last_name,
        email: input.email,
        password: input.password,
        favourite_genres: input.favourite_genres,
        role: "USER",
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Registration failed');
    }
  };

  const logout = async () => {
    await fetch('http://localhost:8080/v1/logout', {
      method: 'POST',
      body: JSON.stringify({ user_id: user?.id }),
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });

    setUser(null);
    localStorage.removeItem('user');
  };

  return (
    <AuthContext.Provider value={{ user, login, register, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};