import React, { useMemo, useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  Box,
  Paper,
  Typography,
  CircularProgress,
  Alert,
  Table,
  TableHead,
  TableBody,
  TableCell,
  TableRow,
  TableContainer,
  Pagination as MuiPagination,
  Stack,
  Collapse,
  IconButton,
  Chip,
  Tooltip,
  Tabs,
  Tab,
} from '@mui/material';
import { ExpandMore as ExpandMoreIcon, ExpandLess as ExpandLessIcon } from '@mui/icons-material';
import JsonDiffViewer from '../components/audit-log/JsonDiffViewer';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import { Assignment as ChangesIcon } from '@mui/icons-material';
import apiClient from '../api/apiClient';
import { useRBAC } from '../auth/permissions';
import AuditLogFilter, { type AuditLogFilterValue } from '../components/audit-log/AuditLogFilter';

const DEFAULT_PER_PAGE = 20;

interface AuditLogRowProps {
  log: any;
}

const AuditLogRow: React.FC<AuditLogRowProps> = ({ log }) => {
  const [expanded, setExpanded] = useState(false);
  const [viewMode, setViewMode] = useState<'diff' | 'separate'>('diff');

  const formatValue = (value: any) => {
    if (!value) return null;
    try {
      return JSON.stringify(value, null, 2);
    } catch {
      return String(value);
    }
  };

  const hasDetails = log.old_value || log.new_value;
  const isUpdateAction = log.action === 'update' && log.old_value && log.new_value;

  return (
    <>
      <TableRow hover>
        <TableCell>{new Date(log.created_at).toLocaleString()}</TableCell>
        <TableCell>
          <Chip 
            label={log.environment_key ?? 'N/A'} 
            size="small" 
            variant="outlined"
            color="primary"
          />
        </TableCell>
        <TableCell>{log.actor}</TableCell>
        <TableCell>{log.username ?? '—'}</TableCell>
        <TableCell>
          <Chip 
            label={log.entity} 
            size="small" 
            color="secondary"
          />
        </TableCell>
        <TableCell>
          <Chip 
            label={log.action} 
            size="small" 
            color={log.action === 'create' ? 'success' : log.action === 'delete' ? 'error' : 'default'}
          />
        </TableCell>
        <TableCell>
          {log.old_value ? (
            <Tooltip title="Click to expand details">
              <Box 
                sx={{ 
                  maxWidth: 200, 
                  overflow: 'hidden', 
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                  cursor: hasDetails ? 'pointer' : 'default'
                }}
                onClick={() => hasDetails && setExpanded(!expanded)}
              >
                {Object.keys(log.old_value).length} field(s)
              </Box>
            </Tooltip>
          ) : (
            <em style={{ color: '#666' }}>—</em>
          )}
        </TableCell>
        <TableCell>
          {log.new_value ? (
            <Tooltip title="Click to expand details">
              <Box 
                sx={{ 
                  maxWidth: 200, 
                  overflow: 'hidden', 
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                  cursor: hasDetails ? 'pointer' : 'default'
                }}
                onClick={() => hasDetails && setExpanded(!expanded)}
              >
                {Object.keys(log.new_value).length} field(s)
              </Box>
            </Tooltip>
          ) : (
            <em style={{ color: '#666' }}>—</em>
          )}
        </TableCell>
        <TableCell>
          {hasDetails && (
            <IconButton 
              size="small" 
              onClick={() => setExpanded(!expanded)}
              aria-label="expand details"
            >
              {expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            </IconButton>
          )}
        </TableCell>
      </TableRow>
      {hasDetails && (
        <TableRow>
          <TableCell colSpan={9} sx={{ py: 0 }}>
            <Collapse in={expanded} timeout="auto" unmountOnExit>
              <Box sx={{ 
                p: 2, 
                bgcolor: 'background.paper',
                border: 1,
                borderColor: 'divider',
                borderRadius: 1
              }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="subtitle2">
                    Change Details
                  </Typography>
                  {isUpdateAction && (
                    <Tabs 
                      value={viewMode} 
                      onChange={(_, newValue) => setViewMode(newValue as 'diff' | 'separate')}
                    >
                      <Tab label="Diff View" value="diff" />
                      <Tab label="Separate View" value="separate" />
                    </Tabs>
                  )}
                </Box>
                
                {isUpdateAction && viewMode === 'diff' ? (
                  <JsonDiffViewer
                    oldValue={log.old_value}
                    newValue={log.new_value}
                  />
                ) : (
                  <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                    {log.old_value && (
                      <Box sx={{ flex: 1, minWidth: 300 }}>
                        <Typography variant="caption" color="text.secondary" gutterBottom>
                          Previous Value:
                        </Typography>
                        <Paper sx={{ 
                          p: 1, 
                          bgcolor: 'error.dark', 
                          color: 'error.contrastText',
                          border: 1,
                          borderColor: 'error.main'
                        }}>
                          <pre style={{ 
                            fontSize: '0.75rem', 
                            margin: 0, 
                            whiteSpace: 'pre-wrap',
                            fontFamily: 'monospace'
                          }}>
                            {formatValue(log.old_value)}
                          </pre>
                        </Paper>
                      </Box>
                    )}
                    {log.new_value && (
                      <Box sx={{ flex: 1, minWidth: 300 }}>
                        <Typography variant="caption" color="text.secondary" gutterBottom>
                          New Value:
                        </Typography>
                        <Paper sx={{ 
                          p: 1, 
                          bgcolor: 'success.dark', 
                          color: 'success.contrastText',
                          border: 1,
                          borderColor: 'success.main'
                        }}>
                          <pre style={{ 
                            fontSize: '0.75rem', 
                            margin: 0, 
                            whiteSpace: 'pre-wrap',
                            fontFamily: 'monospace'
                          }}>
                            {formatValue(log.new_value)}
                          </pre>
                        </Paper>
                      </Box>
                    )}
                  </Box>
                )}
              </Box>
            </Collapse>
          </TableCell>
        </TableRow>
      )}
    </>
  );
};

const AuditLogPage: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const rbac = useRBAC(projectId);
  const [searchParams, setSearchParams] = useSearchParams();

  // Permission guards
  if (!rbac.canViewProject()) {
    return (
      <AuthenticatedLayout showBackButton backTo="/dashboard">
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="error" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You don't have permission to view this project.
          </Typography>
        </Box>
      </AuthenticatedLayout>
    );
  }

  if (!rbac.canViewAudit()) {
    return (
      <AuthenticatedLayout showBackButton backTo="/dashboard">
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="error" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You don't have permission to view audit logs.
          </Typography>
        </Box>
      </AuthenticatedLayout>
    );
  }

  // Load environments for environment filter
  const { data: envResp } = useQuery({
    queryKey: ['project-environments', projectId],
    queryFn: async () => (await apiClient.listProjectEnvironments(projectId || '')).data,
    enabled: !!projectId,
  });
  const environments = envResp?.items ?? [];

  // Filters from query string
  const filter: AuditLogFilterValue = useMemo(() => ({
    environmentKey: searchParams.get('environment_key') || undefined,
    entity: searchParams.get('entity') || undefined,
    entityId: searchParams.get('entity_id') || undefined,
    actor: searchParams.get('actor') || undefined,
    from: searchParams.get('from') || undefined,
    to: searchParams.get('to') || undefined,
    sortBy: searchParams.get('sort_by') || 'created_at',
    sortOrder: (searchParams.get('sort_order') as 'asc' | 'desc') || 'desc',
  }), [searchParams]);

  const page = parseInt(searchParams.get('page') || '1', 10);
  const perPage = parseInt(searchParams.get('per_page') || String(DEFAULT_PER_PAGE), 10);

  const setFilter = (next: AuditLogFilterValue) => {
    const params: Record<string, string> = {};
    if (next.environmentKey) params.environment_key = next.environmentKey;
    if (next.entity) params.entity = next.entity;
    if (next.entityId) params.entity_id = next.entityId;
    if (next.actor) params.actor = next.actor;
    if (next.from) params.from = next.from;
    if (next.to) params.to = next.to;
    if (next.sortBy) params.sort_by = next.sortBy;
    if (next.sortOrder) params.sort_order = next.sortOrder;
    params.page = String(1); // reset page on filter change
    params.per_page = String(perPage);
    setSearchParams(params);
  };

  const { data, isLoading, error } = useQuery({
    queryKey: ['audit-log', projectId, filter, page, perPage],
    queryFn: async () => {
      const res = await apiClient.listProjectAuditLogs(
        projectId || '',
        filter.environmentKey,
        filter.entity,
        filter.entityId,
        filter.actor,
        filter.from ? new Date(filter.from).toISOString() : undefined,
        filter.to ? new Date(filter.to).toISOString() : undefined,
        filter.sortBy as any,
        filter.sortOrder as any,
        page,
        perPage,
      );
      return res.data;
    },
    enabled: !!projectId,
  });

  const logs = (data as any)?.items ?? [];
  const total = (data as any)?.pagination?.total ? Number((data as any).pagination.total) : 0;
  const totalPages = Math.max(1, Math.ceil(total / perPage));

  const handlePageChange = (_: any, p: number) => {
    const next = new URLSearchParams(searchParams);
    next.set('page', String(p));
    next.set('per_page', String(perPage));
    setSearchParams(next);
  };

  return (
    <AuthenticatedLayout showBackButton backTo={`/projects/${projectId}`}>
      <PageHeader title="Audit Log" icon={<ChangesIcon />} />

      <Paper sx={{ p: 2, mb: 2 }}>
        <AuditLogFilter value={filter} environments={environments as any} onChange={setFilter} />
      </Paper>

      <Paper sx={{ p: 0 }}>
        {isLoading ? (
          <Box sx={{ p: 4, textAlign: 'center' }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Box sx={{ p: 2 }}>
            <Alert severity="error">Failed to load audit logs</Alert>
          </Box>
        ) : (
          <>
            <TableContainer>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Time</TableCell>
                    <TableCell>Environment</TableCell>
                    <TableCell>Actor</TableCell>
                    <TableCell>Username</TableCell>
                    <TableCell>Entity</TableCell>
                    <TableCell>Action</TableCell>
                    <TableCell>Old Value</TableCell>
                    <TableCell>New Value</TableCell>
                    <TableCell></TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {logs.map((log: any) => (
                    <AuditLogRow key={log.id} log={log} />
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
            <Stack direction="row" justifyContent="center" sx={{ py: 2 }}>
              <MuiPagination count={totalPages} page={page} onChange={handlePageChange} />
            </Stack>
          </>
        )}
      </Paper>
    </AuthenticatedLayout>
  );
};

export default AuditLogPage;
