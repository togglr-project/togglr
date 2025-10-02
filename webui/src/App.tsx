import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from './theme/ThemeContext';
import { AuthProvider } from './auth/AuthContext';
import { LicenseProvider } from './auth/LicenseContext';
import { AccessibilityProvider } from './components/AccessibilityProvider';
import { ConfigProvider } from './config/ConfigContext';
import { createContext, useContext, type ReactNode } from 'react';
import { type AlertColor } from '@mui/material';
import NotificationSnackbar from './components/NotificationSnackbar';
import LicenseGuard from './components/License/LicenseGuard';
import LoginPage from './pages/LoginPage';
import ForgotPasswordPage from './pages/ForgotPasswordPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import SSOCallbackPage from './pages/SSOCallbackPage';
import SAMLSuccessHandler from './components/SAMLSuccessHandler';
import DashboardPage from './pages/DashboardPage';
import ProjectPage from './pages/ProjectPage';
import ProjectSchedulingPage from './pages/ProjectSchedulingPage';
import ProjectSegmentsPage from './pages/ProjectSegmentsPage';
import ProjectSettingsPage from './pages/ProjectSettingsPage';
import ProjectsPage from './pages/ProjectsPage';
import CategoriesPage from './pages/CategoriesPage';
import ProjectTagsPage from './pages/ProjectTagsPage';
import ProjectPermissionsPage from './pages/ProjectPermissionsPage';
import PendingChangesPage from './pages/PendingChangesPage';
import AuditLogPage from './pages/AuditLogPage';
import AdminPage from './pages/AdminPage';
import AccountPage from './pages/AccountPage';
import ErrorBoundary from './components/ErrorBoundary';
import './App.css';

// Create a client for React Query
const queryClient = new QueryClient();
// Expose globally for realtime handlers (minimal integration without refactor)
if (typeof window !== 'undefined') {
  window.__RQ = queryClient;
}

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

  const handleClose = () => {
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

import { useEffect, useState } from 'react';
import { initRealtime, stopRealtime } from './realtime';

function App() {
  const [currentPath, setCurrentPath] = useState(window.location.pathname);

  useEffect(() => {
    // Listen for route changes
    const handleRouteChange = () => {
      setCurrentPath(window.location.pathname);
    };

    // Listen for popstate events (back/forward navigation)
    window.addEventListener('popstate', handleRouteChange);
    
    // Also listen for pushstate/replacestate (programmatic navigation)
    const originalPushState = window.history.pushState;
    const originalReplaceState = window.history.replaceState;
    
    window.history.pushState = function(...args) {
      originalPushState.apply(this, args);
      handleRouteChange();
    };
    
    window.history.replaceState = function(...args) {
      originalReplaceState.apply(this, args);
      handleRouteChange();
    };

    return () => {
      window.removeEventListener('popstate', handleRouteChange);
      window.history.pushState = originalPushState;
      window.history.replaceState = originalReplaceState;
    };
  }, []);

  useEffect(() => {
    // Try to detect project/env from storage and URL
    const fromStorage = (keys: string[]): string | '' => {
      for (const k of keys) {
        const v = localStorage.getItem(k);
        if (v) return v;
      }
      return '' as const;
    };

    const path = currentPath;
    let projectId = fromStorage(['currentProjectId', 'selectedProjectId', 'projectId']);

    // Fallback: parse from URL like /projects/:projectId/...
    if (!projectId) {
      const m = path.match(/\/projects\/([^\/]+)/);
      if (m && m[1]) projectId = m[1];
    }

    let envId = fromStorage([
      'currentEnvId',
      'selectedEnvironmentId',
      'environmentId',
      'env_id',
      'activeEnvironmentId',
    ]);

    // Fallback: parse envId from URL like /projects/:projectId/... (if not found in storage)
    if (!envId) {
      const envMatch = path.match(/\/projects\/[^\/]+\/([^\/]+)/);
      if (envMatch && envMatch[1]) {
        envId = envMatch[1];
        console.log('[Realtime] Found envId in URL:', envId);
      }
    }
    
    // Additional fallback: try to get envId from query parameters
    if (!envId) {
      const urlParams = new URLSearchParams(window.location.search);
      const queryEnvId = urlParams.get('env_id') || urlParams.get('environment_id');
      if (queryEnvId) {
        envId = queryEnvId;
        console.log('[Realtime] Found envId in query params:', envId);
      }
    }
    
    // Additional fallback: try to get envId from DOM (environment selector)
    if (!envId) {
      // Look for environment selector in DOM
      const envSelect = document.querySelector('select[name*="env"], select[id*="env"], [data-testid*="env"]') as HTMLSelectElement;
      if (envSelect && envSelect.value) {
        // The select value is environment.key, we need to find the corresponding environment.id
        const environmentKey = envSelect.value;
        console.log('[Realtime] Found environment key in DOM selector:', environmentKey);
        
        // Try to find the environment ID by looking for the selected option's data attributes
        const selectedOption = envSelect.querySelector(`option[value="${environmentKey}"]`);
        if (selectedOption) {
          const envIdFromData = selectedOption.getAttribute('data-env-id') || selectedOption.getAttribute('data-environment-id');
          if (envIdFromData) {
            envId = envIdFromData;
            console.log('[Realtime] Found envId from selected option data:', envId);
          }
        }
        
        // If still no envId, try to find it by looking for environment data in the page
        if (!envId) {
          // Look for environment data in script tags or global variables
          const envDataScript = document.querySelector('script[type="application/json"][data-env-data]');
          if (envDataScript) {
            try {
              const envData = JSON.parse(envDataScript.textContent || '{}');
              const currentEnv = envData.find((env: any) => env.key === environmentKey);
              if (currentEnv && currentEnv.id) {
                envId = currentEnv.id.toString();
                console.log('[Realtime] Found envId from environment data:', envId);
              }
            } catch (e) {
              console.log('[Realtime] Failed to parse environment data:', e);
            }
          }
        }
      }
      
      // Also try to find any element with environment ID
      const envElements = document.querySelectorAll('[data-environment-id], [data-env-id]');
      if (envElements.length > 0) {
        const firstEnvElement = envElements[0];
        envId = firstEnvElement.getAttribute('data-environment-id') || firstEnvElement.getAttribute('data-env-id') || '';
        if (envId) {
          console.log('[Realtime] Found envId in DOM data attributes:', envId);
        }
      }
    }

    const token = localStorage.getItem('accessToken') || undefined;

    let stop: (() => void) | undefined;

    // Check if we're on a project-related page
    const isProjectPage = path.startsWith('/projects/') && path !== '/projects';
    const shouldConnectWS = isProjectPage && projectId && envId;

    console.log('[Realtime] Debug info:', { 
      projectId, 
      envId, 
      token: token ? 'present' : 'missing', 
      path,
      isProjectPage,
      shouldConnectWS,
      localStorage: {
        currentProjectId: localStorage.getItem('currentProjectId'),
        currentEnvId: localStorage.getItem('currentEnvId'),
        selectedEnvironmentId: localStorage.getItem('selectedEnvironmentId'),
        environmentId: localStorage.getItem('environmentId'),
        env_id: localStorage.getItem('env_id'),
        activeEnvironmentId: localStorage.getItem('activeEnvironmentId'),
        accessToken: localStorage.getItem('accessToken') ? 'present' : 'missing'
      },
      urlSearch: window.location.search
    });

    if (shouldConnectWS) {
      console.log('[Realtime] Connecting WS', { projectId, envId });
      stop = initRealtime({ projectId, envId, token });
    } else if (isProjectPage && projectId && !envId) {
      // If we're on a project page but don't have envId, wait a bit and try again
      console.log('[Realtime] Project page but no envId, waiting for DOM to load...');
      setTimeout(() => {
        // Try to find envId again after DOM loads
        const envSelect = document.querySelector('select[name*="env"], select[id*="env"], [data-testid*="env"]') as HTMLSelectElement;
        if (envSelect && envSelect.value) {
          const delayedEnvId = envSelect.value;
          console.log('[Realtime] Found envId after delay:', delayedEnvId);
          stop = initRealtime({ projectId, envId: delayedEnvId, token });
        } else {
          console.log('[Realtime] Still no envId found after delay');
        }
      }, 1000);
    } else {
      console.log('[Realtime] Skipped WS init:', { 
        reason: !isProjectPage ? 'not a project page' : 'no project/env id',
        projectId, 
        envId, 
        path 
      });
    }

    return () => {
      try { stop?.(); } catch {}
      stopRealtime();
    };
  }, [currentPath]);

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

                        <Route path="/dashboard" element={<LicenseGuard><DashboardPage /></LicenseGuard>} />
                        <Route path="/projects" element={<LicenseGuard><ProjectsPage /></LicenseGuard>} />
                        <Route path="/categories" element={<LicenseGuard><CategoriesPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId" element={<LicenseGuard><ProjectPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId/scheduling" element={<LicenseGuard><ProjectSchedulingPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId/segments" element={<LicenseGuard><ProjectSegmentsPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId/tags" element={<LicenseGuard><ProjectTagsPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId/permissions" element={<LicenseGuard><ProjectPermissionsPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId/pending-changes" element={<LicenseGuard><PendingChangesPage /></LicenseGuard>} />
                                                <Route path="/projects/:projectId/audit-log" element={<LicenseGuard><AuditLogPage /></LicenseGuard>} />
                        <Route path="/projects/:projectId/settings" element={<LicenseGuard><ProjectSettingsPage /></LicenseGuard>} />
                        <Route path="/admin" element={<LicenseGuard><AdminPage /></LicenseGuard>} />
                        <Route path="/account" element={<LicenseGuard><AccountPage /></LicenseGuard>} />
                        <Route path="/" element={<Navigate to="/dashboard" replace />} />
                        <Route path="*" element={<Navigate to="/dashboard" replace />} />
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
