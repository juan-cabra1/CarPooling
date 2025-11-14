import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import { Search, Plus, Car, Calendar, User, LogOut, ChevronDown } from 'lucide-react'
import { IconMenu2, IconX } from '@tabler/icons-react'
import { motion, useScroll, useMotionValueEvent, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuth()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false)
  const location = useLocation()
  const { scrollY } = useScroll()
  const [visible, setVisible] = useState<boolean>(false)

  useMotionValueEvent(scrollY, 'change', (latest) => {
    if (latest > 100) {
      setVisible(true)
    } else {
      setVisible(false)
    }
  })

  const isActive = (path: string) => location.pathname === path

  const handleMobileMenuClose = () => {
    setIsMobileMenuOpen(false)
  }

  const handleLogout = () => {
    logout()
    setIsProfileMenuOpen(false)
    setIsMobileMenuOpen(false)
  }

  return (
    <motion.div
      className={cn('sticky inset-x-0 top-0 z-40 w-full')}
    >
      {/* Desktop Navbar */}
      <motion.div
        animate={{
          backdropFilter: visible ? 'blur(10px)' : 'none',
          boxShadow: visible
            ? '0 0 24px rgba(34, 42, 53, 0.06), 0 1px 1px rgba(0, 0, 0, 0.05), 0 0 0 1px rgba(34, 42, 53, 0.04), 0 0 4px rgba(34, 42, 53, 0.08), 0 16px 68px rgba(47, 48, 55, 0.05), 0 1px 0 rgba(255, 255, 255, 0.1) inset'
            : 'none',
          width: visible ? '40%' : '100%',
          y: visible ? 20 : 0,
        }}
        transition={{
          type: 'spring',
          stiffness: 200,
          damping: 50,
        }}
        style={{
          minWidth: '800px',
        }}
        className={cn(
          'relative z-[60] mx-auto hidden w-full max-w-7xl flex-row items-center justify-between self-start rounded-full bg-transparent px-4 py-2 lg:flex dark:bg-transparent',
          visible && 'bg-white/80 dark:bg-neutral-950/80'
        )}
      >
        {/* Logo */}
        <Link to="/" className="font-bold text-xl text-primary">
          CarPooling
        </Link>

        {/* Center Navigation */}
        <div className="flex items-center gap-2">
          <Link
            to="/search"
            className={cn(
              'relative px-4 py-2 text-sm font-medium rounded-full transition-colors flex items-center gap-2',
              isActive('/search')
                ? 'bg-primary/10 text-primary'
                : 'text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800'
            )}
          >
            <Search className="w-4 h-4" />
            Buscar Viajes
          </Link>

          {isAuthenticated && (
            <>
              <Link
                to="/create-trip"
                className={cn(
                  'relative px-4 py-2 text-sm font-medium rounded-full transition-colors flex items-center gap-2',
                  isActive('/create-trip')
                    ? 'bg-primary text-white'
                    : 'bg-primary/90 text-white hover:bg-primary'
                )}
              >
                <Plus className="w-4 h-4" />
                Publicar Viaje
              </Link>

              <Link
                to="/my-trips"
                className={cn(
                  'relative px-4 py-2 text-sm font-medium rounded-full transition-colors flex items-center gap-2',
                  isActive('/my-trips')
                    ? 'bg-primary/10 text-primary'
                    : 'text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800'
                )}
              >
                <Car className="w-4 h-4" />
                Mis Viajes
              </Link>

              <Link
                to="/my-bookings"
                className={cn(
                  'relative px-4 py-2 text-sm font-medium rounded-full transition-colors flex items-center gap-2',
                  isActive('/my-bookings')
                    ? 'bg-primary/10 text-primary'
                    : 'text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800'
                )}
              >
                <Calendar className="w-4 h-4" />
                Mis Reservas
              </Link>
            </>
          )}
        </div>

        {/* Right Side Actions */}
        <div className="flex items-center gap-2">
          {!isAuthenticated ? (
            <>
              <Link to="/login">
                <Button variant="outline" size="sm">
                  Login
                </Button>
              </Link>
              <Link to="/register">
                <Button size="sm">Registrarse</Button>
              </Link>
            </>
          ) : (
            <div className="relative">
              <button
                onClick={() => setIsProfileMenuOpen(!isProfileMenuOpen)}
                className="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-full text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800 transition-colors"
              >
                <User className="w-4 h-4" />
                {user?.name || user?.email}
                <ChevronDown className="w-4 h-4" />
              </button>

              <AnimatePresence>
                {isProfileMenuOpen && (
                  <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    className="absolute right-0 mt-2 w-48 bg-white dark:bg-neutral-950 rounded-lg shadow-lg border border-gray-200 dark:border-neutral-800 overflow-hidden"
                  >
                    <Link
                      to="/profile"
                      onClick={() => setIsProfileMenuOpen(false)}
                      className="flex items-center gap-2 px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-neutral-800"
                    >
                      <User className="w-4 h-4" />
                      Perfil
                    </Link>
                    <button
                      onClick={handleLogout}
                      className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20"
                    >
                      <LogOut className="w-4 h-4" />
                      Cerrar Sesión
                    </button>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          )}
        </div>
      </motion.div>

      {/* Mobile Navbar */}
      <motion.div
        animate={{
          backdropFilter: visible ? 'blur(10px)' : 'none',
          boxShadow: visible
            ? '0 0 24px rgba(34, 42, 53, 0.06), 0 1px 1px rgba(0, 0, 0, 0.05), 0 0 0 1px rgba(34, 42, 53, 0.04), 0 0 4px rgba(34, 42, 53, 0.08), 0 16px 68px rgba(47, 48, 55, 0.05), 0 1px 0 rgba(255, 255, 255, 0.1) inset'
            : 'none',
          width: visible ? '90%' : '100%',
          paddingRight: visible ? '12px' : '0px',
          paddingLeft: visible ? '12px' : '0px',
          borderRadius: visible ? '4px' : '2rem',
          y: visible ? 20 : 0,
        }}
        transition={{
          type: 'spring',
          stiffness: 200,
          damping: 50,
        }}
        className={cn(
          'relative z-50 mx-auto flex w-full max-w-[calc(100vw-2rem)] flex-col items-center justify-between bg-transparent px-0 py-2 lg:hidden',
          visible && 'bg-white/80 dark:bg-neutral-950/80'
        )}
      >
        {/* Mobile Header */}
        <div className="flex w-full flex-row items-center justify-between px-4">
          <Link to="/" className="font-bold text-xl text-primary">
            CarPooling
          </Link>

          <button onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
            {isMobileMenuOpen ? (
              <IconX className="text-black dark:text-white" />
            ) : (
              <IconMenu2 className="text-black dark:text-white" />
            )}
          </button>
        </div>

        {/* Mobile Menu */}
        <AnimatePresence>
          {isMobileMenuOpen && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="absolute inset-x-0 top-16 z-50 flex w-full flex-col items-start justify-start gap-2 rounded-lg bg-white px-4 py-6 shadow-[0_0_24px_rgba(34,_42,_53,_0.06),_0_1px_1px_rgba(0,_0,_0,_0.05),_0_0_0_1px_rgba(34,_42,_53,_0.04),_0_0_4px_rgba(34,_42,_53,_0.08),_0_16px_68px_rgba(47,_48,_55,_0.05),_0_1px_0_rgba(255,_255,_255,_0.1)_inset] dark:bg-neutral-950"
            >
              <Link
                to="/search"
                onClick={handleMobileMenuClose}
                className={cn(
                  'w-full flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
                  isActive('/search')
                    ? 'bg-primary/10 text-primary'
                    : 'text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800'
                )}
              >
                <Search className="w-4 h-4" />
                Buscar Viajes
              </Link>

              {isAuthenticated ? (
                <>
                  <Link
                    to="/create-trip"
                    onClick={handleMobileMenuClose}
                    className="w-full"
                  >
                    <Button className="w-full justify-start gap-2">
                      <Plus className="w-4 h-4" />
                      Publicar Viaje
                    </Button>
                  </Link>

                  <Link
                    to="/my-trips"
                    onClick={handleMobileMenuClose}
                    className={cn(
                      'w-full flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
                      isActive('/my-trips')
                        ? 'bg-primary/10 text-primary'
                        : 'text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800'
                    )}
                  >
                    <Car className="w-4 h-4" />
                    Mis Viajes
                  </Link>

                  <Link
                    to="/my-bookings"
                    onClick={handleMobileMenuClose}
                    className={cn(
                      'w-full flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-colors',
                      isActive('/my-bookings')
                        ? 'bg-primary/10 text-primary'
                        : 'text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800'
                    )}
                  >
                    <Calendar className="w-4 h-4" />
                    Mis Reservas
                  </Link>

                  <div className="w-full border-t border-gray-200 dark:border-neutral-800 my-2" />

                  <Link
                    to="/profile"
                    onClick={handleMobileMenuClose}
                    className="w-full flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg text-neutral-600 dark:text-neutral-300 hover:bg-gray-100 dark:hover:bg-neutral-800 transition-colors"
                  >
                    <User className="w-4 h-4" />
                    Perfil
                  </Link>

                  <button
                    onClick={handleLogout}
                    className="w-full flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                  >
                    <LogOut className="w-4 h-4" />
                    Cerrar Sesión
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/login"
                    onClick={handleMobileMenuClose}
                    className="w-full"
                  >
                    <Button variant="outline" className="w-full">
                      Login
                    </Button>
                  </Link>
                  <Link
                    to="/register"
                    onClick={handleMobileMenuClose}
                    className="w-full"
                  >
                    <Button className="w-full">Registrarse</Button>
                  </Link>
                </>
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </motion.div>
  )
}
