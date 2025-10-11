import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  CircularProgress,
  Divider,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  Collapse,
  Card,
  CardContent,
  CardActions,
  Chip,
  Tabs,
  Tab,
  Tooltip
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  Edit as EditIcon,
  Refresh as RefreshIcon,
  Close as CloseIcon,
  Send as SendIcon
} from '@mui/icons-material';
import apiClient from '../api/apiClient';
import type {
  NotificationSetting,
  CreateNotificationSettingRequest,
} from '../generated/api/client';

// Interface for our project data
interface Project {
  id: string;
  name: string;
  team_id?: number | null;
  team_name?: string | null;
  created_at: string;
}

interface NotificationSettingsProps {
  projectId: string;
  environmentKey: string;
  environments: any[];
  loadingEnvironments: boolean;
}

// Use the standard NotificationSetting interface
type ExtendedNotificationSetting = NotificationSetting;

interface EmailConfig {
  email_to: string;
}

interface MattermostConfig {
  webhook_url: string;
  channel_name: string;
}

interface WebhookConfig {
  webhook_url: string;
}

interface TelegramConfig {
  bot_token: string;
  chat_id: string;
}

interface SlackConfig {
  webhook_url: string;
  channel_name: string;
}

interface PachcaConfig {
  webhook_url: string;
}

const NotificationSettings: React.FC<NotificationSettingsProps> = ({ projectId, environmentKey, environments, loadingEnvironments }) => {
  const [settings, setSettings] = useState<ExtendedNotificationSetting[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [expandedSettings, setExpandedSettings] = useState<number[]>([]);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [openValidationDialog, setOpenValidationDialog] = useState(false);

  // State for field validation
  const [emailFieldError, setEmailFieldError] = useState(false);
  const [mattermostWebhookError, setMattermostWebhookError] = useState(false);
  const [mattermostChannelError, setMattermostChannelError] = useState(false);
  const [webhookUrlError, setWebhookUrlError] = useState(false);
  const [telegramBotTokenError, setTelegramBotTokenError] = useState(false);
  const [telegramChatIdError, setTelegramChatIdError] = useState(false);
  const [slackWebhookError, setSlackWebhookError] = useState(false);
  const [slackChannelError, setSlackChannelError] = useState(false);
  const [pachcaWebhookUrlError, setPachcaWebhookUrlError] = useState(false);
  const [project, setProject] = useState<Project>({
    id: projectId,
    name: `Project ${projectId}`,
    team_id: null,
    team_name: null,
    created_at: new Date().toISOString()
  });

  // State for add/edit setting dialog
  const [openSettingDialog, setOpenSettingDialog] = useState(false);
  const [settingDialogMode, setSettingDialogMode] = useState<'add' | 'edit'>('add');
  const [currentSetting, setCurrentSetting] = useState<ExtendedNotificationSetting | null>(null);
  const [settingType, setSettingType] = useState<string>('email');
  const [settingEnabled, setSettingEnabled] = useState(true);
  const [emailConfig, setEmailConfig] = useState<EmailConfig>({ email_to: '' });
  const [mattermostConfig, setMattermostConfig] = useState<MattermostConfig>({ webhook_url: '', channel_name: '' });
  const [webhookConfig, setWebhookConfig] = useState<WebhookConfig>({ webhook_url: '' });
  const [telegramConfig, setTelegramConfig] = useState<TelegramConfig>({ bot_token: '', chat_id: '' });
  const [slackConfig, setSlackConfig] = useState<SlackConfig>({ webhook_url: '', channel_name: '' });
  const [pachcaConfig, setPachcaConfig] = useState<PachcaConfig>({ webhook_url: '' });

  // State for tab selection
  const [selectedTab, setSelectedTab] = useState<string>('all');
  
  // State for environment selection
  const [currentEnvironmentKey, setCurrentEnvironmentKey] = useState<string>(environmentKey);


  // Fetch notification settings
  const fetchSettingsAndRules = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.listNotificationSettings(projectId, currentEnvironmentKey);
      const fetchedSettings = response.data.notification_settings;
      setSettings(fetchedSettings);
    } catch (err) {
      console.error('Error fetching notification settings:', err);
      setError('Failed to load notification settings');
    } finally {
      setLoading(false);
    }
  };

  // Fetch project details
  const fetchProjectDetails = async () => {
    try {
      const response = await apiClient.getProject(projectId);
      setProject(response.data.project);
    } catch (err) {
      console.error('Error fetching project details:', err);
    }
  };

  // Load data on component mount
  useEffect(() => {
    fetchProjectDetails();
    fetchSettingsAndRules();
  }, [projectId, currentEnvironmentKey]);

  const hasEmailChannel = () => {
    return settings.some(setting => setting.type === 'email');
  };

  const handleAddSetting = (channelType?: string) => {
    setSettingDialogMode('add');
    setCurrentSetting(null);
    setSettingType(channelType || 'email');
    setSettingEnabled(true);
    setEmailConfig({ email_to: '' });
    setMattermostConfig({ webhook_url: '', channel_name: '' });
    setWebhookConfig({ webhook_url: '' });
    setTelegramConfig({ bot_token: '', chat_id: '' });
    setSlackConfig({ webhook_url: '', channel_name: '' });
    setPachcaConfig({ webhook_url: '' });
    setOpenSettingDialog(true);
  };

  // Helper functions for each channel type
  const handleAddEmailSetting = () => handleAddSetting('email');
  const handleAddTelegramSetting = () => handleAddSetting('telegram');
  const handleAddMattermostSetting = () => handleAddSetting('mattermost');
  const handleAddSlackSetting = () => handleAddSetting('slack');
  const handleAddWebhookSetting = () => handleAddSetting('webhook');
  const handleAddPachcaSetting = () => handleAddSetting('pachca');

  const handleEditSetting = (setting: ExtendedNotificationSetting) => {
    setSettingDialogMode('edit');
    setCurrentSetting(setting);
    setSettingType(setting.type);
    setSettingEnabled(setting.enabled);

    // Parse the config based on the setting type
    try {
      const config = JSON.parse(setting.config);
      switch (setting.type) {
        case 'email':
          setEmailConfig(config);
          break;
        case 'mattermost':
          setMattermostConfig(config);
          break;
        case 'webhook':
          setWebhookConfig(config);
          break;
        case 'telegram':
          setTelegramConfig(config);
          break;
        case 'slack':
          setSlackConfig(config);
          break;
        case 'pachca':
          setPachcaConfig(config);
          break;
      }
    } catch (e) {
      console.error('Error parsing setting config:', e);
    }

    setOpenSettingDialog(true);
  };

  const handleDeleteSetting = async (settingId: number) => {
    if (!confirm('Are you sure you want to delete this notification setting?')) {
      return;
    }

    try {
      await apiClient.deleteNotificationSetting(projectId, currentEnvironmentKey, settingId);
      setSuccess('Notification setting deleted successfully');
      fetchSettingsAndRules(); // Refresh the list
    } catch (err) {
      console.error('Error deleting notification setting:', err);
      setError('Failed to delete notification setting');
    }
  };

  const handleSendTestNotification = async (settingId: number) => {
    try {
      await apiClient.sendTestNotification(projectId, currentEnvironmentKey, settingId);
      setSuccess('Test notification sent successfully');
    } catch (err) {
      console.error('Error sending test notification:', err);
      setError('Failed to send test notification');
    }
  };

  const handleSaveSetting = async () => {
    // Validate the configuration based on the setting type
    let config = {};
    let isValid = true;

    switch (settingType) {
      case 'email':
        if (!emailConfig.email_to || !isValidEmail(emailConfig.email_to)) {
          setEmailFieldError(true);
          isValid = false;
        } else {
          setEmailFieldError(false);
          config = emailConfig;
        }
        break;
      case 'mattermost':
        if (!mattermostConfig.webhook_url || !isValidUrl(mattermostConfig.webhook_url)) {
          setMattermostWebhookError(true);
          isValid = false;
        } else {
          setMattermostWebhookError(false);
        }
        if (!mattermostConfig.channel_name) {
          setMattermostChannelError(true);
          isValid = false;
        } else {
          setMattermostChannelError(false);
        }
        if (isValid) {
          config = mattermostConfig;
        }
        break;
      case 'webhook':
        if (!webhookConfig.webhook_url || !isValidUrl(webhookConfig.webhook_url)) {
          setWebhookUrlError(true);
          isValid = false;
        } else {
          setWebhookUrlError(false);
          config = webhookConfig;
        }
        break;
      case 'telegram':
        if (!telegramConfig.bot_token) {
          setTelegramBotTokenError(true);
          isValid = false;
        } else {
          setTelegramBotTokenError(false);
        }
        if (!telegramConfig.chat_id) {
          setTelegramChatIdError(true);
          isValid = false;
        } else {
          setTelegramChatIdError(false);
        }
        if (isValid) {
          config = telegramConfig;
        }
        break;
      case 'slack':
        if (!slackConfig.webhook_url || !isValidUrl(slackConfig.webhook_url)) {
          setSlackWebhookError(true);
          isValid = false;
        } else {
          setSlackWebhookError(false);
        }
        if (!slackConfig.channel_name) {
          setSlackChannelError(true);
          isValid = false;
        } else {
          setSlackChannelError(false);
        }
        if (isValid) {
          config = slackConfig;
        }
        break;
      case 'pachca':
        if (!pachcaConfig.webhook_url || !isValidUrl(pachcaConfig.webhook_url)) {
          setPachcaWebhookUrlError(true);
          isValid = false;
        } else {
          setPachcaWebhookUrlError(false);
          config = pachcaConfig;
        }
        break;
    }

    if (!isValid) {
      setValidationError('Please fix the validation errors above');
      setOpenValidationDialog(true);
      return;
    }

    try {
      const request: CreateNotificationSettingRequest = {
        type: settingType as any,
        config: JSON.stringify(config),
        enabled: settingEnabled
      };

      if (settingDialogMode === 'add') {
        await apiClient.createNotificationSetting(projectId, currentEnvironmentKey, request);
        setSuccess('Notification setting created successfully');
      } else {
        await apiClient.updateNotificationSetting(projectId, currentEnvironmentKey, currentSetting!.id, request);
        setSuccess('Notification setting updated successfully');
      }

      setOpenSettingDialog(false);
      fetchSettingsAndRules(); // Refresh the list
    } catch (err) {
      console.error('Error saving notification setting:', err);
      setError('Failed to save notification setting');
    }
  };

  const renderSettingConfigForm = () => {
    switch (settingType) {
      case 'email':
        return (
          <TextField
            fullWidth
            margin="normal"
            label="Email Address"
            value={emailConfig.email_to}
            onChange={(e) => setEmailConfig({ ...emailConfig, email_to: e.target.value })}
            error={emailFieldError}
            helperText={emailFieldError ? 'Please enter a valid email address' : ''}
          />
        );
      case 'mattermost':
        return (
          <>
            <TextField
              fullWidth
              margin="normal"
              label="Webhook URL"
              value={mattermostConfig.webhook_url}
              onChange={(e) => setMattermostConfig({ ...mattermostConfig, webhook_url: e.target.value })}
              error={mattermostWebhookError}
              helperText={mattermostWebhookError ? 'Please enter a valid webhook URL' : ''}
            />
            <TextField
              fullWidth
              margin="normal"
              label="Channel Name"
              value={mattermostConfig.channel_name}
              onChange={(e) => setMattermostConfig({ ...mattermostConfig, channel_name: e.target.value })}
              error={mattermostChannelError}
              helperText={mattermostChannelError ? 'Please enter a channel name' : ''}
            />
          </>
        );
      case 'webhook':
        return (
          <TextField
            fullWidth
            margin="normal"
            label="Webhook URL"
            value={webhookConfig.webhook_url}
            onChange={(e) => setWebhookConfig({ ...webhookConfig, webhook_url: e.target.value })}
            error={webhookUrlError}
            helperText={webhookUrlError ? 'Please enter a valid webhook URL' : ''}
          />
        );
      case 'telegram':
        return (
          <>
            <TextField
              fullWidth
              margin="normal"
              label="Bot Token"
              value={telegramConfig.bot_token}
              onChange={(e) => setTelegramConfig({ ...telegramConfig, bot_token: e.target.value })}
              error={telegramBotTokenError}
              helperText={telegramBotTokenError ? 'Please enter a bot token' : ''}
            />
            <TextField
              fullWidth
              margin="normal"
              label="Chat ID"
              value={telegramConfig.chat_id}
              onChange={(e) => setTelegramConfig({ ...telegramConfig, chat_id: e.target.value })}
              error={telegramChatIdError}
              helperText={telegramChatIdError ? 'Please enter a chat ID' : ''}
            />
          </>
        );
      case 'slack':
        return (
          <>
            <TextField
              fullWidth
              margin="normal"
              label="Webhook URL"
              value={slackConfig.webhook_url}
              onChange={(e) => setSlackConfig({ ...slackConfig, webhook_url: e.target.value })}
              error={slackWebhookError}
              helperText={slackWebhookError ? 'Please enter a valid webhook URL' : ''}
            />
            <TextField
              fullWidth
              margin="normal"
              label="Channel Name"
              value={slackConfig.channel_name}
              onChange={(e) => setSlackConfig({ ...slackConfig, channel_name: e.target.value })}
              error={slackChannelError}
              helperText={slackChannelError ? 'Please enter a channel name' : ''}
            />
          </>
        );
      case 'pachca':
        return (
          <TextField
            fullWidth
            margin="normal"
            label="Webhook URL"
            value={pachcaConfig.webhook_url}
            onChange={(e) => setPachcaConfig({ ...pachcaConfig, webhook_url: e.target.value })}
            error={pachcaWebhookUrlError}
            helperText={pachcaWebhookUrlError ? 'Please enter a valid webhook URL' : ''}
          />
        );
      default:
        return null;
    }
  };

  const isValidEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const isValidUrl = (url: string): boolean => {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  };

  const toggleExpandSetting = (settingId: number) => {
    setExpandedSettings(prev =>
      prev.includes(settingId)
        ? prev.filter(id => id !== settingId)
        : [...prev, settingId]
    );
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: string) => {
    setSelectedTab(newValue);
  };

  const handleEnvironmentChange = (newEnvironmentKey: string) => {
    setCurrentEnvironmentKey(newEnvironmentKey);
    localStorage.setItem('currentEnvironmentKey', newEnvironmentKey);
    // Find the environment ID and save it to localStorage
    const selectedEnv = environments.find((env: any) => env.key === newEnvironmentKey);
    if (selectedEnv) {
      localStorage.setItem('currentEnvId', selectedEnv.id.toString());
    }
  };

  return (
    <Box>
      {/* Success and error messages */}
      <Collapse in={!!success || !!error}>
        <Box sx={{ mb: 2 }}>
          {success && (
            <Alert
              severity="success"
              action={
                <IconButton
                  aria-label="close"
                  color="inherit"
                  size="small"
                  onClick={() => setSuccess(null)}
                >
                  <CloseIcon fontSize="inherit" />
                </IconButton>
              }
            >
              {success}
            </Alert>
          )}
          {error && (
            <Alert
              severity="error"
              action={
                <IconButton
                  aria-label="close"
                  color="inherit"
                  size="small"
                  onClick={() => setError(null)}
                >
                  <CloseIcon fontSize="inherit" />
                </IconButton>
              }
            >
              {error}
            </Alert>
          )}
        </Box>
      </Collapse>

      {/* Header with environment selector, refresh and add buttons */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="subtitle1">
          Configure notification channels for this project
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <FormControl size="small" sx={{ minWidth: 200 }}>
            <InputLabel>Environment</InputLabel>
            <Select
              value={currentEnvironmentKey}
              label="Environment"
              size="small"
              onChange={(e) => handleEnvironmentChange(e.target.value)}
              disabled={loadingEnvironments}
            >
              {environments.map((env: any) => (
                <MenuItem key={env.id} value={env.key} data-env-id={env.id}>
                  {env.name} ({env.key})
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Button
            size="small"
            startIcon={<RefreshIcon />}
            onClick={fetchSettingsAndRules}
            disabled={loading}
          >
            Refresh
          </Button>
          <Button
            size="small"
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => handleAddSetting()}
            disabled={loading}
          >
            Add Channel
          </Button>
        </Box>
      </Box>

      {/* Tabs for different notification channel types */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs
          value={selectedTab}
          onChange={handleTabChange}
          aria-label="notification channel tabs"
          variant="scrollable"
          scrollButtons="auto"
        >
          <Tab label="All" value="all" />
          <Tab label="Email" value="email" />
          <Tab label="Telegram" value="telegram" />
          <Tab label="Mattermost" value="mattermost" />
          <Tab label="Slack" value="slack" />
          <Tab label="Webhook" value="webhook" />
          <Tab label="Pachca" value="pachca" />
        </Tabs>
      </Box>

      {/* Loading indicator */}
      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', my: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {/* No settings message */}
      {!loading && (settings.length === 0 || settings.filter(setting => selectedTab === 'all' || setting.type === selectedTab).length === 0) && (
        <Paper sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="body1" sx={{ mb: 2 }}>
            {settings.length === 0
              ? 'No notification settings configured for this project.'
              : `No ${selectedTab !== 'all' ? selectedTab : ''} notification settings configured for this project.`}
          </Typography>
          <Button
            size="small"
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              switch (selectedTab) {
                case 'email':
                  handleAddEmailSetting();
                  break;
                case 'telegram':
                  handleAddTelegramSetting();
                  break;
                case 'mattermost':
                  handleAddMattermostSetting();
                  break;
                case 'slack':
                  handleAddSlackSetting();
                  break;
                case 'webhook':
                  handleAddWebhookSetting();
                  break;
                case 'pachca':
                  handleAddPachcaSetting();
                  break;
                default:
                  handleAddEmailSetting();
                  break;
              }
            }}
          >
            Add {selectedTab !== 'all' ? selectedTab.charAt(0).toUpperCase() + selectedTab.slice(1) : 'Email'} Channel
          </Button>
        </Paper>
      )}

      {/* Settings list */}
      {!loading && settings.length > 0 && (
        <Box>
          {settings
            .filter(setting => selectedTab === 'all' || setting.type === selectedTab)
            .map(setting => (
              <Paper key={setting.id} sx={{ p: 3, mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', flex: 1 }}>
                    <Box>
                      <Typography variant="h6">
                        {setting.type.charAt(0).toUpperCase() + setting.type.slice(1)} Notifications #{setting.id}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {setting.enabled ? 'Enabled' : 'Disabled'}
                      </Typography>
                    </Box>
                    <Box sx={{ ml: 4, flex: 1 }}>
                      {setting.type === 'email' && (
                        project.team_id ? (
                          <Typography variant="body2">
                            Notifications are sent to all project team members.
                          </Typography>
                        ) : (
                          <Typography variant="body2">
                            Email: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as EmailConfig;
                                return config.email_to;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        )
                      )}
                      {setting.type === 'mattermost' && (
                        <Typography variant="body2">
                          Mattermost: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as MattermostConfig;
                              return config.channel_name;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                      {setting.type === 'webhook' && (
                        <Typography variant="body2">
                          Webhook configured
                        </Typography>
                      )}
                      {setting.type === 'telegram' && (
                        <Typography variant="body2">
                          Telegram bot configured
                        </Typography>
                      )}
                      {setting.type === 'slack' && (
                        <Typography variant="body2">
                          Slack: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as SlackConfig;
                              return config.channel_name;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                      {setting.type === 'pachca' && (
                        <Typography variant="body2">
                          Pachca bot configured
                        </Typography>
                      )}
                    </Box>
                  </Box>
                  <Box>
                    <Button
                      size="small"
                      variant="outlined"
                      onClick={() => toggleExpandSetting(setting.id)}
                      sx={{ mr: 1 }}
                    >
                      {expandedSettings.includes(setting.id) ? 'Hide Details' : 'Show Details'}
                    </Button>
                    <IconButton
                      aria-label="edit"
                      onClick={() => handleEditSetting(setting)}
                      size="small"
                      sx={{ mr: 1 }}
                    >
                      <EditIcon />
                    </IconButton>
                    <Tooltip title="Send test notification">
                      <IconButton
                        aria-label="send test notification"
                        onClick={() => handleSendTestNotification(setting.id)}
                        size="small"
                        sx={{ mr: 1 }}
                        color="primary"
                      >
                        <SendIcon />
                      </IconButton>
                    </Tooltip>
                    <IconButton
                      aria-label="delete"
                      onClick={() => handleDeleteSetting(setting.id)}
                      size="small"
                      color="error"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                </Box>

                <Collapse in={expandedSettings.includes(setting.id)} timeout="auto" unmountOnExit>
                  <Divider sx={{ my: 2 }} />

                  <Box sx={{ mt: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Configuration Details:
                    </Typography>
                    <Box sx={{ pl: 2 }}>
                      {setting.type === 'email' && (
                        project.team_id ? (
                          <Typography variant="body2">
                            Notifications are sent to all project team members.
                          </Typography>
                        ) : (
                          <Typography variant="body2">
                            Email: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as EmailConfig;
                                return config.email_to;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        )
                      )}
                      {setting.type === 'mattermost' && (
                        <>
                          <Typography variant="body2">
                            Webhook URL: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as MattermostConfig;
                                return config.webhook_url;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                          <Typography variant="body2">
                            Channel: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as MattermostConfig;
                                return config.channel_name;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        </>
                      )}
                      {setting.type === 'webhook' && (
                        <Typography variant="body2">
                          Webhook URL: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as WebhookConfig;
                              return config.webhook_url;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                      {setting.type === 'telegram' && (
                        <>
                          <Typography variant="body2">
                            Bot Token: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as TelegramConfig;
                                return config.bot_token;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                          <Typography variant="body2">
                            Chat ID: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as TelegramConfig;
                                return config.chat_id;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        </>
                      )}
                      {setting.type === 'slack' && (
                        <>
                          <Typography variant="body2">
                            Webhook URL: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as SlackConfig;
                                return config.webhook_url;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                          <Typography variant="body2">
                            Channel: {
                            (() => {
                              try {
                                const config = JSON.parse(setting.config) as SlackConfig;
                                return config.channel_name;
                              } catch (e) {
                                return 'Invalid configuration';
                              }
                            })()
                          }
                          </Typography>
                        </>
                      )}
                      {setting.type === 'pachca' && (
                        <Typography variant="body2">
                          Webhook URL: {
                          (() => {
                            try {
                              const config = JSON.parse(setting.config) as PachcaConfig;
                              return config.webhook_url;
                            } catch (e) {
                              return 'Invalid configuration';
                            }
                          })()
                        }
                        </Typography>
                      )}
                    </Box>
                  </Box>
                </Collapse>
              </Paper>
            ))}
        </Box>
      )}

      {/* Add/Edit Setting Dialog */}
      <Dialog open={openSettingDialog} onClose={() => setOpenSettingDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle className="gradient-text-purple">
          {settingDialogMode === 'add' ? 'Add Notification Channel' : 'Edit Notification Channel'}
        </DialogTitle>
        <DialogContent>
          <FormControl fullWidth margin="normal">
            <InputLabel id="setting-type-label">Notification Type</InputLabel>
            <Select
              labelId="setting-type-label"
              value={settingType}
              label="Notification Type"
              onChange={(e) => setSettingType(e.target.value)}
              disabled={settingDialogMode === 'edit'}
            >
              {/* Only show email option if we're editing an existing email channel or if no email channel exists */}
              {(settingDialogMode === 'edit' && currentSetting?.type === 'email') || (settingDialogMode === 'add' && !hasEmailChannel()) ? (
                <MenuItem value="email">Email</MenuItem>
              ) : null}
              <MenuItem value="mattermost">Mattermost</MenuItem>
              <MenuItem value="webhook">Webhook</MenuItem>
              <MenuItem value="telegram">Telegram</MenuItem>
              <MenuItem value="slack">Slack</MenuItem>
              <MenuItem value="pachca">Pachca</MenuItem>
            </Select>
          </FormControl>

          {renderSettingConfigForm()}

          <FormControlLabel
            control={
              <Switch
                checked={settingEnabled}
                onChange={(e) => setSettingEnabled(e.target.checked)}
              />
            }
            label="Enabled"
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button size="small" onClick={() => setOpenSettingDialog(false)}>Cancel</Button>
          <Button
            size="small"
            onClick={handleSaveSetting}
            variant="contained"
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Validation Error Dialog */}
      <Dialog open={openValidationDialog} onClose={() => setOpenValidationDialog(false)}>
        <DialogTitle className="gradient-text-purple">Validation Error</DialogTitle>
        <DialogContent>
          <Typography>{validationError}</Typography>
        </DialogContent>
        <DialogActions>
          <Button size="small" onClick={() => setOpenValidationDialog(false)} color="primary">
            OK
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default NotificationSettings;
