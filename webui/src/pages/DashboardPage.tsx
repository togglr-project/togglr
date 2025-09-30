import React, { useMemo, useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  CircularProgress,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Stack,
  Tooltip,
  Divider,
  Card,
  CardContent,
  LinearProgress,
  Skeleton,
  Alert,
  IconButton,
  Collapse,
} from '@mui/material';
import { 
  Dashboard as DashboardIcon,
  TrendingUp,
  TrendingDown,
  Warning,
  CheckCircle,
  Error,
  Schedule,
  Flag,
  People,
  Settings,
  ExpandMore,
  ExpandLess,
  Refresh,
  Visibility,
  VisibilityOff,
} from '@mui/icons-material';
import { Navigate, useLocation, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import { useAuth } from '../auth/AuthContext';
import apiClient from '../api/apiClient';
import type { DashboardOverviewResponse, Project } from '../generated/api/client';
import { useRBAC } from '../auth/permissions';

const useQueryParam = (key: string) => {
  const location = useLocation();
  const params = useMemo(() => new URLSearchParams(location.search), [location.search]);
  return params.get(key);
};

const HealthChip: React.FC<{ status?: string }> = ({ status }) => {
  const label = (status || 'unknown').toUpperCase();
  const color: 'success' | 'warning' | 'error' | 'default' =
    status === 'green' ? 'success' : status === 'yellow' ? 'warning' : status === 'red' ? 'error' : 'default';
  const icon = status === 'green' ? <CheckCircle /> : status === 'yellow' ? <Warning /> : status === 'red' ? <Error /> : undefined;
  
  return (
    <Chip 
      label={label} 
      color={color} 
      size="small" 
      icon={icon}
      sx={{ 
        fontWeight: 600,
        '& .MuiChip-icon': {
          fontSize: '0.9rem'
        }
      }}
    />
  );
};

// Metric Card Component
const MetricCard: React.FC<{
  title: string;
  value: string | number;
  subtitle?: string;
  icon: React.ReactNode;
  color: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info';
  trend?: 'up' | 'down' | 'neutral';
  trendValue?: string;
  loading?: boolean;
}> = ({ title, value, subtitle, icon, color, trend, trendValue, loading = false }) => {
  const getTrendIcon = () => {
    if (trend === 'up') return <TrendingUp sx={{ fontSize: '0.9rem' }} />;
    if (trend === 'down') return <TrendingDown sx={{ fontSize: '0.9rem' }} />;
    return null;
  };

  const getTrendColor = () => {
    if (trend === 'up') return 'success.main';
    if (trend === 'down') return 'error.main';
    return 'text.secondary';
  };

  return (
    <Card 
      sx={{ 
        height: '100%',
        background: 'linear-gradient(135deg, rgba(130, 82, 255, 0.02) 0%, rgba(130, 82, 255, 0.05) 100%)',
        border: '1px solid rgba(130, 82, 255, 0.1)',
        transition: 'all 0.3s ease',
        '&:hover': {
          transform: 'translateY(-2px)',
          boxShadow: '0 8px 25px rgba(130, 82, 255, 0.15)',
        }
      }}
    >
      <CardContent sx={{ p: 2.5 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1.5 }}>
          <Box
            sx={{
              p: 1.5,
              borderRadius: 2,
              backgroundColor: `${color}.main`,
              color: 'white',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            {icon}
          </Box>
          {trend && trendValue && (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
              {getTrendIcon()}
              <Typography 
                variant="caption" 
                sx={{ 
                  color: getTrendColor(),
                  fontWeight: 600,
                  fontSize: '0.75rem'
                }}
              >
                {trendValue}
              </Typography>
            </Box>
          )}
        </Box>
        
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5, color: 'text.primary' }}>
          {loading ? <Skeleton width="60%" /> : value}
        </Typography>
        
        <Typography variant="body2" sx={{ color: 'text.secondary', mb: 0.5 }}>
          {loading ? <Skeleton width="80%" /> : title}
        </Typography>
        
        {subtitle && (
          <Typography variant="caption" sx={{ color: 'text.secondary', fontSize: '0.75rem' }}>
            {loading ? <Skeleton width="60%" /> : subtitle}
          </Typography>
        )}
      </CardContent>
    </Card>
  );
};

// Enhanced Table Component
const EnhancedTable: React.FC<{
  title: string;
  subtitle?: string;
  data: any[];
  columns: Array<{
    key: string;
    label: string;
    align?: 'left' | 'center' | 'right';
    render?: (value: any, row: any) => React.ReactNode;
    width?: string;
  }>;
  loading?: boolean;
  emptyMessage?: string;
  collapsible?: boolean;
  defaultExpanded?: boolean;
}> = ({ title, subtitle, data, columns, loading = false, emptyMessage = "No data", collapsible = false, defaultExpanded = true }) => {
  const [expanded, setExpanded] = useState(defaultExpanded);

  const TableContent = () => (
    <TableContainer>
      <Table size="small">
        <TableHead>
          <TableRow sx={{ backgroundColor: 'rgba(130, 82, 255, 0.02)' }}>
            {columns.map((column) => (
              <TableCell 
                key={column.key}
                align={column.align || 'left'}
                sx={{ 
                  fontWeight: 600, 
                  color: 'text.primary',
                  borderBottom: '2px solid rgba(130, 82, 255, 0.1)',
                  width: column.width
                }}
              >
                {column.label}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {loading ? (
            Array.from({ length: 3 }).map((_, idx) => (
              <TableRow key={idx}>
                {columns.map((column) => (
                  <TableCell key={column.key} align={column.align || 'left'}>
                    <Skeleton width="80%" />
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : data.length > 0 ? (
            data.map((row, idx) => (
              <TableRow 
                key={idx} 
                hover
                sx={{ 
                  '&:hover': {
                    backgroundColor: 'rgba(130, 82, 255, 0.04)',
                  },
                  '&:nth-of-type(even)': {
                    backgroundColor: 'rgba(130, 82, 255, 0.01)',
                  }
                }}
              >
                {columns.map((column) => (
                  <TableCell 
                    key={column.key}
                    align={column.align || 'left'}
                    sx={{ 
                      borderBottom: '1px solid rgba(130, 82, 255, 0.05)',
                      py: 1.5
                    }}
                  >
                    {column.render ? column.render(row[column.key], row) : row[column.key]}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} align="center" sx={{ py: 4 }}>
                <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    {emptyMessage}
                  </Typography>
                </Box>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );

  return (
    <Paper sx={{ overflow: 'hidden' }}>
      <Box sx={{ p: 2, borderBottom: '1px solid rgba(130, 82, 255, 0.1)' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 600 }}>
              {title}
            </Typography>
            {subtitle && (
              <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                {subtitle}
              </Typography>
            )}
          </Box>
          {collapsible && (
            <IconButton 
              onClick={() => setExpanded(!expanded)}
              size="small"
              sx={{ 
                color: 'primary.main',
                '&:hover': {
                  backgroundColor: 'rgba(130, 82, 255, 0.08)'
                }
              }}
            >
              {expanded ? <ExpandLess /> : <ExpandMore />}
            </IconButton>
          )}
        </Box>
      </Box>
      
      {collapsible ? (
        <Collapse in={expanded}>
          <TableContent />
        </Collapse>
      ) : (
        <TableContent />
      )}
    </Paper>
  );
};

const SectionTitle: React.FC<{ title: string; subtitle?: string }> = ({ title, subtitle }) => (
  <Box sx={{ display: 'flex', alignItems: 'baseline', justifyContent: 'space-between', mb: 1 }}>
    <Typography variant="h6" sx={{ color: 'primary.light' }}>{title}</Typography>
    {subtitle && (
      <Typography variant="caption" color="text.secondary">{subtitle}</Typography>
    )}
  </Box>
);

const DashboardPage: React.FC = () => {
  const { isAuthenticated, user } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  // URL-synced filters
  const qEnv = useQueryParam('environment_key') || 'prod';
  const qProjectId = useQueryParam('project_id') || '';
  const [environmentKey, setEnvironmentKey] = useState<string>(qEnv);
  const [projectId, setProjectId] = useState<string>(qProjectId);

  const setQuery = (nextEnv: string, nextProjectId: string) => {
    const params = new URLSearchParams(location.search);
    params.set('environment_key', nextEnv);
    if (nextProjectId) {
      params.set('project_id', nextProjectId);
    } else {
      params.delete('project_id');
    }
    navigate({ pathname: location.pathname, search: params.toString() }, { replace: true });
  };

  // Projects for selector
  const { data: projects } = useQuery<Project[]>({
    queryKey: ['projects'],
    queryFn: async () => (await apiClient.listProjects()).data,
  });

  // Filter projects by access permissions
  const accessibleProjects = useMemo(() => {
    if (!projects || !user) return [];
    
    // If user is superuser, show all projects
    if (user.is_superuser) return projects;
    
    // Otherwise filter by project_permissions
    return projects.filter(project => {
      const permissions = user.project_permissions?.[project.id];
      return permissions && permissions.includes('project.view');
    });
  }, [projects, user]);

  // Dashboard data
  const { data, isLoading, error, refetch, isFetching } = useQuery<DashboardOverviewResponse>({
    queryKey: ['dashboard_overview', environmentKey, projectId],
    queryFn: async () => (await apiClient.getDashboardOverview(environmentKey, projectId || undefined, 20)).data,
  });

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  const handleEnvChange = (val: string) => {
    setEnvironmentKey(val);
    setQuery(val, projectId);
    // react-query will refetch due to key change
  };
  const handleProjectChange = (val: string) => {
    setProjectId(val);
    setQuery(environmentKey, val);
  };

  // Calculate metrics
  const totalFeatures = data?.projects?.reduce((sum, p) => sum + (p.total_features || 0), 0) || 0;
  const totalEnabled = data?.projects?.reduce((sum, p) => sum + (p.enabled_features || 0), 0) || 0;
  const totalPending = data?.projects?.reduce((sum, p) => sum + (p.pending_features || 0), 0) || 0;
  const totalGuarded = data?.projects?.reduce((sum, p) => sum + (p.guarded_features || 0), 0) || 0;
  const riskyFeaturesCount = data?.risky_features?.length || 0;
  const recentActivityCount = data?.recent_activity?.length || 0;

  // Calculate health percentage
  const healthPercentage = totalFeatures > 0 ? Math.round((totalEnabled / totalFeatures) * 100) : 0;

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="Dashboard"
        subtitle="Overview of projects health and recent activity"
        icon={<DashboardIcon />}
      >
        <IconButton 
          onClick={() => refetch()} 
          disabled={isFetching}
          sx={{ 
            color: 'primary.main',
            '&:hover': {
              backgroundColor: 'rgba(130, 82, 255, 0.08)'
            }
          }}
        >
          <Refresh sx={{ transform: isFetching ? 'rotate(360deg)' : 'rotate(0deg)', transition: 'transform 0.5s ease' }} />
        </IconButton>
      </PageHeader>

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3, background: 'linear-gradient(135deg, rgba(130, 82, 255, 0.02) 0%, rgba(130, 82, 255, 0.05) 100%)' }}>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ xs: 'stretch', sm: 'center' }}>
          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel id="env-select-label">Environment</InputLabel>
            <Select
              labelId="env-select-label"
              label="Environment"
              value={environmentKey}
              onChange={(e) => handleEnvChange(e.target.value)}
            >
              <MenuItem value="prod">Production</MenuItem>
              <MenuItem value="stage">Staging</MenuItem>
              <MenuItem value="dev">Development</MenuItem>
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 240 }}>
            <InputLabel id="project-select-label">Project (optional)</InputLabel>
            <Select
              labelId="project-select-label"
              label="Project (optional)"
              value={projectId}
              onChange={(e) => handleProjectChange(e.target.value)}
            >
              <MenuItem value="">
                <em>All projects</em>
              </MenuItem>
              {accessibleProjects?.map((p) => (
                <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>
              ))}
            </Select>
          </FormControl>

          <Box sx={{ flex: 1 }} />
          {(isFetching || isLoading) && (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <CircularProgress size={20} />
              <Typography variant="caption" color="text.secondary">Loading...</Typography>
            </Box>
          )}
        </Stack>
      </Paper>

      {isLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Alert severity="error" sx={{ mb: 2 }}>
          <Typography variant="h6" sx={{ mb: 1 }}>Failed to load dashboard</Typography>
          <Typography variant="body2">{(error as any)?.message || 'Unknown error'}</Typography>
        </Alert>
      ) : data ? (
        <Box>
          {/* Metrics Cards */}
          <Grid container spacing={3} sx={{ mb: 4 }}>
            <Grid item xs={12} sm={6} md={3}>
              <MetricCard
                title="Total Features"
                value={totalFeatures}
                subtitle={`Across ${data.projects?.length || 0} projects`}
                icon={<Flag />}
                color="primary"
                loading={isLoading}
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <MetricCard
                title="Enabled Features"
                value={totalEnabled}
                subtitle={`${healthPercentage}% of total`}
                icon={<CheckCircle />}
                color="success"
                trend={healthPercentage > 80 ? 'up' : healthPercentage < 60 ? 'down' : 'neutral'}
                trendValue={healthPercentage > 80 ? 'Good' : healthPercentage < 60 ? 'Low' : 'Normal'}
                loading={isLoading}
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <MetricCard
                title="Pending Changes"
                value={totalPending}
                subtitle="Awaiting approval"
                icon={<Schedule />}
                color="warning"
                trend={totalPending > 10 ? 'up' : 'neutral'}
                trendValue={totalPending > 10 ? 'High' : 'Normal'}
                loading={isLoading}
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <MetricCard
                title="Risky Features"
                value={riskyFeaturesCount}
                subtitle="Critical or guarded"
                icon={<Warning />}
                color="error"
                trend={riskyFeaturesCount > 5 ? 'up' : 'neutral'}
                trendValue={riskyFeaturesCount > 5 ? 'High Risk' : 'Safe'}
                loading={isLoading}
              />
            </Grid>
          </Grid>

          {/* Health Progress Bar */}
          <Paper sx={{ p: 2, mb: 3, background: 'linear-gradient(135deg, rgba(130, 82, 255, 0.02) 0%, rgba(130, 82, 255, 0.05) 100%)' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
              <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 600 }}>
                System Health
              </Typography>
              <Typography variant="h6" sx={{ color: 'text.primary', fontWeight: 700 }}>
                {healthPercentage}%
              </Typography>
            </Box>
            <LinearProgress 
              variant="determinate" 
              value={healthPercentage} 
              sx={{ 
                height: 8, 
                borderRadius: 4,
                backgroundColor: 'rgba(130, 82, 255, 0.1)',
                '& .MuiLinearProgress-bar': {
                  backgroundColor: healthPercentage > 80 ? '#4CAF50' : healthPercentage > 60 ? '#FF9800' : '#F44336',
                  borderRadius: 4,
                }
              }} 
            />
            <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
              {totalEnabled} of {totalFeatures} features are enabled
            </Typography>
          </Paper>

          {/* Project Overview Table */}
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <EnhancedTable
                title="Project Overview"
                subtitle="Health and feature totals per project"
                data={data.projects || []}
                columns={[
                  { key: 'project_name', label: 'Project', width: '25%', render: (value, row) => value || row.project_id || '—' },
                  { key: 'total_features', label: 'Total', align: 'right' },
                  { key: 'enabled_features', label: 'Enabled', align: 'right' },
                  { key: 'disabled_features', label: 'Disabled', align: 'right' },
                  { key: 'guarded_features', label: 'Guarded', align: 'right' },
                  { key: 'pending_features', label: 'Pending', align: 'right' },
                  { 
                    key: 'health_status', 
                    label: 'Health', 
                    render: (value) => <HealthChip status={value} />
                  },
                ]}
                loading={isLoading}
                emptyMessage="No projects found"
              />
          </Grid>

          {/* Category Health */}
          <Grid item xs={12}>
              <EnhancedTable
                title="Category Health"
                subtitle="Per-category health inside project"
                data={data.categories || []}
                columns={[
                  { key: 'category_name', label: 'Category', width: '25%', render: (value, row) => value || row.category_slug || '—' },
                  { key: 'total_features', label: 'Total', align: 'right' },
                  { key: 'enabled_features', label: 'Enabled', align: 'right' },
                  { key: 'disabled_features', label: 'Disabled', align: 'right' },
                  { key: 'pending_features', label: 'Pending', align: 'right' },
                  { key: 'guarded_features', label: 'Guarded', align: 'right' },
                  { 
                    key: 'health_status', 
                    label: 'Health', 
                    render: (value) => <HealthChip status={value} />
                  },
                ]}
                loading={isLoading}
                emptyMessage="No categories found"
                collapsible={true}
                defaultExpanded={false}
              />
          </Grid>

          {/* Feature Activity */}
          <Grid item xs={12} md={6}>
              <Paper sx={{ height: '100%', background: 'linear-gradient(135deg, rgba(130, 82, 255, 0.02) 0%, rgba(130, 82, 255, 0.05) 100%)' }}>
                <Box sx={{ p: 2, borderBottom: '1px solid rgba(130, 82, 255, 0.1)' }}>
                  <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 600 }}>
                    Feature Activity — Upcoming
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                    Scheduled changes
                  </Typography>
                </Box>
                <Box sx={{ p: 2 }}>
                  <Stack spacing={2}>
                    {isLoading ? (
                      Array.from({ length: 3 }).map((_, i) => (
                        <Box key={i} sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                          <Skeleton width="60%" height={20} />
                          <Skeleton width="30%" height={20} />
                        </Box>
                      ))
                    ) : data.feature_activity?.upcoming?.length ? (
                  data.feature_activity.upcoming.map((u, i) => (
                        <Box key={i} sx={{ 
                          display: 'flex', 
                          justifyContent: 'space-between', 
                          alignItems: 'center',
                          p: 1.5,
                          borderRadius: 1,
                          backgroundColor: 'rgba(130, 82, 255, 0.02)',
                          border: '1px solid rgba(130, 82, 255, 0.05)',
                          '&:hover': {
                            backgroundColor: 'rgba(130, 82, 255, 0.04)',
                          }
                        }}>
                          <Typography variant="body2" sx={{ fontWeight: 500 }}>
                            {u.feature_name ?? '—'}
                          </Typography>
                      <Stack direction="row" spacing={1} alignItems="center">
                            <Chip 
                              size="small" 
                              label={(u.next_state as any)?.toString()} 
                              color={u.next_state === 'enabled' ? 'success' : 'default'}
                              sx={{ fontWeight: 600 }}
                            />
                            <Typography variant="caption" color="text.secondary">
                              {u.at as any}
                            </Typography>
                      </Stack>
                    </Box>
                  ))
                ) : (
                      <Box sx={{ textAlign: 'center', py: 3 }}>
                        <Schedule sx={{ fontSize: 48, color: 'text.secondary', mb: 1 }} />
                        <Typography variant="body2" color="text.secondary">
                          No upcoming events
                        </Typography>
                      </Box>
                )}
              </Stack>
                </Box>
            </Paper>
          </Grid>
          <Grid item xs={12} md={6}>
              <Paper sx={{ height: '100%', background: 'linear-gradient(135deg, rgba(130, 82, 255, 0.02) 0%, rgba(130, 82, 255, 0.05) 100%)' }}>
                <Box sx={{ p: 2, borderBottom: '1px solid rgba(130, 82, 255, 0.1)' }}>
                  <Typography variant="h6" sx={{ color: 'primary.main', fontWeight: 600 }}>
                    Feature Activity — Recent
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                    Recent changes
                  </Typography>
                </Box>
                <Box sx={{ p: 2 }}>
                  <Stack spacing={2}>
                    {isLoading ? (
                      Array.from({ length: 3 }).map((_, i) => (
                        <Box key={i} sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                          <Skeleton width="60%" height={20} />
                          <Skeleton width="30%" height={20} />
                        </Box>
                      ))
                    ) : data.feature_activity?.recent?.length ? (
                  data.feature_activity.recent.map((r, i) => (
                        <Box key={i} sx={{ 
                          display: 'flex', 
                          justifyContent: 'space-between', 
                          alignItems: 'center',
                          p: 1.5,
                          borderRadius: 1,
                          backgroundColor: 'rgba(130, 82, 255, 0.02)',
                          border: '1px solid rgba(130, 82, 255, 0.05)',
                          '&:hover': {
                            backgroundColor: 'rgba(130, 82, 255, 0.04)',
                          }
                        }}>
                          <Typography variant="body2" sx={{ fontWeight: 500 }}>
                            {r.feature_name ?? '—'}
                          </Typography>
                      <Stack direction="row" spacing={1} alignItems="center">
                            <Chip 
                              size="small" 
                              label={r.action ?? ''} 
                              color="primary"
                              variant="outlined"
                              sx={{ fontWeight: 600 }}
                            />
                        <Tooltip title={r.actor ?? ''}>
                              <Typography variant="caption" color="text.secondary">
                                {r.at as any}
                              </Typography>
                        </Tooltip>
                      </Stack>
                    </Box>
                  ))
                ) : (
                      <Box sx={{ textAlign: 'center', py: 3 }}>
                        <Flag sx={{ fontSize: 48, color: 'text.secondary', mb: 1 }} />
                        <Typography variant="body2" color="text.secondary">
                          No recent feature activity
                        </Typography>
                      </Box>
                )}
              </Stack>
                </Box>
            </Paper>
          </Grid>

          {/* Top Risky Features */}
          <Grid item xs={12}>
              <EnhancedTable
                title="Top Risky Features"
                subtitle="Features with critical, guarded, or auto-disable tags"
                data={data.risky_features || []}
                columns={[
                  { key: 'feature_name', label: 'Feature', width: '25%' },
                  { key: 'project_name', label: 'Project', width: '20%', render: (value, row) => value || row.project_id || '—' },
                  { 
                    key: 'risky_tags', 
                    label: 'Tags', 
                    render: (value) => (
                      <Chip 
                        size="small" 
                        label={value || '—'} 
                        color="error" 
                        variant="outlined"
                        sx={{ fontWeight: 600 }}
                      />
                    )
                  },
                  { 
                    key: 'enabled', 
                    label: 'Enabled', 
                    align: 'right',
                    render: (value) => (
                      <Chip 
                        size="small" 
                        label={String(value)} 
                        color={value ? 'success' : 'default'}
                        sx={{ fontWeight: 600 }}
                      />
                    )
                  },
                  { 
                    key: 'has_pending', 
                    label: 'Pending', 
                    align: 'right',
                    render: (value) => (
                      <Chip 
                        size="small" 
                        label={String(value)} 
                        color={value ? 'warning' : 'default'}
                        sx={{ fontWeight: 600 }}
                      />
                    )
                  },
                ]}
                loading={isLoading}
                emptyMessage="No risky features found"
                collapsible={true}
                defaultExpanded={false}
              />
          </Grid>

          {/* Recent Activity */}
          <Grid item xs={12}>
              <EnhancedTable
                title="Recent Activity"
                subtitle="Batched changes by request"
                data={data.recent_activity || []}
                columns={[
                  { key: 'project_name', label: 'Project', width: '20%', render: (value, row) => value || row.project_id || '—' },
                  { key: 'actor', label: 'Actor', width: '15%' },
                  { 
                    key: 'status', 
                    label: 'Status', 
                    width: '10%',
                    render: (value) => (
                      <Chip 
                        size="small" 
                        label={value || '—'} 
                        color={value === 'applied' ? 'success' : value === 'pending' ? 'warning' : 'default'}
                        sx={{ fontWeight: 600 }}
                      />
                    )
                  },
                  { key: 'created_at', label: 'When', width: '20%' },
                  { 
                    key: 'changes', 
                    label: 'Changes', 
                    render: (value) => (
                      <Stack direction="row" spacing={0.5} sx={{ flexWrap: 'wrap', gap: 0.5 }}>
                        {value?.map((c: any, i: number) => (
                          <Chip 
                            key={i} 
                            size="small" 
                            label={`${c.entity}:${c.action}`} 
                            color="primary"
                            variant="outlined"
                            sx={{ fontSize: '0.7rem' }}
                          />
                            ))}
                          </Stack>
                    )
                  },
                ]}
                loading={isLoading}
                emptyMessage="No recent activity"
                collapsible={true}
                defaultExpanded={true}
              />
          </Grid>

          {/* Pending Summary */}
          <Grid item xs={12}>
              <EnhancedTable
                title="Pending Summary"
                subtitle="Summary of pending changes across projects"
                data={data.pending_summary || []}
                columns={[
                  { key: 'project_name', label: 'Project', width: '25%', render: (value, row) => value || row.project_id || '—' },
                  { key: 'total_pending', label: 'Total', align: 'right' },
                  { key: 'pending_feature_changes', label: 'Feature Changes', align: 'right' },
                  { key: 'pending_guarded_changes', label: 'Guarded Changes', align: 'right' },
                  { key: 'oldest_request_at', label: 'Oldest Request', width: '20%' },
                ]}
                loading={isLoading}
                emptyMessage="No pending changes"
                collapsible={true}
                defaultExpanded={false}
              />
            </Grid>
          </Grid>
        </Box>
      ) : null}
    </AuthenticatedLayout>
  );
};

export default DashboardPage;
