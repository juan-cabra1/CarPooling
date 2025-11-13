import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from '@/contexts/AuthContext';
import { MainLayout } from '@/components/layout/MainLayout';
import { ProtectedRoute } from '@/components/common/ProtectedRoute';

// Pages
import { HomePage } from '@/pages/HomePage';
import { LoginPage } from '@/pages/auth/LoginPage';
import { RegisterPage } from '@/pages/auth/RegisterPage';
import { RegisterDebug } from '@/pages/auth/RegisterDebug';
import { TripsPage } from '@/pages/trips/TripsPage';
import { TripsListPage } from '@/pages/trips/TripsListPage';
import { TripDetailsPage } from '@/pages/trips/TripDetailsPage';
import { TripFormPage } from '@/pages/trips/TripFormPage';
import { ProfilePage } from '@/pages/profile/ProfilePage';
import { NotFoundPage } from '@/pages/NotFoundPage';

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          {/* Public routes without MainLayout */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/register-debug" element={<RegisterDebug />} />

          {/* Public routes with MainLayout */}
          <Route
            path="/"
            element={
              <MainLayout>
                <HomePage />
              </MainLayout>
            }
          />

          {/* Trips routes - public and protected */}
          <Route
            path="/trips"
            element={
              <MainLayout>
                <TripsListPage />
              </MainLayout>
            }
          />

          <Route
            path="/trips/new"
            element={
              <MainLayout>
                <ProtectedRoute>
                  <TripFormPage />
                </ProtectedRoute>
              </MainLayout>
            }
          />

          <Route
            path="/trips/:id"
            element={
              <MainLayout>
                <TripDetailsPage />
              </MainLayout>
            }
          />

          <Route
            path="/trips/:id/edit"
            element={
              <MainLayout>
                <ProtectedRoute>
                  <TripFormPage />
                </ProtectedRoute>
              </MainLayout>
            }
          />

          {/* Protected routes with MainLayout */}
          <Route
            path="/profile"
            element={
              <MainLayout>
                <ProtectedRoute>
                  <ProfilePage />
                </ProtectedRoute>
              </MainLayout>
            }
          />

          {/* 404 Page */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;
