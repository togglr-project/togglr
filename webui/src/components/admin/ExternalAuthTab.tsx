import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
} from '@mui/material';
import {
  Sync as SyncIcon,
  Settings as SettingsIcon,
  ListAlt as LogsIcon,
} from '@mui/icons-material';
import TabPanel from './TabPanel';
import LDAPSyncTab from './ldap/LDAPSyncTab';
import LDAPConfigTab from './ldap/LDAPConfigTab';
import LDAPLogsTab from './ldap/LDAPLogsTab';

const ExternalAuthTab: React.FC = () => {
  const [tabValue, setTabValue] = useState(0);

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Box sx={{ mb: 3 }}>
        <Typography
          variant="h5"
          component="h2"
          gutterBottom
          sx={{
            fontWeight: 600,
            background: (theme) => theme.palette.mode === 'dark'
              ? 'linear-gradient(45deg, #8352ff 10%, #5e72e4 90%)'
              : 'linear-gradient(45deg, #5e72e4 30%, #8352ff 90%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            mb: 1
          }}
        >
          External Authentication
        </Typography>
        <Typography
          variant="body1"
          sx={{
            color: 'primary.light',
            maxWidth: '800px',
            fontSize: '1rem'
          }}
        >
          Configure and manage external authentication providers like LDAP/Active Directory.
        </Typography>
      </Box>

      <Paper
        sx={{
          width: '100%',
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(to bottom, rgba(65, 68, 74, 0.3), rgba(55, 58, 64, 0.3))'
            : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.7), rgba(245, 245, 245, 0.7))',
          backdropFilter: 'blur(10px)',
          boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.03)',
          borderRadius: 2
        }}
      >
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            aria-label="external auth tabs"
            sx={{
              '& .MuiTab-root': {
                fontWeight: 500,
                transition: 'all 0.2s ease-in-out',
                '&:hover': {
                  color: 'primary.main',
                  opacity: 0.8
                }
              },
              '& .Mui-selected': {
                fontWeight: 600
              }
            }}
          >
            <Tab 
              label="LDAP Sync" 
              icon={<SyncIcon />} 
              iconPosition="start"
            />
            <Tab 
              label="LDAP Config" 
              icon={<SettingsIcon />} 
              iconPosition="start"
            />
            <Tab 
              label="LDAP Logs" 
              icon={<LogsIcon />} 
              iconPosition="start"
            />
          </Tabs>
        </Box>

        {/* LDAP Sync Tab */}
        <TabPanel value={tabValue} index={0}>
          <LDAPSyncTab />
        </TabPanel>

        {/* LDAP Config Tab */}
        <TabPanel value={tabValue} index={1}>
          <LDAPConfigTab />
        </TabPanel>

        {/* LDAP Logs Tab */}
        <TabPanel value={tabValue} index={2}>
          <LDAPLogsTab />
        </TabPanel>
      </Paper>
    </Box>
  );
};

export default ExternalAuthTab; 