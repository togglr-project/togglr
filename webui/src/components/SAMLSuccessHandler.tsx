import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

const SAMLSuccessHandler: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { loginWithTokens } = useAuth();
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const handleSuccess = async () => {
      try {
        // Получаем токены из URL параметров
        const accessToken = searchParams.get('access_token');
        const refreshToken = searchParams.get('refresh_token');

        if (!accessToken || !refreshToken) {
          setError('Missing authentication tokens');
          return;
        }

        // Вызываем функцию логина с токенами из AuthContext
        await loginWithTokens(accessToken, refreshToken);

        // Перенаправляем на dashboard
        navigate('/dashboard', { replace: true });
      } catch (err) {
        console.error('Error during SAML success handling:', err);
        setError('Failed to complete authentication');
      }
    };

    handleSuccess();
  }, [searchParams, navigate, loginWithTokens]);

  if (error) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        flexDirection: 'column',
        gap: '20px'
      }}>
        <div>Authentication Error: {error}</div>
        <button onClick={() => navigate('/login')}>
          Return to Login
        </button>
      </div>
    );
  }

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      height: '100vh' 
    }}>
      <div>Completing authentication...</div>
    </div>
  );
};

export default SAMLSuccessHandler; 