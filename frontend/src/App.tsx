import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from '@/context/AuthContext'
import Layout from '@/components/layout/Layout'
import AdminLayout from '@/components/admin/AdminLayout'
import AdminRoute from '@/components/routes/AdminRoute'

// Pages
import HomePage from '@/pages/HomePage'
import LoginPage from '@/pages/LoginPage'
import RegisterPage from '@/pages/RegisterPage'
import VerifyEmailPage from '@/pages/VerifyEmailPage'
import ResendVerificationPage from '@/pages/ResendVerificationPage'
import ForgotPasswordPage from '@/pages/ForgotPasswordPage'
import ResetPasswordPage from '@/pages/ResetPasswordPage'
import SearchPage from '@/pages/SearchPage'
import CreateTripPage from '@/pages/CreateTripPage'
import MyTripsPage from '@/pages/MyTripsPage'
import MyBookingsPage from '@/pages/MyBookingsPage'
import ProfilePage from '@/pages/ProfilePage'
import TripDetailPage from '@/pages/TripDetailPage'
import EditTripPage from '@/pages/EditTripPage'
import UnauthorizedPage from '@/pages/UnauthorizedPage'

// Admin Pages
import AdminDashboardPage from '@/pages/admin/AdminDashboardPage'
import AdminUsersPage from '@/pages/admin/AdminUsersPage'
import AdminTripsPage from '@/pages/admin/AdminTripsPage'
import AdminBookingsPage from '@/pages/admin/AdminBookingsPage'

// Protected Route wrapper
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, loading, user } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Cargando...</div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  // Si el usuario no tiene el email verificado, mostrar mensaje
  if (user && !user.email_verified) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="max-w-md w-full bg-yellow-50 border border-yellow-200 rounded-lg p-6 text-center">
          <div className="mb-4">
            <svg className="w-16 h-16 text-yellow-500 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h2 className="text-xl font-bold text-gray-900 mb-2">
            Email no verificado
          </h2>
          <p className="text-gray-700 mb-4">
            Debes verificar tu correo electrónico antes de acceder a esta funcionalidad.
          </p>
          <div className="space-y-2">
            <a
              href="/resend-verification"
              className="block w-full bg-primary text-white py-2 px-4 rounded hover:bg-primary-600 transition"
            >
              Reenviar email de verificación
            </a>
            <a
              href="/"
              className="block w-full bg-gray-200 text-gray-800 py-2 px-4 rounded hover:bg-gray-300 transition"
            >
              Volver al inicio
            </a>
          </div>
        </div>
      </div>
    )
  }

  return <>{children}</>
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="login" element={<LoginPage />} />
            <Route path="register" element={<RegisterPage />} />
            <Route path="verify-email" element={<VerifyEmailPage />} />
            <Route path="resend-verification" element={<ResendVerificationPage />} />
            <Route path="forgot-password" element={<ForgotPasswordPage />} />
            <Route path="reset-password" element={<ResetPasswordPage />} />
            <Route path="search" element={<SearchPage />} />
            <Route path="unauthorized" element={<UnauthorizedPage />} />

            {/* Trip routes - More specific routes first */}
            <Route path="trips/:id/edit" element={
              <ProtectedRoute>
                <EditTripPage />
              </ProtectedRoute>
            } />
            <Route path="trips/:id" element={<TripDetailPage />} />

            {/* Protected routes */}
            <Route path="create-trip" element={
              <ProtectedRoute>
                <CreateTripPage />
              </ProtectedRoute>
            } />
            <Route path="my-trips" element={
              <ProtectedRoute>
                <MyTripsPage />
              </ProtectedRoute>
            } />
            <Route path="my-bookings" element={
              <ProtectedRoute>
                <MyBookingsPage />
              </ProtectedRoute>
            } />
            <Route path="profile" element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            } />
          </Route>

          {/* Admin routes */}
          <Route path="/admin" element={
            <AdminRoute>
              <AdminLayout />
            </AdminRoute>
          }>
            <Route index element={<AdminDashboardPage />} />
            <Route path="users" element={<AdminUsersPage />} />
            <Route path="trips" element={<AdminTripsPage />} />
            <Route path="bookings" element={<AdminBookingsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App
