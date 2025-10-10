import React from 'react';
import {
  Popover,
  Typography,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  Divider,
  Button,
  Box,
  CircularProgress,
  Alert,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  Notifications as NotificationsIcon,
  CheckCircle as CheckCircleIcon,
  Group as GroupIcon,
  ArrowForward as ArrowForwardIcon
} from '@mui/icons-material';
import { useNotifications } from '../hooks/useNotifications';
import { formatDistanceToNow } from 'date-fns';
import type { UserNotification } from '../generated/api/client';

interface NotificationPopoverProps {
  anchorEl: HTMLElement | null;
  onClose: () => void;
  onViewAll: () => void;
}

const getNotificationIcon = (type: string) => {
  switch (type) {
    case 'project_added':
    case 'project_removed':
    case 'role_changed':
      return <GroupIcon />;
    default:
      return <NotificationsIcon />;
  }
};

const getNotificationTitle = (notification: UserNotification): string => {
  switch (notification.type) {
    case 'project_added':
      return 'You have been added to a project';
    case 'project_removed':
      return 'You have been removed from a project';
    case 'role_changed':
      return 'Your role has been changed';
    case 'need_approve':
      return 'You need to approve entity change';
    default:
      return 'Notification';
  }
};

const getNotificationMessage = (notification: UserNotification): string => {
  const content = notification.content;
  
  switch (notification.type) {
    case 'project_added':
      return `You have been added to project '${content.project_name}' as ${content.role}`;
    case 'project_removed':
      return `You have been removed from project '${content.project_name}'`;
    case 'role_changed':
      return `Your role in project '${content.project_name}' has been changed from ${content.old_role} to ${content.new_role}`;
    case 'need_approve':
      return `You need to approve entity change`;
    default:
      return 'You have a new notification';
  }
};

const NotificationPopover: React.FC<NotificationPopoverProps> = ({
  anchorEl,
  onClose,
  onViewAll
}) => {
  const { notifications, unreadCount, loading, error, markAsRead } = useNotifications();

  const handleMarkAsRead = async (notificationId: number) => {
    await markAsRead(notificationId);
  };

  const handleViewAll = () => {
    onViewAll();
    onClose();
  };

  return (
    <Popover
      open={Boolean(anchorEl)}
      anchorEl={anchorEl}
      onClose={onClose}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      PaperProps={{ 
        sx: { 
          minWidth: 400, 
          maxWidth: 500, 
          maxHeight: 600,
          overflow: 'hidden'
        } 
      }}
    >
      <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">
            Notifications
          </Typography>
          {unreadCount > 0 && (
            <Typography variant="caption" color="primary">
              {unreadCount} unread
            </Typography>
          )}
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ m: 2 }}>
          {error}
        </Alert>
      )}

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress size={24} />
        </Box>
      ) : (
        <>
          <List sx={{ p: 0, maxHeight: 400, overflow: 'auto' }}>
            {notifications.length === 0 ? (
              <ListItem>
                <ListItemText 
                  primary="No notifications" 
                  secondary="You're all caught up!"
                />
              </ListItem>
            ) : (
              notifications.map((notification, index) => (
                <React.Fragment key={notification.id}>
                  <ListItem 
                    sx={{ 
                      bgcolor: notification.is_read ? 'background.paper' : 'primary.50',
                      borderLeft: notification.is_read ? 'none' : '4px solid',
                      borderLeftColor: 'primary.main',
                      '&:hover': {
                        bgcolor: notification.is_read ? 'action.hover' : 'primary.100'
                      }
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar 
                        sx={{ 
                          bgcolor: notification.is_read ? 'grey.300' : 'primary.main',
                          width: 40,
                          height: 40,
                          opacity: notification.is_read ? 0.7 : 1
                        }}
                      >
                        {getNotificationIcon(notification.type)}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography 
                          variant="subtitle2" 
                          color={notification.is_read ? 'text.secondary' : 'text.primary'}
                          sx={{ fontWeight: notification.is_read ? 400 : 600 }}
                        >
                          {getNotificationTitle(notification)}
                        </Typography>
                      }
                      secondary={
                        <Box>
                          <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                            {getNotificationMessage(notification)}
                          </Typography>
                          <Typography variant="caption" color="text.disabled">
                            {formatDistanceToNow(new Date(notification.created_at), { addSuffix: true })}
                          </Typography>
                        </Box>
                      }
                    />
                    {!notification.is_read && (
                      <Tooltip title="Mark as read">
                        <IconButton
                          size="small"
                          onClick={() => handleMarkAsRead(notification.id)}
                          sx={{ ml: 1 }}
                        >
                          <CheckCircleIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    )}
                  </ListItem>
                  {index < notifications.length - 1 && <Divider />}
                </React.Fragment>
              ))
            )}
          </List>

          {notifications.length > 0 && (
            <Box sx={{ p: 2, borderTop: 1, borderColor: 'divider' }}>
              <Button
                fullWidth
                variant="outlined"
                endIcon={<ArrowForwardIcon />}
                onClick={handleViewAll}
              >
                View All Notifications
              </Button>
            </Box>
          )}
        </>
      )}
    </Popover>
  );
};

export default NotificationPopover;
