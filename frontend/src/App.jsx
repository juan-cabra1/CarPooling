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

          {/* Protected routes with MainLayout */}
          <Route
            path="/trips"
            element={
              <MainLayout>
                <ProtectedRoute>
                  <TripsPage />
                </ProtectedRoute>
              </MainLayout>
            }
          />

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
