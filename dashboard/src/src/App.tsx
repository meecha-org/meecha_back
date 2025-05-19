import { Routes, Route, Navigate } from "react-router-dom"
import { Toaster } from "./components/ui/toaster"
import { AuthProvider } from "./contexts/AuthContext"
import ProtectedRoute from "./components/auth/ProtectedRoute"
import DashboardLayout from "./layouts/DashboardLayout"
import LoginPage from "./pages/LoginPage"
import SignupPage from "./pages/SignupPage"
import UsersPage from "./pages/dashboard/UsersPage"
import LabelsPage from "./pages/dashboard/LabelsPage"
import ProvidersPage from "./pages/dashboard/ProvidersPage"
import SessionsPage from "./pages/dashboard/SessionsPage"

function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/users" replace />} />
          <Route path="users" element={<UsersPage />} />
          <Route path="labels" element={<LabelsPage />} />
          <Route path="providers" element={<ProvidersPage />} />
          <Route path="sessions" element={<SessionsPage />} />
        </Route>
        <Route path="*" element={<Navigate to="/signup" replace />} />
      </Routes>
      <Toaster />
    </AuthProvider>
  )
}

export default App
