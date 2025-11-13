import { useState, useId } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const id = useId();

  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login({ email, password });
      navigate('/');
    } catch (err) {
      setError('Email o contraseña incorrectos');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 via-blue-50 to-primary-100 py-12 px-4">
      <div className="w-full max-w-[400px]">
        {/* Header con ícono */}
        <div className="flex flex-col items-center gap-4 mb-8 animate-fade-in">
          <div
            className="flex size-14 shrink-0 items-center justify-center rounded-full border-2 border-primary bg-white shadow-lg"
            aria-hidden="true"
          >
            <svg
              className="stroke-primary"
              xmlns="http://www.w3.org/2000/svg"
              width="28"
              height="28"
              viewBox="0 0 32 32"
              aria-hidden="true"
            >
              <circle cx="16" cy="16" r="12" fill="none" strokeWidth="3" />
            </svg>
          </div>
          <div className="text-center">
            <h1 className="text-3xl font-bold tracking-tight text-gray-900">
              Iniciar Sesión
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              Bienvenido de nuevo a CarPooling
            </p>
          </div>
        </div>

        {/* Card con formulario */}
        <div className="rounded-xl border bg-card p-6 shadow-lg">
          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="bg-destructive/10 border border-destructive/20 text-destructive px-4 py-3 rounded-lg animate-shake" role="alert">
                <p className="font-medium text-sm">Error</p>
                <p className="text-sm">{error}</p>
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor={`${id}-email`}>Correo electrónico</Label>
              <Input
                id={`${id}-email`}
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="tu@email.com"
                autoComplete="email"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor={`${id}-password`}>Contraseña</Label>
              <Input
                id={`${id}-password`}
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                placeholder="Ingresa tu contraseña"
                autoComplete="current-password"
              />
            </div>

            <div className="flex items-center justify-end">
              <Link
                to="/forgot-password"
                className="text-sm text-primary hover:underline font-medium"
              >
                ¿Olvidaste tu contraseña?
              </Link>
            </div>

            <Button
              type="submit"
              className="w-full"
              disabled={isLoading}
            >
              {isLoading ? 'Iniciando sesión...' : 'Iniciar sesión'}
            </Button>
          </form>

          {/* Separador "Or" */}
          <div className="flex items-center gap-3 my-5 before:h-px before:flex-1 before:bg-border after:h-px after:flex-1 after:bg-border">
            <span className="text-xs text-muted-foreground">O</span>
          </div>

          {/* Botón Google (placeholder) */}
          <Button variant="outline" className="w-full" type="button">
            <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
              <path
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                fill="#4285F4"
              />
              <path
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                fill="#34A853"
              />
              <path
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                fill="#FBBC05"
              />
              <path
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                fill="#EA4335"
              />
            </svg>
          </Button>

          {/* Link a registro */}
          <p className="text-center text-sm text-muted-foreground mt-5">
            ¿No tienes una cuenta?{' '}
            <Link to="/register" className="text-primary font-semibold hover:underline">
              Regístrate gratis
            </Link>
          </p>
        </div>

        {/* Términos */}
        <p className="text-center text-xs text-muted-foreground mt-6">
          Al iniciar sesión, aceptas nuestros{' '}
          <Link to="/terms" className="underline hover:no-underline">
            Términos
          </Link>
          {' y '}
          <Link to="/privacy" className="underline hover:no-underline">
            Privacidad
          </Link>
          .
        </p>
      </div>
    </div>
  );
};
