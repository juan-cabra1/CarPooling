import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Search, Plus, TrendingUp, Zap, Clock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import TripCard from '@/components/trips/TripCard'
import { searchService } from '@/services'
import type { SearchTrip } from '@/types'

export default function HomePage() {
  const [featuredTrips, setFeaturedTrips] = useState<SearchTrip[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchFeaturedTrips = async () => {
      try {
        // Obtener viajes destacados (limitamos a 6)
        const response = await searchService.searchTrips({
          sort_by: 'popularity',
          limit: 6,
        })
        setFeaturedTrips(response.trips)
      } catch (error) {
        console.error('Error fetching featured trips:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchFeaturedTrips()
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 via-white to-secondary-50">
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Decorative background elements */}
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-primary-200 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob" />
          <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-secondary-200 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000" />
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-accent-200 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000" />
        </div>

        <div className="relative container mx-auto px-4 py-16 md:py-24">
          <div className="max-w-4xl mx-auto text-center">
            {/* Hero Title */}
            <h1 className="text-5xl sm:text-6xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-primary-600 via-secondary-600 to-accent-600 bg-clip-text text-transparent">
              Viaja Compartiendo
            </h1>

            {/* Subtitle */}
            <p className="text-xl sm:text-2xl text-gray-700 mb-12 max-w-2xl mx-auto">
              Conecta con conductores y pasajeros. Ahorra dinero, conoce personas y cuida el medio ambiente.
            </p>

            {/* CTA Buttons */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <Link to="/search" className="w-full sm:w-auto">
                <Button
                  size="lg"
                  variant="outline"
                  className="w-full sm:w-auto text-lg px-8 py-6 border-2 border-primary bg-white hover:bg-primary hover:text-white transition-all duration-300"
                >
                  <Search className="w-5 h-5 mr-2" />
                  Buscar Viajes
                </Button>
              </Link>

              <Link to="/create-trip" className="w-full sm:w-auto">
                <Button
                  size="lg"
                  className="w-full sm:w-auto text-lg px-8 py-6 bg-gradient-to-r from-primary to-secondary hover:from-primary-600 hover:to-secondary-600 text-white shadow-lg hover:shadow-xl transition-all duration-300"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Publicar Viaje
                </Button>
              </Link>
            </div>

            {/* Quick Stats */}
            <div className="mt-16 grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-3xl mx-auto">
              <div className="bg-white/80 backdrop-blur-sm rounded-lg p-4 shadow-md">
                <div className="text-3xl font-bold text-primary mb-1">1000+</div>
                <div className="text-sm text-gray-600">Viajes realizados</div>
              </div>
              <div className="bg-white/80 backdrop-blur-sm rounded-lg p-4 shadow-md">
                <div className="text-3xl font-bold text-secondary mb-1">500+</div>
                <div className="text-sm text-gray-600">Usuarios activos</div>
              </div>
              <div className="bg-white/80 backdrop-blur-sm rounded-lg p-4 shadow-md">
                <div className="text-3xl font-bold text-accent mb-1">4.8★</div>
                <div className="text-sm text-gray-600">Calificación promedio</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Trips Section */}
      <section className="container mx-auto px-4 py-16">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2 flex items-center gap-2">
              <TrendingUp className="w-8 h-8 text-primary" />
              Viajes Destacados
            </h2>
            <p className="text-gray-600">Los viajes más populares del momento</p>
          </div>
          <Link to="/search">
            <Button variant="outline" className="hidden md:flex">
              Ver Todos
              <Search className="w-4 h-4 ml-2" />
            </Button>
          </Link>
        </div>

        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div
                key={i}
                className="h-96 bg-white/50 rounded-lg animate-pulse"
              />
            ))}
          </div>
        ) : featuredTrips.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {featuredTrips.map((trip) => (
              <TripCard key={trip.id} trip={trip} />
            ))}
          </div>
        ) : (
          <div className="text-center py-16 bg-white/50 rounded-lg">
            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Search className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              No hay viajes disponibles
            </h3>
            <p className="text-gray-600 mb-6">
              Sé el primero en publicar un viaje
            </p>
            <Link to="/create-trip">
              <Button>
                <Plus className="w-4 h-4 mr-2" />
                Publicar Viaje
              </Button>
            </Link>
          </div>
        )}

        {/* Link móvil para ver todos */}
        <div className="mt-8 text-center md:hidden">
          <Link to="/search">
            <Button variant="outline" className="w-full">
              Ver Todos los Viajes
              <Search className="w-4 h-4 ml-2" />
            </Button>
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section className="bg-white/30 backdrop-blur-sm py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-5xl mx-auto">
            <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">
              ¿Por qué elegir CarPooling?
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="text-center p-6 rounded-lg bg-white shadow-md hover:shadow-lg transition-shadow">
                <div className="w-16 h-16 bg-gradient-to-br from-primary-100 to-primary-200 rounded-full flex items-center justify-center mx-auto mb-4">
                  <svg className="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <h3 className="text-xl font-semibold mb-3 text-gray-900">Ahorra Dinero</h3>
                <p className="text-gray-600">Comparte gastos de gasolina y peajes con otros pasajeros. Ahorra hasta un 70% en tus viajes.</p>
              </div>

              <div className="text-center p-6 rounded-lg bg-white shadow-md hover:shadow-lg transition-shadow">
                <div className="w-16 h-16 bg-gradient-to-br from-secondary-100 to-secondary-200 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Zap className="w-8 h-8 text-secondary" />
                </div>
                <h3 className="text-xl font-semibold mb-3 text-gray-900">Rápido y Fácil</h3>
                <p className="text-gray-600">Encuentra o publica un viaje en minutos. Nuestra plataforma es simple e intuitiva.</p>
              </div>

              <div className="text-center p-6 rounded-lg bg-white shadow-md hover:shadow-lg transition-shadow">
                <div className="w-16 h-16 bg-gradient-to-br from-accent-100 to-accent-200 rounded-full flex items-center justify-center mx-auto mb-4">
                  <svg className="w-8 h-8 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <h3 className="text-xl font-semibold mb-3 text-gray-900">Cuida el Planeta</h3>
                <p className="text-gray-600">Reduce tu huella de carbono y contribuye a un futuro más sostenible compartiendo el viaje.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto bg-gradient-to-r from-primary to-secondary rounded-2xl p-8 md:p-12 text-white shadow-2xl">
          <div className="text-center">
            <Clock className="w-16 h-16 mx-auto mb-4 opacity-90" />
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              ¿Listo para tu próximo viaje?
            </h2>
            <p className="text-xl mb-8 opacity-90">
              Únete a nuestra comunidad y empieza a viajar de manera inteligente
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link to="/register">
                <Button
                  size="lg"
                  variant="outline"
                  className="w-full sm:w-auto bg-white text-primary hover:bg-gray-100 border-0"
                >
                  Crear Cuenta
                </Button>
              </Link>
              <Link to="/search">
                <Button
                  size="lg"
                  className="w-full sm:w-auto bg-white/20 backdrop-blur-sm hover:bg-white/30 text-white border-2 border-white"
                >
                  Explorar Viajes
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
