import React from 'react';
import { 
  Breadcrumbs as MuiBreadcrumbs, 
  Link, 
  Typography, 
  Box,
  useTheme
} from '@mui/material';
import { 
  NavigateNext as NavigateNextIcon,
  Home as HomeIcon
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../api/apiClient';

interface BreadcrumbItem {
  label: string;
  path: string;
  icon?: React.ReactNode;
}

const Breadcrumbs: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();

  // Get project data for breadcrumbs
  const { data: projects, isLoading: projectsLoading } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const response = await apiClient.listProjects();
      return response.data;
    }
  });

  // Get current project data
  const projectId = location.pathname.match(/\/projects\/([^\/]+)/)?.[1];
  const currentProject = projects?.find(p => p.id === projectId);

  const generateBreadcrumbs = (): BreadcrumbItem[] => {
    const pathSegments = location.pathname.split('/').filter(Boolean);
    const breadcrumbs: BreadcrumbItem[] = [];
    
    // Debug logging
    console.log('Breadcrumbs debug:', {
      pathname: location.pathname,
      pathSegments,
      currentProject: currentProject?.name,
      projectId,
      projectsLoading,
      projects: projects?.map(p => ({ id: p.id, name: p.name })),
      projectsCount: projects?.length
    });

    // Always add home
    breadcrumbs.push({
      label: 'Dashboard',
      path: '/dashboard',
      icon: <HomeIcon fontSize="small" />
    });

    // Handle different route patterns
    if (pathSegments[0] === 'projects') {
      // Add Projects
      breadcrumbs.push({
        label: 'Projects',
        path: '/projects'
      });

      // If we have a project ID
      if (pathSegments[1]) {
        const projectPath = `/projects/${pathSegments[1]}`;
        
        // Add project name (or fallback if project not found)
        breadcrumbs.push({
          label: currentProject?.name || (projectsLoading ? 'Loading...' : `Project ${pathSegments[1]}`),
          path: projectPath
        });

        // Add subpage if exists
        if (pathSegments[2]) {
          const subpage = pathSegments[2];
          const subpagePath = `${projectPath}/${subpage}`;
          
          switch (subpage) {
            case 'scheduling':
              breadcrumbs.push({
                label: 'Scheduling',
                path: subpagePath
              });
              break;
            case 'segments':
              breadcrumbs.push({
                label: 'Segments',
                path: subpagePath
              });
              break;
            case 'settings':
              breadcrumbs.push({
                label: 'Settings',
                path: subpagePath
              });
              break;
            default:
              // Unknown subpage
              breadcrumbs.push({
                label: subpage.charAt(0).toUpperCase() + subpage.slice(1),
                path: subpagePath
              });
              break;
          }
        } else {
          // Main project page - add Features
          breadcrumbs.push({
            label: 'Features',
            path: projectPath
          });
        }
      }
    } else if (pathSegments[0] === 'issues') {
      breadcrumbs.push({
        label: 'Issues',
        path: '/issues'
      });
    } else if (pathSegments[0] === 'admin') {
      breadcrumbs.push({
        label: 'Admin',
        path: '/admin'
      });
    }

    return breadcrumbs;
  };

  const breadcrumbs = generateBreadcrumbs();

  // Don't show breadcrumbs on login pages or if only one item
  if (breadcrumbs.length <= 1 || location.pathname.startsWith('/login') || location.pathname.startsWith('/forgot-password') || location.pathname.startsWith('/reset-password')) {
    return null;
  }

  return (
    <Box sx={{ mb: 3, width: '100%' }}>
      <MuiBreadcrumbs 
        separator={<NavigateNextIcon fontSize="small" />}
        aria-label="breadcrumb"
        sx={{
          '& .MuiBreadcrumbs-separator': {
            color: theme.palette.text.secondary,
          },
        }}
      >
        {breadcrumbs.map((item, index) => {
          const isLast = index === breadcrumbs.length - 1;
          
          if (isLast) {
            return (
              <Typography 
                key={item.path}
                color="text.primary" 
                sx={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  gap: 0.5,
                  fontWeight: 500,
                  fontSize: '0.875rem'
                }}
              >
                {item.icon}
                {item.label}
              </Typography>
            );
          }

          return (
            <Link
              key={item.path}
              color="inherit"
              href="#"
              onClick={(e) => {
                e.preventDefault();
                navigate(item.path);
              }}
              sx={{ 
                display: 'flex', 
                alignItems: 'center', 
                gap: 0.5,
                textDecoration: 'none',
                color: theme.palette.text.secondary,
                fontSize: '0.875rem',
                '&:hover': {
                  color: theme.palette.primary.main,
                  textDecoration: 'underline',
                },
              }}
            >
              {item.icon}
              {item.label}
            </Link>
          );
        })}
      </MuiBreadcrumbs>
    </Box>
  );
};

export default Breadcrumbs; 