import { Link } from 'react-router-dom';
import { ShieldAlert } from 'lucide-react';
import { Button } from '@/components/ui/button';

export default function UnauthorizedPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
      <div className="max-w-md w-full text-center">
        <div className="flex justify-center mb-6">
          <div className="rounded-full bg-destructive-100 dark:bg-destructive-900/20 p-6">
            <ShieldAlert className="h-16 w-16 text-destructive-600 dark:text-destructive-400" />
          </div>
        </div>

        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-3">
          Acceso Denegado
        </h1>

        <p className="text-gray-600 dark:text-gray-400 mb-8">
          No tienes permisos para acceder a esta sección. Solo los administradores pueden ver este contenido.
        </p>

        <div className="space-y-3">
          <Link to="/" className="block">
            <Button className="w-full">
              Volver al Inicio
            </Button>
          </Link>

          <Link to="/profile" className="block">
            <Button variant="outline" className="w-full">
              Ir a Mi Perfil
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
