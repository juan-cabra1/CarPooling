import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Input, Card } from '@/components/common';

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login({ email, password });
      navigate('/');
    } catch (err: unknown) {
      setError('Email o contraseña incorrectos');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 via-blue-50 to-primary-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <div className="text-center mb-8 animate-fade-in">
          <div className="inline-block p-4 bg-white rounded-2xl shadow-lg mb-4">
            <span className="text-6xl">🚗</span>
          </div>
          <h1 className="text-5xl font-bold text-primary-600 mb-3">CarPooling</h1>
          <h2 className="text-2xl font-semibold text-gray-900 mb-2">Bienvenido de nuevo</h2>
          <p className="text-gray-600">Inicia sesión para continuar</p>
        </div>

        <Card className="shadow-2xl backdrop-blur-sm bg-white/95">
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="bg-red-50 border-l-4 border-red-500 text-red-700 px-4 py-3 rounded-r-lg animate-shake" role="alert">
                <p className="font-medium">Error</p>
                <p className="text-sm">{error}</p>
              </div>
            )}

            <Input
              label="Correo electrónico"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              fullWidth
              placeholder="tu@email.com"
              className="transition-all focus:scale-105"
            />

            <Input
              label="Contraseña"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              fullWidth
              placeholder="Ingresa tu contraseña"
              className="transition-all focus:scale-105"
            />

            <div className="flex items-center justify-between">
              <Link to="/forgot-password" className="text-sm text-primary-600 hover:text-primary-700 hover:underline font-medium">
                ¿Olvidaste tu contraseña?
              </Link>
            </div>

            <Button
              type="submit"
              variant="primary"
              fullWidth
              isLoading={isLoading}
              className="transform hover:scale-105 transition-all shadow-lg hover:shadow-xl"
            >
              Iniciar sesión
            </Button>

            <div className="text-center text-sm pt-4 border-t border-gray-200">
              <span className="text-gray-600">¿No tienes una cuenta? </span>
              <Link to="/register" className="text-primary-600 hover:text-primary-700 font-semibold hover:underline">
                Regístrate gratis
              </Link>
            </div>
          </form>
        </Card>

        <div className="mt-8 text-center">
          <p className="text-sm text-gray-500">
            Al iniciar sesión, aceptas nuestros{' '}
            <Link to="/terms" className="text-primary-600 hover:underline">
              Términos de Servicio
            </Link>{' '}
            y{' '}
            <Link to="/privacy" className="text-primary-600 hover:underline">
              Política de Privacidad
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};
