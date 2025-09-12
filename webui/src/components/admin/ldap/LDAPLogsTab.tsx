import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Chip,
  TextField,
  Button,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  CircularProgress,
  IconButton,
  Tooltip,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../../../api/apiClient';
import { Visibility as VisibilityIcon } from '@mui/icons-material';

const LDAPLogsTab: React.FC = () => {
  // Filter state
  const [filters, setFilters] = useState({
    level: '',
    syncId: '',
    username: '',
    from: '',
    to: '',
    limit: 50
  });

  // Pagination state
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  // Fetch LDAP sync logs
  const {
    data: logsData,
    isLoading: isLoadingLogs,
    error: logsError,
    refetch: refetchLogs
  } = useQuery({
    queryKey: ['ldapSyncLogs', filters, page, rowsPerPage],
    queryFn: async () => {
      try {
        const response = await apiClient.getLDAPSyncLogs(
          filters.limit,
          filters.level ? filters.level as any : undefined,
          filters.syncId || undefined,
          filters.username || undefined,
          filters.from || undefined,
          filters.to || undefined
        );
        return response.data;
      } catch (error) {
        console.error('Error fetching LDAP sync logs:', error);
        return {
          logs: [],
          total: 0,
          has_more: false
        };
      }
    }
  });

  const handleFilterChange = (field: keyof typeof filters) => (
    event: React.ChangeEvent<HTMLInputElement> | { target: { value: any } }
  ) => {
    const value = event.target.value;
    setFilters(prev => ({
      ...prev,
      [field]: value
    }));
    setPage(0); // Reset to first page when filters change
  };

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const getLevelColor = (level: string) => {
    switch (level.toLowerCase()) {
      case 'error':
        return 'error';
      case 'warning':
        return 'warning';
      case 'info':
        return 'info';
      case 'debug':
        return 'default';
      default:
        return 'default';
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  if (isLoadingLogs) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (logsError) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Error loading LDAP sync logs. Please try again later.
      </Alert>
    );
  }

  const logs = logsData?.logs || [];
  const total = logsData?.total || 0;

  return (
    <>
      <Box 
        sx={{ 
          display: 'flex', 
          flexDirection: 'column',
          mb: 4,
          pb: 2,
          borderBottom: (theme) => `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)'}`
        }}
      >
        <Typography 
          variant="h6"
          sx={{ 
            fontWeight: 600,
            mb: 0.5
          }}
        >
          LDAP Sync Logs
        </Typography>
        <Typography 
          variant="body2" 
          color="text.secondary"
          sx={{ maxWidth: '600px' }}
        >
          View detailed logs of LDAP synchronization operations and troubleshoot issues.
        </Typography>
      </Box>

      {/* Filters */}
      <Paper 
        sx={{ 
          p: 3, 
          mb: 3,
          borderRadius: 2,
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
            : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
          backdropFilter: 'blur(8px)',
          boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <Typography variant="h6" sx={{ mb: 3, fontWeight: 600 }}>Filters</Typography>

        <Grid container spacing={3}>
          <Grid item xs={12} md={3}>
            <FormControl fullWidth>
              <InputLabel>Level</InputLabel>
              <Select
                value={filters.level}
                label="Level"
                onChange={handleFilterChange('level')}
              >
                <MenuItem value="">All Levels</MenuItem>
                <MenuItem value="error">Error</MenuItem>
                <MenuItem value="warning">Warning</MenuItem>
                <MenuItem value="info">Info</MenuItem>
                <MenuItem value="debug">Debug</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              fullWidth
              label="Sync ID"
              value={filters.syncId}
              onChange={handleFilterChange('syncId')}
              placeholder="Enter sync ID"
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              fullWidth
              label="Username"
              value={filters.username}
              onChange={handleFilterChange('username')}
              placeholder="Enter username"
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              fullWidth
              label="From Date"
              type="datetime-local"
              value={filters.from}
              onChange={handleFilterChange('from')}
              InputLabelProps={{
                shrink: true,
              }}
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              fullWidth
              label="To Date"
              type="datetime-local"
              value={filters.to}
              onChange={handleFilterChange('to')}
              InputLabelProps={{
                shrink: true,
              }}
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <TextField
              fullWidth
              label="Limit"
              type="number"
              value={filters.limit}
              onChange={handleFilterChange('limit')}
              inputProps={{ min: 1, max: 1000 }}
            />
          </Grid>
          <Grid item xs={12} md={3}>
            <Button
              variant="outlined"
              onClick={() => refetchLogs()}
              sx={{ mt: 1 }}
            >
              Refresh
            </Button>
          </Grid>
        </Grid>
      </Paper>

      {/* Logs Table */}
      <Paper 
        sx={{ 
          borderRadius: 2,
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
            : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
          backdropFilter: 'blur(8px)',
          boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Timestamp</TableCell>
                <TableCell>Level</TableCell>
                <TableCell>Message</TableCell>
                <TableCell>Username</TableCell>
                <TableCell>Group</TableCell>
                <TableCell>Sync ID</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {logs.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <Typography variant="body2" color="text.secondary">
                      No logs found
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                logs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell>
                      <Typography variant="body2">
                        {formatDate(log.timestamp || '')}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={log.level || 'unknown'}
                        color={getLevelColor(log.level || '')}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" sx={{ maxWidth: 300 }}>
                        {log.message || ''}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {log.username || '-'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" sx={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>
                        {log.sync_session_id || '-'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {log.details && (
                        <Tooltip title="View Details">
                          <IconButton size="small">
                            <VisibilityIcon />
                          </IconButton>
                        </Tooltip>
                      )}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        <TablePagination
          rowsPerPageOptions={[10, 25, 50, 100]}
          component="div"
          count={total}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </Paper>
    </>
  );
};

export default LDAPLogsTab; 