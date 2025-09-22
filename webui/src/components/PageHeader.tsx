import React from 'react';
import { 
  Box, 
  Typography, 
  Chip,
  useTheme,
  Skeleton
} from '@mui/material';

interface PageHeaderProps {
  title: string;
  subtitle?: string;
  icon?: React.ReactNode;
  badge?: {
    label: string;
    color: 'primary' | 'secondary' | 'error' | 'warning' | 'info' | 'success';
    variant?: 'filled' | 'outlined';
  };
  loading?: boolean;
  children?: React.ReactNode;
}

const PageHeader: React.FC<PageHeaderProps> = ({ 
  title, 
  subtitle, 
  icon, 
  badge, 
  loading = false,
  children
}) => {
  const theme = useTheme();

  if (loading) {
    return (
      <Box sx={{ mb: 1.5, width: '100%' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 0.5 }}>
          <Skeleton variant="circular" width={28} height={28} />
          <Skeleton variant="text" width={200} height={20} />
          {badge && <Skeleton variant="rectangular" width={80} height={20} sx={{ borderRadius: 1 }} />}
        </Box>
        {subtitle && <Skeleton variant="text" width={300} height={16} />}
      </Box>
    );
  }

  return (
    <Box sx={{ 
      mb: 1.5, 
      width: '100%',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'flex-start',
      flexWrap: 'wrap',
      gap: 1.5
    }}>
      <Box sx={{ flex: 1, minWidth: 0 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 0.5 }}>
          {icon && (
            <Box sx={{ 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'center',
              width: 28,
              height: 28,
              borderRadius: 1.5,
              backgroundColor: theme.palette.mode === 'dark' 
                ? 'rgba(130, 82, 255, 0.1)' 
                : 'rgba(130, 82, 255, 0.08)',
              color: theme.palette.primary.main,
            }}>
              {icon}
            </Box>
          )}
          <Typography 
            variant="h6" 
            component="h1" 
            sx={{ 
              fontWeight: 700,
              lineHeight: 1.2,
              fontSize: '1.2rem',
              color: 'primary.main'
            }}
          >
            {title}
          </Typography>
          {badge && (
            <Chip 
              label={badge.label} 
              color={badge.color} 
              variant={badge.variant || 'filled'}
              size="small"
              sx={{ 
                height: 20,
                fontSize: '0.7rem',
                fontWeight: 500,
              }}
            />
          )}
        </Box>
        {subtitle && (
          <Typography 
            variant="body2" 
            sx={{ 
              ml: icon ? 5 : 0,
              fontSize: '0.9rem',
              lineHeight: 1.4,
              whiteSpace: 'pre-wrap',
              color: 'primary.light'
            }}
          >
            {subtitle}
          </Typography>
        )}
      </Box>
      {children && (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {children}
        </Box>
      )}
    </Box>
  );
};

export default PageHeader; 