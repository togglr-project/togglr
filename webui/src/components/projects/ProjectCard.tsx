import React from 'react';
import { Card, CardActionArea, CardContent, Chip, Typography, Box } from '@mui/material';

export interface ProjectCardProps {
  id: string;
  name: string;
  description?: string;
  onClick?: () => void;
}

const ProjectCard: React.FC<ProjectCardProps> = ({ id, name, description, onClick }) => {
  return (
    <Card
      sx={{
        background: (theme) =>
          theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
            : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
        backdropFilter: 'blur(8px)',
        boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)',
        transition: 'all 0.2s ease-in-out',
        '&:hover': {
          background: (theme) =>
            theme.palette.mode === 'dark'
              ? 'linear-gradient(135deg, rgba(65, 68, 75, 0.7) 0%, rgba(60, 63, 70, 0.7) 100%)'
              : 'linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(250, 250, 250, 1) 100%)',
          boxShadow: '0 5px 15px 0 rgba(0, 0, 0, 0.1)',
          transform: 'translateY(-3px)'
        }
      }}
    >
      <CardActionArea onClick={onClick}>
        <CardContent>
          <Typography variant="h6" component="div">
            {name}
          </Typography>
          {description ? (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, whiteSpace: 'pre-line' }}>
              {description}
            </Typography>
          ) : null}
          <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
            <Chip label={`ID: ${id}`} size="small" variant="outlined" />
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

export default ProjectCard;
