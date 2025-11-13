import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Input, Card } from '@/components/common';

export const RegisterPage: React.FC = () => {
  const [formData, setFormData] = useState({
    name: '',
    lastname: '',
    email: '',
    phone: '',
    street: '',
    number: '',
    sex: 'otro' as 'hombre' | 'mujer' | 'otro',
    birthdate: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const { register } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
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
      console.log('Datos a enviar:', {
        name: formData.name,
        lastname: formData.lastname,
        email: formData.email,
        phone: formData.phone,
        street: formData.street,
        number: Number(formData.number),
        sex: formData.sex,
        birthdate: formData.birthdate,
        password: '***',
      });

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
    } catch (err: any) {
      console.error('Error completo:', err);

      // Intentar extraer mensaje de error del backend
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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl w-full">
        <div className="text-center mb-8">
          <h1 className="text-5xl font-bold text-primary-600 mb-2">🚗 CarPooling</h1>
          <h2 className="text-3xl font-semibold text-gray-900 mb-2">Crea tu cuenta</h2>
          <p className="text-gray-600">Únete a nuestra comunidad de viajes compartidos</p>
        </div>

        <Card className="shadow-2xl">
          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="bg-red-50 border-l-4 border-red-500 text-red-700 px-4 py-3 rounded-r-lg" role="alert">
                <p className="font-medium">Error</p>
                <p className="text-sm">{error}</p>
              </div>
            )}

            <div className="grid md:grid-cols-2 gap-4">
              <Input
                label="Nombre"
                name="name"
                type="text"
                value={formData.name}
                onChange={handleChange}
                required
                fullWidth
                placeholder="Juan"
                className="transition-all focus:scale-105"
              />

              <Input
                label="Apellido"
                name="lastname"
                type="text"
                value={formData.lastname}
                onChange={handleChange}
                required
                fullWidth
                placeholder="Pérez"
                className="transition-all focus:scale-105"
              />
            </div>

            <Input
              label="Correo electrónico"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              required
              fullWidth
              placeholder="tu@email.com"
              className="transition-all focus:scale-105"
            />

            <div className="grid md:grid-cols-2 gap-4">
              <Input
                label="Teléfono"
                name="phone"
                type="tel"
                value={formData.phone}
                onChange={handleChange}
                required
                fullWidth
                placeholder="+54 9 11 1234-5678"
                className="transition-all focus:scale-105"
              />

              <div className="w-full">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Sexo
                </label>
                <select
                  name="sex"
                  value={formData.sex}
                  onChange={handleChange}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 transition-all"
                >
                  <option value="hombre">Masculino</option>
                  <option value="mujer">Femenino</option>
                  <option value="otro">Otro</option>
                </select>
              </div>
            </div>

            <div className="grid md:grid-cols-3 gap-4">
              <div className="md:col-span-2">
                <Input
                  label="Calle"
                  name="street"
                  type="text"
                  value={formData.street}
                  onChange={handleChange}
                  required
                  fullWidth
                  placeholder="Av. Siempre Viva"
                  className="transition-all focus:scale-105"
                />
              </div>

              <Input
                label="Número"
                name="number"
                type="number"
                value={formData.number}
                onChange={handleChange}
                required
                fullWidth
                placeholder="123"
                className="transition-all focus:scale-105"
              />
            </div>

            <Input
              label="Fecha de nacimiento"
              name="birthdate"
              type="date"
              value={formData.birthdate}
              onChange={handleChange}
              required
              fullWidth
              className="transition-all focus:scale-105"
            />

            <div className="grid md:grid-cols-2 gap-4">
              <Input
                label="Contraseña"
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                required
                fullWidth
                placeholder="Mínimo 8 caracteres"
                helperText="Debe tener al menos 8 caracteres"
                className="transition-all focus:scale-105"
              />

              <Input
                label="Confirmar contraseña"
                name="confirmPassword"
                type="password"
                value={formData.confirmPassword}
                onChange={handleChange}
                required
                fullWidth
                placeholder="Confirma tu contraseña"
                className="transition-all focus:scale-105"
              />
            </div>

            <Button
              type="submit"
              variant="primary"
              fullWidth
              isLoading={isLoading}
              className="mt-6 transform hover:scale-105 transition-transform"
            >
              Crear cuenta
            </Button>

            <div className="text-center text-sm pt-4 border-t border-gray-200">
              <span className="text-gray-600">¿Ya tienes una cuenta? </span>
              <Link to="/login" className="text-primary-600 hover:text-primary-700 font-semibold hover:underline">
                Inicia sesión
              </Link>
            </div>
          </form>
        </Card>
      </div>
    </div>
  );
};
