import React, { useEffect, useState } from 'react';
import { useSearchParams, Navigate } from 'react-router-dom';
import { 
  Box, 
  Container, 
  CircularProgress, 
  Alert, 
  Typography,
  Paper 
} from '@mui/material';
import { useAuth } from '../auth/AuthContext';
import LogoImg from '../components/LogoImg';

const SSOCallbackPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const { handleSSOCallback, isLoading, error, isAuthenticated } = useAuth();
  const [isProcessing, setIsProcessing] = useState(true);

  useEffect(() => {
    const processCallback = async () => {
      // For SAML, the response comes as a POST parameter or in the URL
      const samlResponse = searchParams.get('SAMLResponse');
      const state = searchParams.get('RelayState'); // SAML uses RelayState instead of state

      if (!samlResponse || !state) {
        setIsProcessing(false);
        return;
      }

      try {
        await handleSSOCallback(samlResponse, state);
      } catch (error) {
        console.error('SSO callback processing failed:', error);
      } finally {
        setIsProcessing(false);
      }
    };

    processCallback();
  }, [searchParams, handleSSOCallback]);

  // If already authenticated, redirect to dashboard
  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  // If no SAML response or state parameters, show error
  if (!isProcessing && (!searchParams.get('SAMLResponse') || !searchParams.get('RelayState'))) {
    return (
      <Container maxWidth="sm">
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '100vh',
          }}
        >
          <Paper
            elevation={3}
            sx={{
              p: 4,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              width: '100%',
            }}
          >
            <Box sx={{ mb: 4 }}>
              <LogoImg size="large" />
            </Box>
            <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
              Invalid SSO callback. Missing required parameters.
            </Alert>
            <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center' }}>
              Please try signing in again.
            </Typography>
          </Paper>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
        }}
      >
        <Paper
          elevation={3}
          sx={{
            p: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            width: '100%',
          }}
        >
          <Box sx={{ mb: 4 }}>
            <LogoImg size="large" />
          </Box>
          
          {error && (
            <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
              {error}
            </Alert>
          )}
          
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <CircularProgress size={48} sx={{ mb: 2 }} />
            <Typography variant="h6" sx={{ mb: 1 }}>
              Completing SSO Authentication
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center' }}>
              Please wait while we complete your sign-in...
            </Typography>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default SSOCallbackPage; 