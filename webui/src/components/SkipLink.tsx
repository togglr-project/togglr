import React from 'react';
import { Box, Link } from '@mui/material';

interface SkipLinkProps {
  href: string;
  children: React.ReactNode;
}

const SkipLink: React.FC<SkipLinkProps> = ({ href, children }) => {
  return (
    <Box
      component={Link}
      href={href}
      className="skip-link"
      sx={{
        position: 'absolute',
        top: '6px',
        left: '-9999px',
        background: 'primary.main',
        color: 'primary.contrastText',
        padding: 1,
        textDecoration: 'none',
        borderRadius: 0.5,
        zIndex: 10000,
        transition: 'left 0.2s ease',
        '&:focus': {
          left: '6px',
        },
        '&:hover': {
          textDecoration: 'none',
        },
      }}
    >
      {children}
    </Box>
  );
};

export default SkipLink; 