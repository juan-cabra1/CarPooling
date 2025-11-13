import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/common';

export const Navbar = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [profileMenuOpen, setProfileMenuOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    setProfileMenuOpen(false);
    setMobileMenuOpen(false);
    navigate('/login');
  };

  const closeMobileMenu = () => setMobileMenuOpen(false);

  return (
    <nav className="bg-white shadow-lg sticky top-0 z-50">
      <div className="container mx-auto px-4 lg:px-6">
        <div className="flex justify-between items-center h-16 lg:h-18">
          {/* Logo */}
          <div className="flex items-center">
            <Link to="/" className="text-2xl lg:text-3xl font-bold text-primary-600 hover:text-primary-700 transition-colors" onClick={closeMobileMenu}>
              CarPooling
            </Link>
          </div>

          {/* Desktop Navigation */}
          {isAuthenticated && (
            <div className="hidden lg:flex items-center space-x-6">
              <Link to="/trips" className="text-gray-700 hover:text-primary-600 font-medium transition-colors py-2 px-3 rounded-lg hover:bg-primary-50">
                Buscar Viajes
              </Link>
              <Link to="/my-trips" className="text-gray-700 hover:text-primary-600 font-medium transition-colors py-2 px-3 rounded-lg hover:bg-primary-50">
                Mis Viajes
              </Link>
              <Link to="/my-bookings" className="text-gray-700 hover:text-primary-600 font-medium transition-colors py-2 px-3 rounded-lg hover:bg-primary-50">
                Mis Reservas
              </Link>
            </div>
          )}

          {/* Desktop Auth Actions */}
          <div className="hidden md:flex items-center space-x-3">
            {isAuthenticated ? (
              <>
                <Link to="/trips/create">
                  <Button variant="primary" size="sm">
                    Ofrecer Viaje
                  </Button>
                </Link>
                <div className="relative">
                  <button
                    onClick={() => setProfileMenuOpen(!profileMenuOpen)}
                    className="flex items-center space-x-2 text-gray-700 hover:text-primary-600 transition-colors py-2 px-3 rounded-lg hover:bg-gray-50"
                  >
                    <div className="w-10 h-10 bg-gradient-to-br from-primary-600 to-primary-700 rounded-full flex items-center justify-center text-white font-semibold shadow-md">
                      {user?.name?.charAt(0).toUpperCase() || 'U'}{user?.lastname?.charAt(0).toUpperCase() || ''}
                    </div>
                    <span className="hidden lg:block font-medium">{user?.name}</span>
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                  </button>
                  {profileMenuOpen && (
                    <div className="absolute right-0 mt-2 w-52 bg-white rounded-xl shadow-xl py-2 border border-gray-100 animate-fade-in">
                      <div className="px-4 py-3 border-b border-gray-100">
                        <p className="text-sm font-semibold text-gray-900">{user?.name} {user?.lastname}</p>
                        <p className="text-xs text-gray-500 truncate">{user?.email}</p>
                      </div>
                      <Link
                        to="/profile"
                        onClick={() => setProfileMenuOpen(false)}
                        className="block px-4 py-2 text-sm text-gray-700 hover:bg-primary-50 hover:text-primary-600 transition-colors"
                      >
                        Mi Perfil
                      </Link>
                      <button
                        onClick={handleLogout}
                        className="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                      >
                        Cerrar Sesión
                      </button>
                    </div>
                  )}
                </div>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="outline" size="sm">
                    Iniciar Sesión
                  </Button>
                </Link>
                <Link to="/register">
                  <Button variant="primary" size="sm">
                    Registrarse
                  </Button>
                </Link>
              </>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden">
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="text-gray-700 hover:text-primary-600 focus:outline-none p-2"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                {mobileMenuOpen ? (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                ) : (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                )}
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <div className="md:hidden border-t border-gray-100 py-4 animate-slide-up">
            {isAuthenticated ? (
              <>
                <div className="flex items-center space-x-3 px-4 py-3 bg-gray-50 rounded-lg mb-3">
                  <div className="w-12 h-12 bg-gradient-to-br from-primary-600 to-primary-700 rounded-full flex items-center justify-center text-white font-semibold shadow-md">
                    {user?.name?.charAt(0).toUpperCase() || 'U'}{user?.lastname?.charAt(0).toUpperCase() || ''}
                  </div>
                  <div>
                    <p className="font-semibold text-gray-900">{user?.name} {user?.lastname}</p>
                    <p className="text-xs text-gray-500">{user?.email}</p>
                  </div>
                </div>
                <Link to="/trips" onClick={closeMobileMenu} className="block px-4 py-3 text-gray-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors font-medium">
                  Buscar Viajes
                </Link>
                <Link to="/my-trips" onClick={closeMobileMenu} className="block px-4 py-3 text-gray-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors font-medium">
                  Mis Viajes
                </Link>
                <Link to="/my-bookings" onClick={closeMobileMenu} className="block px-4 py-3 text-gray-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors font-medium">
                  Mis Reservas
                </Link>
                <Link to="/trips/create" onClick={closeMobileMenu} className="block px-4 py-3 mt-2">
                  <Button variant="primary" size="sm" fullWidth>
                    Ofrecer Viaje
                  </Button>
                </Link>
                <Link to="/profile" onClick={closeMobileMenu} className="block px-4 py-3 text-gray-700 hover:bg-primary-50 hover:text-primary-600 rounded-lg transition-colors font-medium">
                  Mi Perfil
                </Link>
                <button
                  onClick={handleLogout}
                  className="block w-full text-left px-4 py-3 text-red-600 hover:bg-red-50 rounded-lg transition-colors font-medium"
                >
                  Cerrar Sesión
                </button>
              </>
            ) : (
              <div className="space-y-3 px-4">
                <Link to="/login" onClick={closeMobileMenu}>
                  <Button variant="outline" size="md" fullWidth>
                    Iniciar Sesión
                  </Button>
                </Link>
                <Link to="/register" onClick={closeMobileMenu}>
                  <Button variant="primary" size="md" fullWidth>
                    Registrarse
                  </Button>
                </Link>
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
  );
};
