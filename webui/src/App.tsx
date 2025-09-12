import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from './theme/ThemeContext';
import { AuthProvider } from './auth/AuthContext';
import { LicenseProvider } from './auth/LicenseContext';
import { AccessibilityProvider } from './components/AccessibilityProvider';
import { ConfigProvider } from './config/ConfigContext';
import { createContext, useContext, useState, type ReactNode } from 'react';
import { type AlertColor } from '@mui/material';
import NotificationSnackbar from './components/NotificationSnackbar';
import LicenseGuard from './components/License/LicenseGuard';
import LoginPage from './pages/LoginPage';
import ForgotPasswordPage from './pages/ForgotPasswordPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import SSOCallbackPage from './pages/SSOCallbackPage';
import SAMLSuccessHandler from './components/SAMLSuccessHandler';
import ErrorBoundary from './components/ErrorBoundary';
import './App.css';

// Create a client for React Query
const queryClient = new QueryClient();

// Global notification context
interface Notification {
  id: string;
  message: string;
  type: AlertColor;
  duration?: number;
}

interface NotificationContextType {
  showNotification: (message: string, type: AlertColor, duration?: number) => void;
}

const NotificationContext = createContext<NotificationContextType>({
  showNotification: () => {},
});

export const useNotification = () => useContext(NotificationContext);

interface NotificationProviderProps {
  children: ReactNode;
}

const NotificationProvider: React.FC<NotificationProviderProps> = ({ children }) => {
  const [notification, setNotification] = useState<Notification | null>(null);

  const showNotification = (message: string, type: AlertColor, duration = 6000) => {
    const id = Date.now().toString();
    setNotification({ id, message, type, duration });
  };

  const handleClose = (_id: string) => {
    setNotification(null);
  };

  return (
    <NotificationContext.Provider value={{ showNotification }}>
      {children}
      <NotificationSnackbar
        notification={notification}
        onClose={handleClose}
        position="bottom"
      />
    </NotificationContext.Provider>
  );
};

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider>
          <AccessibilityProvider>
            <NotificationProvider>
              <ConfigProvider>
                <AuthProvider>
                  <LicenseProvider>
                    <Router>
                      <Routes>
                        <Route path="/login" element={<LoginPage />} />
                        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
                        <Route path="/reset-password" element={<ResetPasswordPage />} />
                        <Route path="/api/v1/auth/sso/callback" element={<SSOCallbackPage />} />
                        <Route path="/auth/saml/callback" element={<SSOCallbackPage />} />
                        <Route path="/auth/saml/logout" element={<SSOCallbackPage />} />
                        <Route path="/auth/saml/success" element={<SAMLSuccessHandler />} />
                      </Routes>
                    </Router>
                  </LicenseProvider>
                </AuthProvider>
              </ConfigProvider>
            </NotificationProvider>
          </AccessibilityProvider>
        </ThemeProvider>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}

export default App;
