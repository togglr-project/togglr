import React, { useState } from 'react';
import {
  Box,
  Typography,
  Tabs,
  Tab,
  Alert,
  CircularProgress,
  Button,
  Paper,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { Refresh as RefreshIcon } from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { usePendingChanges } from '../hooks/usePendingChanges';
import PendingChangeCard from '../components/pending-changes/PendingChangeCard';
import apiClient from '../api/apiClient';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import { Assignment as ChangesIcon } from '@mui/icons-material';
import type { PendingChangeResponseStatusEnum } from '../generated/api/client';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`pending-tabpanel-${index}`}
      aria-labelledby={`pending-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

const PendingChangesPage: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const [tabValue, setTabValue] = useState(0);
  const [statusFilter, setStatusFilter] = useState<PendingChangeResponseStatusEnum | undefined>('pending');
  const [environmentKey, setEnvironmentKey] = useState<string>('prod');

  // Load environments for the project
  const { data: environmentsResp, isLoading: loadingEnvironments } = useQuery({
    queryKey: ['project-environments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectEnvironments(projectId || '');
      return res.data;
    },
    enabled: !!projectId,
  });
  const environments = environmentsResp?.items ?? [];
  const selectedEnv = environments.find((e: any) => e.key === environmentKey);
  const environmentId = selectedEnv?.id as number | undefined;

  const { data, isLoading, error, refetch } = usePendingChanges(projectId || '', statusFilter, environmentId);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
    switch (newValue) {
      case 0:
        setStatusFilter('pending');
        break;
      case 1:
        setStatusFilter('approved');
        break;
      case 2:
        setStatusFilter('rejected');
        break;
      case 3:
        setStatusFilter(undefined); // All
        break;
      default:
        setStatusFilter('pending');
    }
  };

  const handleRefresh = () => {
    refetch();
  };

  if (!projectId) {
    return (
      <AuthenticatedLayout>
        <PageHeader
          title="Change Requests"
          subtitle="Manage pending changes and approvals"
          icon={<ChangesIcon />}
        />
        <Box sx={{ p: 3 }}>
          <Alert severity="error">Project ID is required</Alert>
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (isLoading) {
    return (
      <AuthenticatedLayout>
        <PageHeader
          title="Change Requests"
          subtitle="Manage pending changes and approvals"
          icon={<ChangesIcon />}
        />
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (error) {
    return (
      <AuthenticatedLayout>
        <PageHeader
          title="Change Requests"
          subtitle="Manage pending changes and approvals"
          icon={<ChangesIcon />}
        />
        <Box sx={{ p: 3 }}>
          <Alert severity="error">
            Failed to load pending changes: {error.message}
          </Alert>
          <Button onClick={handleRefresh} startIcon={<RefreshIcon />} sx={{ mt: 2 }}>
            Retry
          </Button>
        </Box>
      </AuthenticatedLayout>
    );
  }

  const pendingChanges = data?.data || [];
  const pendingCount = pendingChanges.filter(pc => pc.status === 'pending').length;
  const approvedCount = pendingChanges.filter(pc => pc.status === 'approved').length;
  const rejectedCount = pendingChanges.filter(pc => pc.status === 'rejected').length;

  // Debug info (remove in production)
  // console.log('PendingChangesPage debug:', {
  //   projectId,
  //   statusFilter,
  //   totalChanges: pendingChanges.length,
  //   pendingCount,
  //   approvedCount,
  //   rejectedCount,
  //   changes: pendingChanges
  // });

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="Change Requests"
        subtitle="Manage pending changes and approvals"
        icon={<ChangesIcon />}
      />
      {/* Environment selector */}
      <Box sx={{ display: 'flex', justifyContent: 'flex-start', mb: 2 }}>
        <FormControl size="small" sx={{ minWidth: 220 }} disabled={loadingEnvironments}>
          <InputLabel>Environment</InputLabel>
          <Select
            label="Environment"
            value={environmentKey}
            onChange={(e) => setEnvironmentKey(e.target.value)}
          >
            {(environments || []).map((env: any) => (
              <MenuItem key={env.id} value={env.key}>
                {env.name} ({env.key})
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
      
      <Box sx={{ width: '100%' }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Tabs value={tabValue} onChange={handleTabChange} aria-label="pending changes tabs">
            <Tab 
              label={`Pending (${pendingCount})`} 
              id="pending-tab-0"
              aria-controls="pending-tabpanel-0"
            />
            <Tab 
              label={`Approved (${approvedCount})`} 
              id="pending-tab-1"
              aria-controls="pending-tabpanel-1"
            />
            <Tab 
              label={`Rejected (${rejectedCount})`} 
              id="pending-tab-2"
              aria-controls="pending-tabpanel-2"
            />
            <Tab 
              label="All" 
              id="pending-tab-3"
              aria-controls="pending-tabpanel-3"
            />
          </Tabs>
          <Button onClick={handleRefresh} startIcon={<RefreshIcon />} sx={{ mr: 2 }}>
            Refresh
          </Button>
        </Box>

      <TabPanel value={tabValue} index={0}>
        <Typography variant="h6" gutterBottom>
          Pending Changes ({pendingCount})
        </Typography>
        {pendingCount === 0 ? (
          <Paper sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary">
              No pending changes found
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              Pending changes will appear here when features with "guarded" tag are modified.
            </Typography>
          </Paper>
        ) : (
          pendingChanges
            .filter(pc => pc.status === 'pending')
            .map((pendingChange) => (
              <PendingChangeCard
                key={pendingChange.id}
                pendingChange={pendingChange}
                onStatusChange={handleRefresh}
              />
            ))
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={1}>
        <Typography variant="h6" gutterBottom>
          Approved Changes ({approvedCount})
        </Typography>
        {approvedCount === 0 ? (
          <Paper sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary">
              No approved changes found
            </Typography>
          </Paper>
        ) : (
          pendingChanges
            .filter(pc => pc.status === 'approved')
            .map((pendingChange) => (
              <PendingChangeCard
                key={pendingChange.id}
                pendingChange={pendingChange}
                onStatusChange={handleRefresh}
              />
            ))
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={2}>
        <Typography variant="h6" gutterBottom>
          Rejected Changes ({rejectedCount})
        </Typography>
        {rejectedCount === 0 ? (
          <Paper sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary">
              No rejected changes found
            </Typography>
          </Paper>
        ) : (
          pendingChanges
            .filter(pc => pc.status === 'rejected')
            .map((pendingChange) => (
              <PendingChangeCard
                key={pendingChange.id}
                pendingChange={pendingChange}
                onStatusChange={handleRefresh}
              />
            ))
        )}
      </TabPanel>

      <TabPanel value={tabValue} index={3}>
        <Typography variant="h6" gutterBottom>
          All Changes ({pendingChanges.length})
        </Typography>
        {pendingChanges.length === 0 ? (
          <Paper sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary">
              No changes found
            </Typography>
          </Paper>
        ) : (
          pendingChanges.map((pendingChange) => (
            <PendingChangeCard
              key={pendingChange.id}
              pendingChange={pendingChange}
              onStatusChange={handleRefresh}
            />
          ))
        )}
      </TabPanel>
      </Box>
    </AuthenticatedLayout>
  );
};

export default PendingChangesPage;
