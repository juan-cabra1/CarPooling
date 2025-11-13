import { useState, useId } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export const RegisterPage = () => {
  const [formData, setFormData] = useState({
    name: '',
    lastname: '',
    email: '',
    phone: '',
    street: '',
    number: '',
    sex: 'otro',
    birthdate: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const id = useId();

  const { register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirmPassword) {
      setError('Las contraseñas no coinciden');
      return;
    }

    if (formData.password.length < 8) {
      setError('La contraseña debe tener al menos 8 caracteres');
      return;
    }

    if (!formData.number || isNaN(Number(formData.number))) {
      setError('El número de calle debe ser válido');
      return;
    }

    setIsLoading(true);

    try {
      await register({
        name: formData.name,
        lastname: formData.lastname,
        email: formData.email,
        phone: formData.phone,
        street: formData.street,
        number: Number(formData.number),
        sex: formData.sex,
        birthdate: formData.birthdate,
        password: formData.password,
      });
      navigate('/');
    } catch (err) {
      const errorMessage = err?.response?.data?.error ||
                          err?.response?.data?.message ||
                          err?.message ||
                          'Error al registrarse. Por favor, intenta de nuevo.';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 py-12 px-4">
      <div className="w-full max-w-2xl">
        {/* Header */}
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
              Crear Cuenta
            </h1>
            <p className="mt-2 text-sm text-muted-foreground">
              Únete a nuestra comunidad de viajes compartidos
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

            {/* Nombre y Apellido */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor={`${id}-name`}>Nombre</Label>
                <Input
                  id={`${id}-name`}
                  name="name"
                  type="text"
                  value={formData.name}
                  onChange={handleChange}
                  required
                  placeholder="Juan"
                  autoComplete="given-name"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${id}-lastname`}>Apellido</Label>
                <Input
                  id={`${id}-lastname`}
                  name="lastname"
                  type="text"
                  value={formData.lastname}
                  onChange={handleChange}
                  required
                  placeholder="Pérez"
                  autoComplete="family-name"
                />
              </div>
            </div>

            {/* Email */}
            <div className="space-y-2">
              <Label htmlFor={`${id}-email`}>Correo electrónico</Label>
              <Input
                id={`${id}-email`}
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                required
                placeholder="tu@email.com"
                autoComplete="email"
              />
            </div>

            {/* Teléfono y Sexo */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor={`${id}-phone`}>Teléfono</Label>
                <Input
                  id={`${id}-phone`}
                  name="phone"
                  type="tel"
                  value={formData.phone}
                  onChange={handleChange}
                  required
                  placeholder="+54 9 11 1234-5678"
                  autoComplete="tel"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${id}-sex`}>Sexo</Label>
                <select
                  id={`${id}-sex`}
                  name="sex"
                  value={formData.sex}
                  onChange={handleChange}
                  required
                  className="flex h-9 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground shadow-sm shadow-black/5 transition-shadow focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20"
                >
                  <option value="hombre">Masculino</option>
                  <option value="mujer">Femenino</option>
                  <option value="otro">Otro</option>
                </select>
              </div>
            </div>

            {/* Calle y Número */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="md:col-span-2 space-y-2">
                <Label htmlFor={`${id}-street`}>Calle</Label>
                <Input
                  id={`${id}-street`}
                  name="street"
                  type="text"
                  value={formData.street}
                  onChange={handleChange}
                  required
                  placeholder="Av. Siempre Viva"
                  autoComplete="street-address"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${id}-number`}>Número</Label>
                <Input
                  id={`${id}-number`}
                  name="number"
                  type="number"
                  value={formData.number}
                  onChange={handleChange}
                  required
                  placeholder="123"
                />
              </div>
            </div>

            {/* Fecha de nacimiento */}
            <div className="space-y-2">
              <Label htmlFor={`${id}-birthdate`}>Fecha de nacimiento</Label>
              <Input
                id={`${id}-birthdate`}
                name="birthdate"
                type="date"
                value={formData.birthdate}
                onChange={handleChange}
                required
                autoComplete="bday"
              />
            </div>

            {/* Contraseñas */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor={`${id}-password`}>Contraseña</Label>
                <Input
                  id={`${id}-password`}
                  name="password"
                  type="password"
                  value={formData.password}
                  onChange={handleChange}
                  required
                  placeholder="Mínimo 8 caracteres"
                  autoComplete="new-password"
                />
                <p className="text-xs text-muted-foreground">
                  Debe tener al menos 8 caracteres
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${id}-confirmPassword`}>Confirmar contraseña</Label>
                <Input
                  id={`${id}-confirmPassword`}
                  name="confirmPassword"
                  type="password"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                  required
                  placeholder="Confirma tu contraseña"
                  autoComplete="new-password"
                />
              </div>
            </div>

            {/* Botón submit */}
            <Button
              type="submit"
              className="w-full mt-2"
              disabled={isLoading}
            >
              {isLoading ? 'Creando cuenta...' : 'Crear cuenta'}
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
            Continuar con Google
          </Button>

          {/* Link a login */}
          <p className="text-center text-sm text-muted-foreground mt-5">
            ¿Ya tienes una cuenta?{' '}
            <Link to="/login" className="text-primary font-semibold hover:underline">
              Inicia sesión
            </Link>
          </p>
        </div>

        {/* Términos */}
        <p className="text-center text-xs text-muted-foreground mt-6">
          Al registrarte, aceptas nuestros{' '}
          <Link to="/terms" className="underline hover:no-underline">
            Términos
          </Link>
          .
        </p>
      </div>
    </div>
  );
};
