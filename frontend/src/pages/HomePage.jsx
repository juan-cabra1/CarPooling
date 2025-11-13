import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button, Card } from '@/components/common';

export const HomePage = () => {
  const { isAuthenticated } = useAuth();

  return (
    <div className="space-y-16 lg:space-y-20">
      {/* Hero Section */}
      <section className="text-center py-12 lg:py-20 animate-fade-in">
        <div className="inline-block mb-6 lg:mb-8">
          <span className="text-6xl lg:text-8xl">🚗</span>
        </div>
        <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-gray-900 mb-4 lg:mb-6 px-4">
          Compartí tu Viaje
        </h1>
        <p className="text-lg sm:text-xl lg:text-2xl text-gray-600 mb-8 lg:mb-10 max-w-3xl mx-auto px-4 leading-relaxed">
          Conectate con personas que viajan en tu misma dirección. Ahorrá dinero, reducí tu huella de carbono y hacé nuevos amigos en el camino.
        </p>
        <div className="flex flex-col sm:flex-row justify-center gap-4 px-4">
          {isAuthenticated ? (
            <>
              <Link to="/trips">
                <Button variant="primary" size="lg" className="w-full sm:w-auto">
                  Buscar Viaje
                </Button>
              </Link>
              <Link to="/trips/create">
                <Button variant="outline" size="lg" className="w-full sm:w-auto">
                  Ofrecer Viaje
                </Button>
              </Link>
            </>
          ) : (
            <>
              <Link to="/register" className="w-full sm:w-auto">
                <Button variant="primary" size="lg" fullWidth className="sm:w-auto">
                  Comenzar Gratis
                </Button>
              </Link>
              <Link to="/login" className="w-full sm:w-auto">
                <Button variant="outline" size="lg" fullWidth className="sm:w-auto">
                  Iniciar Sesión
                </Button>
              </Link>
            </>
          )}
        </div>
      </section>

      {/* Features Section */}
      <section className="py-12 lg:py-16">
        <h2 className="text-3xl lg:text-4xl font-bold text-center text-gray-900 mb-10 lg:mb-16 px-4">
          ¿Por qué elegir CarPooling?
        </h2>
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 lg:gap-8 px-4">
          <Card>
            <div className="text-center">
              <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold mb-2">Save Money</h3>
              <p className="text-gray-600">
                Share travel costs and reduce your expenses significantly
              </p>
            </div>
          </Card>

          <Card>
            <div className="text-center">
              <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold mb-2">Go Green</h3>
              <p className="text-gray-600">
                Reduce your carbon footprint and help protect the environment
              </p>
            </div>
          </Card>

          <Card>
            <div className="text-center">
              <div className="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold mb-2">Meet People</h3>
              <p className="text-gray-600">
                Connect with fellow travelers and make new friends
              </p>
            </div>
          </Card>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-16 bg-gray-50 -mx-4 px-4">
        <div className="container mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">
            How It Works
          </h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center text-xl font-bold mx-auto mb-4">
                1
              </div>
              <h3 className="text-xl font-semibold mb-2">Search for a Ride</h3>
              <p className="text-gray-600">
                Enter your destination and find available rides
              </p>
            </div>

            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center text-xl font-bold mx-auto mb-4">
                2
              </div>
              <h3 className="text-xl font-semibold mb-2">Book Your Seat</h3>
              <p className="text-gray-600">
                Choose your ride and reserve your spot instantly
              </p>
            </div>

            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center text-xl font-bold mx-auto mb-4">
                3
              </div>
              <h3 className="text-xl font-semibold mb-2">Travel Together</h3>
              <p className="text-gray-600">
                Meet at the pickup point and enjoy your journey
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};
