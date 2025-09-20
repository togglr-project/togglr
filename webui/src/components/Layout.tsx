import React, { type ReactNode, useState } from 'react';
import { 
  Box, 
  AppBar, 
  Toolbar, 
  IconButton, 
  Typography, 
  Button,
  Divider,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  useTheme,
  Avatar,
  Tooltip,
  Badge,
  Tabs,
  Tab
} from '@mui/material';
import { APP_NAME } from '../constants/app';
import { 
  ArrowBack as ArrowBackIcon,
  Logout as LogoutIcon,
  ChevronLeft as ChevronLeftIcon,
  ChevronRight as ChevronRightIcon,
  FolderOutlined as ProjectsIcon,
  BugReportOutlined as IssuesIcon,
  SettingsOutlined as SettingsIcon,
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  NotificationsNone as NotificationsIcon,
  AdminPanelSettings as AdminPanelSettingsIcon,
  InsightsOutlined as AnalyticsIcon,
  FlagOutlined as FlagOutlinedIcon,
  Schedule as ScheduleIcon,
  PeopleOutline as PeopleIcon
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';
import ThemeToggle from './ThemeToggle';
import Breadcrumbs from './Breadcrumbs';
import SkipLink from './SkipLink';
import WardenLogo from "./WardenLogo.tsx";

interface LayoutProps {
  children: ReactNode;
  showBackButton?: boolean;
  backTo?: string;
}

// Drawer width constants
const DRAWER_WIDTH = 260;
const DRAWER_COLLAPSED_WIDTH = 72;

// AppBar heights
const APPBAR_MAIN_HEIGHT = 70;
const TABBAR_HEIGHT = 48;

const Layout: React.FC<LayoutProps> = ({ 
  children, 
  showBackButton = false, 
  backTo = '/dashboard' 
}) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const [open, setOpen] = useState(() => {
    // Get saved state from localStorage, default to true (expanded)
    try {
      const saved = localStorage.getItem('sidebarExpanded');
      return saved !== null ? JSON.parse(saved) : true;
    } catch (error) {
      // Fallback to true if localStorage is not available (e.g., private mode)
      console.warn('localStorage not available, using default sidebar state');
      return true;
    }
  });
  const [notificationAnchorEl, setNotificationAnchorEl] = useState<null | HTMLElement>(null);

  const handleDrawerToggle = () => {
    const newOpen = !open;
    setOpen(newOpen);
    // Save state to localStorage
    try {
      localStorage.setItem('sidebarExpanded', JSON.stringify(newOpen));
    } catch (error) {
      // Ignore errors if localStorage is not available
      console.warn('Could not save sidebar state to localStorage');
    }
  };

  const handleNotificationOpen = (event: React.MouseEvent<HTMLElement>) => {
    setNotificationAnchorEl(event.currentTarget);
  };

  const handleNotificationClose = () => {
    setNotificationAnchorEl(null);
  };

  const handleViewAllNotifications = () => {
    navigate('/notifications');
  };

  // Define menu items based on user role
  const getMenuItems = () => {
    const items = [
      { text: 'Dashboard', icon: <DashboardIcon />, path: '/dashboard' },
      { text: 'Projects', icon: <ProjectsIcon />, path: '/projects' },
      { text: 'Analytics', icon: <AnalyticsIcon />, path: '/analytics' },
      { text: 'Settings', icon: <SettingsIcon />, path: '/settings' }
    ];

    // Add Admin menu item for superusers
    if (user?.is_superuser) {
      items.push({ text: 'Admin', icon: <AdminPanelSettingsIcon />, path: '/admin' });
    }

    return items;
  };

  const menuItems = getMenuItems();

    // Detect if we are on a project details page and extract projectId
    const projectMatch = location.pathname.match(/^\/projects\/([^/]+)/);
    const currentProjectId = projectMatch ? projectMatch[1] : null;
    const showProjectSidebar = Boolean(currentProjectId);

  // Check if the current path matches the menu item path
  const isActive = (path: string) => {
    return location.pathname === path || location.pathname.startsWith(`${path}/`);
  };

  // Get user initials for avatar
  const getUserInitials = () => {
    if (!user?.username) return 'U';
    const names = user.username.split(' ');
    if (names.length === 1) return names[0].charAt(0).toUpperCase();
    return (names[0].charAt(0) + names[names.length - 1].charAt(0)).toUpperCase();
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'row', minHeight: '100vh' }}>
      <AppBar 
        position="fixed" 
        elevation={0}
        sx={{ 
          width: '100%',
          zIndex: theme.zIndex.drawer + 1,
          transition: theme.transitions.create(['width'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
          }),
        }}
      >
        <Toolbar sx={{ height: 70, px: { xs: 2, sm: 3 } }}>
          <IconButton
            aria-label="toggle drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ 
              mr: 2, 
              display: showProjectSidebar ? { xs: 'none', sm: 'flex' } : 'none',
              '&:hover': {
                backgroundColor: theme.palette.mode === 'dark' 
                  ? 'rgba(255, 255, 255, 0.08)' 
                  : 'rgba(130, 82, 255, 0.08)',
              },
            }}
          >
            <MenuIcon className="gradient-text" />
          </IconButton>

          {showBackButton && (
            <IconButton
              size="medium"
              aria-label="back"
              sx={{ 
                mr: 2,
                '&:hover': {
                  backgroundColor: theme.palette.mode === 'dark' 
                    ? 'rgba(255, 255, 255, 0.08)' 
                    : 'rgba(130, 82, 255, 0.08)',
                },
              }}
              onClick={() => navigate(backTo)}
            >
              <ArrowBackIcon className="gradient-text" />
            </IconButton>
          )}

          <WardenLogo logoSize={32} />


          <Box sx={{ flexGrow: 1 }} />

          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
            <Tooltip title="Notifications">
              <IconButton
                size="medium"
                aria-label="notifications"
                onClick={handleNotificationOpen}
                sx={{ 
                  '&:hover': {
                    backgroundColor: theme.palette.mode === 'dark' 
                      ? 'rgba(255, 255, 255, 0.08)' 
                      : 'rgba(130, 82, 255, 0.08)',
                  },
                }}
              >
                <Badge badgeContent={12} color="error">
                  <NotificationsIcon className="gradient-text" />
                </Badge>
              </IconButton>
            </Tooltip>


            <ThemeToggle />

            <Divider orientation="vertical" flexItem sx={{ mx: 0.5, my: 1 }} />

            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
              <Tooltip title={user?.username || 'User'}>
                <Avatar 
                  sx={{ 
                    width: 38, 
                    height: 38, 
                    bgcolor: 'primary.main',
                    boxShadow: '0 2px 8px rgba(130, 82, 255, 0.3)',
                    fontWeight: 600,
                    fontSize: '0.9rem',
                  }}
                >
                  {getUserInitials()}
                </Avatar>
              </Tooltip>

              <Box sx={{ display: { xs: 'none', md: 'block' } }}>
                <Typography 
                  variant="body1" 
                  sx={{ 
                    fontWeight: 500,
                    color: theme.palette.mode === 'dark' ? 'inherit' : 'text.primary',
                  }}
                >
                  {user?.username || 'User'}
                </Typography>
              </Box>

              <Button 
                variant="outlined"
                onClick={logout}
                startIcon={<LogoutIcon />}
                size="small"
                sx={{
                  borderColor: theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.2)' : 'rgba(130, 82, 255, 0.3)',
                  color: theme.palette.mode === 'dark' ? 'inherit' : 'primary.main',
                  background: 'transparent',
                  '&:hover': {
                    borderColor: theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.3)' : 'primary.main',
                    background: theme.palette.mode === 'dark' 
                      ? 'rgba(255, 255, 255, 0.05)' 
                      : 'rgba(130, 82, 255, 0.05)',
                  },
                }}
              >
                Logout
              </Button>
            </Box>
          </Box>
        </Toolbar>

        {/* Secondary navigation tabs (below main toolbar) */}
        <Box sx={{ display: { xs: 'none', md: 'block' } }}>
          <Divider sx={{ opacity: theme.palette.mode === 'dark' ? 0.12 : 0.08 }} />
          <Toolbar variant="dense" sx={{ minHeight: TABBAR_HEIGHT, px: { xs: 2, sm: 3 } }}>
            <Tabs
              value={menuItems.find((it) => isActive(it.path))?.path ?? false}
              onChange={(_e, val) => { if (typeof val === 'string') navigate(val); }}
              textColor="inherit"
              indicatorColor="primary"
              aria-label="Primary navigation tabs"
              sx={{
                minHeight: TABBAR_HEIGHT,
                '& .MuiTab-root': {
                  minHeight: TABBAR_HEIGHT,
                  textTransform: 'none',
                  fontWeight: 600,
                  mr: 0,
                  px: 2,
                  borderRadius: 0,
                  backgroundColor: 'transparent',
                  color: theme.palette.mode === 'dark' ? theme.palette.common.white : theme.palette.primary.main,
                  opacity: 1,
                  transition: 'color 0.2s',
                  '&:hover': {
                    color: theme.palette.mode === 'dark' ? '#fff' : theme.palette.primary.dark,
                    backgroundColor: 'transparent',
                  },
                },
                '& .MuiTab-root.Mui-selected': {
                  color: theme.palette.mode === 'dark' ? theme.palette.common.white : theme.palette.primary.main,
                  fontWeight: 700,
                },
              }}
            >
              {menuItems.map((item) => (
                <Tab
                  key={item.text}
                  value={item.path}
                  icon={item.icon}
                  iconPosition="start"
                  label={item.text}
                />
              ))}
            </Tabs>
          </Toolbar>
        </Box>
      </AppBar>

      {showProjectSidebar && (
        <Drawer
          variant="permanent"
          sx={{
            width: open ? DRAWER_WIDTH : DRAWER_COLLAPSED_WIDTH,
            flexShrink: 0,
            '& .MuiDrawer-paper': {
              width: open ? DRAWER_WIDTH : DRAWER_COLLAPSED_WIDTH,
              boxSizing: 'border-box',
              transition: theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
              }),
              overflowX: 'hidden',
              marginTop: { xs: `${APPBAR_MAIN_HEIGHT}px`, md: `${APPBAR_MAIN_HEIGHT + TABBAR_HEIGHT}px` },
              borderRight: 'none',
              boxShadow: theme.palette.mode === 'dark'
                ? '2px 0 10px rgba(0, 0, 0, 0.2)'
                : '2px 0 10px rgba(0, 0, 0, 0.05)',
            },
          }}
        >
          <Box sx={{ py: 2 }}>
            <List sx={{ px: 1.5 }}>
              <ListItem disablePadding sx={{ display: 'block', mb: 0.8 }}>
                <ListItemButton
                  sx={{
                    minHeight: 48,
                    justifyContent: open ? 'initial' : 'center',
                    px: 2.5,
                    py: 1.2,
                    borderRadius: 2,
                    backgroundColor: location.pathname === `/projects/${currentProjectId}` ? (
                      theme.palette.mode === 'dark' ? 'rgba(130, 82, 255, 0.15)' : 'rgba(130, 82, 255, 0.1)'
                    ) : 'transparent',
                    '&:hover': {
                      backgroundColor: theme.palette.mode === 'dark' 
                        ? 'rgba(255, 255, 255, 0.06)'
                        : 'rgba(130, 82, 255, 0.06)',
                    },
                  }}
                  onClick={() => {
                    const el = document.getElementById('features');
                    if (el) {
                      const y = el.getBoundingClientRect().top + window.scrollY - 80;
                      window.scrollTo({ top: y, behavior: 'smooth' });
                    } else if (currentProjectId) {
                      navigate(`/projects/${currentProjectId}#features`);
                    }
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      mr: open ? 3 : 'auto',
                      justifyContent: 'center',
                      color: location.pathname === `/projects/${currentProjectId}` ? 'primary.main' : 'inherit',
                    }}
                  >
                    <FlagOutlinedIcon />
                  </ListItemIcon>
                  <ListItemText 
                    primary={'Features'} 
                    primaryTypographyProps={{
                      fontWeight: location.pathname === `/projects/${currentProjectId}` ? 600 : 500,
                      color: location.pathname === `/projects/${currentProjectId}` ? 'primary.main' : 'inherit',
                    }}
                    sx={{ 
                      opacity: open ? 1 : 0,
                      ml: 0.5,
                    }} 
                  />
                </ListItemButton>
              </ListItem>

              {/* Segments menu item */}
              <ListItem disablePadding sx={{ display: 'block', mb: 0.8 }}>
                <ListItemButton
                  sx={{
                    minHeight: 48,
                    justifyContent: open ? 'initial' : 'center',
                    px: 2.5,
                    py: 1.2,
                    borderRadius: 2,
                    backgroundColor: location.pathname.startsWith(`/projects/${currentProjectId}/segments`) ? (
                      theme.palette.mode === 'dark' ? 'rgba(130, 82, 255, 0.15)' : 'rgba(130, 82, 255, 0.1)'
                    ) : 'transparent',
                    '&:hover': {
                      backgroundColor: theme.palette.mode === 'dark' 
                        ? 'rgba(255, 255, 255, 0.06)' 
                        : 'rgba(130, 82, 255, 0.06)',
                    },
                  }}
                  onClick={() => {
                    if (currentProjectId) navigate(`/projects/${currentProjectId}/segments`);
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      mr: open ? 3 : 'auto',
                      justifyContent: 'center',
                      color: location.pathname.startsWith(`/projects/${currentProjectId}/segments`) ? 'primary.main' : 'inherit',
                    }}
                  >
                    <PeopleIcon />
                  </ListItemIcon>
                  <ListItemText 
                    primary={'Segments'} 
                    primaryTypographyProps={{
                      fontWeight: location.pathname.startsWith(`/projects/${currentProjectId}/segments`) ? 600 : 500,
                      color: location.pathname.startsWith(`/projects/${currentProjectId}/segments`) ? 'primary.main' : 'inherit',
                    }}
                    sx={{ 
                      opacity: open ? 1 : 0,
                      ml: 0.5,
                    }} 
                  />
                </ListItemButton>
              </ListItem>

              {/* Scheduling menu item */}
              <ListItem disablePadding sx={{ display: 'block', mb: 0.8 }}>
                <ListItemButton
                  sx={{
                    minHeight: 48,
                    justifyContent: open ? 'initial' : 'center',
                    px: 2.5,
                    py: 1.2,
                    borderRadius: 2,
                    backgroundColor: location.pathname.startsWith(`/projects/${currentProjectId}/scheduling`) ? (
                      theme.palette.mode === 'dark' ? 'rgba(130, 82, 255, 0.15)' : 'rgba(130, 82, 255, 0.1)'
                    ) : 'transparent',
                    '&:hover': {
                      backgroundColor: theme.palette.mode === 'dark' 
                        ? 'rgba(255, 255, 255, 0.06)'
                        : 'rgba(130, 82, 255, 0.06)',
                    },
                  }}
                  onClick={() => {
                    if (currentProjectId) navigate(`/projects/${currentProjectId}/scheduling`);
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      mr: open ? 3 : 'auto',
                      justifyContent: 'center',
                      color: location.pathname.startsWith(`/projects/${currentProjectId}/scheduling`) ? 'primary.main' : 'inherit',
                    }}
                  >
                    <ScheduleIcon />
                  </ListItemIcon>
                  <ListItemText 
                    primary={'Scheduling'} 
                    primaryTypographyProps={{
                      fontWeight: location.pathname.startsWith(`/projects/${currentProjectId}/scheduling`) ? 600 : 500,
                      color: location.pathname.startsWith(`/projects/${currentProjectId}/scheduling`) ? 'primary.main' : 'inherit',
                    }}
                    sx={{ 
                      opacity: open ? 1 : 0,
                      ml: 0.5,
                    }} 
                  />
                </ListItemButton>
              </ListItem>

              {/* Project Settings menu item */}
              <ListItem disablePadding sx={{ display: 'block', mb: 0.8 }}>
                <ListItemButton
                  sx={{
                    minHeight: 48,
                    justifyContent: open ? 'initial' : 'center',
                    px: 2.5,
                    py: 1.2,
                    borderRadius: 2,
                    backgroundColor: location.pathname.startsWith(`/projects/${currentProjectId}/settings`) ? (
                      theme.palette.mode === 'dark' ? 'rgba(130, 82, 255, 0.15)' : 'rgba(130, 82, 255, 0.1)'
                    ) : 'transparent',
                    '&:hover': {
                      backgroundColor: theme.palette.mode === 'dark' 
                        ? 'rgba(255, 255, 255, 0.06)'
                        : 'rgba(130, 82, 255, 0.06)',
                    },
                  }}
                  onClick={() => {
                    if (currentProjectId) navigate(`/projects/${currentProjectId}/settings`);
                  }}
                >
                  <ListItemIcon
                    sx={{
                      minWidth: 0,
                      mr: open ? 3 : 'auto',
                      justifyContent: 'center',
                      color: location.pathname.startsWith(`/projects/${currentProjectId}/settings`) ? 'primary.main' : 'inherit',
                    }}
                  >
                    <SettingsIcon />
                  </ListItemIcon>
                  <ListItemText 
                    primary={'Project Settings'} 
                    primaryTypographyProps={{
                      fontWeight: location.pathname.startsWith(`/projects/${currentProjectId}/settings`) ? 600 : 500,
                      color: location.pathname.startsWith(`/projects/${currentProjectId}/settings`) ? 'primary.main' : 'inherit',
                    }}
                    sx={{ 
                      opacity: open ? 1 : 0,
                      ml: 0.5,
                    }} 
                  />
                </ListItemButton>
              </ListItem>
            </List>
            <Divider sx={{ my: 2, mx: 2 }} />
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 1 }}>
              <Tooltip title={open ? "Collapse sidebar" : "Expand sidebar"}>
                <IconButton 
                  onClick={handleDrawerToggle}
                  sx={{
                    backgroundColor: theme.palette.mode === 'dark' 
                      ? 'rgba(255, 255, 255, 0.05)' 
                      : 'rgba(130, 82, 255, 0.05)',
                    '&:hover': {
                      backgroundColor: theme.palette.mode === 'dark' 
                        ? 'rgba(255, 255, 255, 0.1)' 
                        : 'rgba(130, 82, 255, 0.1)',
                    },
                  }}
                >
                  {open ? <ChevronLeftIcon /> : <ChevronRightIcon />}
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
        </Drawer>
      )}

      <Box
        component="main"
        sx={{ 
          flexGrow: 1, 
          p: { xs: 2, sm: 3, md: 4 },
          ml: 0,
          mr: '0 !important',
          right: 0,
          width: showProjectSidebar ? { xs: '100%', sm: `calc(100% - ${open ? DRAWER_WIDTH : DRAWER_COLLAPSED_WIDTH}px)` } : '100%',
          boxSizing: 'border-box',
          transition: theme.transitions.create(['width'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
          }),
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <SkipLink href="#main-content">Skip to main content</SkipLink>
        <Box sx={{ height: { xs: `${APPBAR_MAIN_HEIGHT}px`, md: `${APPBAR_MAIN_HEIGHT + TABBAR_HEIGHT}px` } }} />
        <Box 
          id="main-content"
          sx={{ 
            mt: 2,
            mb: 4, 
            width: '100%', 
            maxWidth: '1400px', 
          }}
        >
          <Breadcrumbs />
          {children}
        </Box>

        <Box
          component="footer"
          sx={{
            py: 3,
            px: 2,
            mt: 'auto',
            width: '100%',
            maxWidth: '1400px',
            backgroundColor: 'transparent',
            borderTop: `1px solid ${theme.palette.mode === 'dark' 
              ? 'rgba(255, 255, 255, 0.1)' 
              : 'rgba(0, 0, 0, 0.08)'}`,
          }}
        >
          <Box sx={{ width: '100%' }}>
            <Typography variant="body2" color="text.secondary" align="center">
              {APP_NAME} © {new Date().getFullYear()}
            </Typography>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default Layout;
