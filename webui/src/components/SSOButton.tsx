import React, { useState, useEffect } from 'react';
import { 
  Button, 
  Box, 
  Menu, 
  MenuItem, 
  ListItemIcon, 
  ListItemText,
  CircularProgress,
  Typography
} from '@mui/material';
import { useAuth } from '../auth/AuthContext';
import apiClient from '../api/apiClient';
import type { SSOProvider } from '../generated/api/client';

interface SSOButtonProps {
  fullWidth?: boolean;
  variant?: 'contained' | 'outlined' | 'text';
  size?: 'small' | 'medium' | 'large';
  disabled?: boolean;
}

const SSOButton: React.FC<SSOButtonProps> = ({
  fullWidth = false,
  variant = 'contained',
  size = 'medium',
  disabled = false
}) => {
  const { isLoading } = useAuth();
  const [providers, setProviders] = useState<SSOProvider[]>([]);
  const [loadingProviders, setLoadingProviders] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [initiatingSSO, setInitiatingSSO] = useState(false);

  const open = Boolean(anchorEl);

  useEffect(() => {
    const fetchProviders = async () => {
      setLoadingProviders(true);
      try {
        const response = await apiClient.getSSOProviders();
        setProviders(response.data.providers.filter(p => p.enabled));
      } catch (error) {
        console.error('Failed to fetch SSO providers:', error);
        setProviders([]);
      } finally {
        setLoadingProviders(false);
      }
    };

    fetchProviders();
  }, []);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    if (providers.length === 1) {
      // Если только один провайдер, сразу инициируем SSO
      handleSSOInitiate(providers[0].name);
    } else {
      // Если несколько провайдеров, показываем меню
      setAnchorEl(event.currentTarget);
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleSSOInitiate = async (providerName: string) => {
    setInitiatingSSO(true);
    setAnchorEl(null);
    
    try {
      const response = await apiClient.sSOInitiate(providerName);
      if (response.data && response.data.redirect_url) {
        window.location.href = response.data.redirect_url;
      }
    } catch (error) {
      console.error('SSO initiation failed:', error);
    } finally {
      setInitiatingSSO(false);
    }
  };

  // Если нет доступных провайдеров, не показываем кнопку
  if (providers.length === 0 && !loadingProviders) {
    return null;
  }

  return (
    <>
      <Button
        fullWidth={fullWidth}
        variant={variant}
        size={size}
        onClick={handleClick}
        disabled={disabled || isLoading || loadingProviders || initiatingSSO}
        sx={{
          backgroundColor: '#6366f1',
          color: '#ffffff',
          borderColor: '#6366f1',
          minWidth: '80px',
          '&:hover': {
            backgroundColor: '#4f46e5',
            borderColor: '#4f46e5',
          },
          '&:active': {
            backgroundColor: '#4338ca',
            borderColor: '#4338ca',
          },
          '&.Mui-disabled': {
            backgroundColor: 'rgba(99, 102, 241, 0.3)',
            color: 'rgba(255, 255, 255, 0.5)',
            borderColor: 'rgba(99, 102, 241, 0.3)',
          }
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {loadingProviders || initiatingSSO ? (
            <CircularProgress size={16} color="inherit" />
          ) : (
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="currentColor"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
          )}
          <Typography variant="button" sx={{ color: 'inherit', fontWeight: 500 }}>
            SSO
          </Typography>
        </Box>
      </Button>

      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
      >
        {providers.map((provider) => (
          <MenuItem
            key={provider.name}
            onClick={() => handleSSOInitiate(provider.name)}
            disabled={initiatingSSO}
          >
            {provider.icon_url && (
              <ListItemIcon>
                <img 
                  src={provider.icon_url} 
                  alt={provider.display_name}
                  style={{ width: 20, height: 20 }}
                />
              </ListItemIcon>
            )}
            <ListItemText primary={provider.display_name} />
          </MenuItem>
        ))}
      </Menu>
    </>
  );
};

export default SSOButton; 