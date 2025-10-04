import React from 'react';
import { IconButton, Tooltip } from '@mui/material';
import { Brightness4, Brightness7 } from '@mui/icons-material';
import { useTheme as useAppTheme } from '../theme/ThemeContext';

interface ThemeToggleProps {
  size?: 'small' | 'medium' | 'large';
  showTooltip?: boolean;
}

export const ThemeToggle: React.FC<ThemeToggleProps> = ({
  size = 'medium',
  showTooltip = true
}) => {
  const { mode, toggleTheme } = useAppTheme();

  const iconSize = size === 'small' ? 20 : size === 'large' ? 28 : 24;
  
  // Determine if the current theme is dark (only 'dark' is considered dark)
  const isDarkTheme = mode !== 'light';

  const button = (
    <IconButton
      onClick={toggleTheme}
      color="inherit"
      size={size}
      sx={{
        transition: 'all 0.3s ease',
        '&:hover': {
          transform: 'rotate(180deg)',
          backgroundColor: isDarkTheme 
            ? 'rgba(255, 255, 255, 0.08)' 
            : 'rgba(0, 0, 0, 0.04)'
        },
        '& .MuiSvgIcon-root': {
          transition: 'transform 0.3s ease'
        }
      }}
      aria-label="Switch theme"
    >
      {isDarkTheme ? (
        <Brightness7 sx={{ 
          fontSize: iconSize,
          color: '#FFD700' // Golden color for sun on dark background
        }} />
      ) : (
        <Brightness4 sx={{ 
          fontSize: iconSize,
          color: '#8252FF' // Purple color for moon on light background
        }} />
      )}
    </IconButton>
  );

  if (!showTooltip) {
    return button;
  }

  return (
    <Tooltip
      title={isDarkTheme ? 'Switch to light theme' : 'Switch to dark theme'}
      placement="bottom"
      arrow
    >
      {button}
    </Tooltip>
  );
};

export default ThemeToggle;
