import React, { useState } from 'react';
import { Box, Chip, CircularProgress, Dialog, DialogActions, DialogContent, DialogTitle, Button, Typography, Switch, Tooltip, FormControlLabel, Collapse, IconButton } from '@mui/material';
import { WarningAmber, ExpandMore, ExpandLess, Schedule as ScheduleIcon } from '@mui/icons-material';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { FeatureExtended, FeatureDetailsResponse, Segment } from '../../generated/api/client';
import { useAuth } from '../../auth/AuthContext';
import EditFeatureDialog from './EditFeatureDialog';
import { getNextStateDescription } from '../../utils/timeUtils';
import { useFeatureHasPendingChanges } from '../../hooks/useProjectPendingChanges';
import { Pending as PendingIcon } from '@mui/icons-material';
import GuardResponseHandler from '../pending-changes/GuardResponseHandler';
import { useApprovePendingChange } from '../../hooks/usePendingChanges';
import type { AuthCredentialsMethodEnum } from '../../generated/api/client';
import { useRBAC } from '../../auth/permissions';
import { getHealthStatusColor, getHealthStatusVariant } from '../../utils/healthStatus';
import { getFirstEnabledAlgorithmSlug } from '../../utils/algorithmUtils';

export interface FeatureDetailsDialogProps {
  open: boolean;
  onClose: () => void;
  feature: FeatureExtended | null;
  environmentKey: string;
}

const FeatureDetailsDialog: React.FC<FeatureDetailsDialogProps> = ({ open, onClose, feature, environmentKey }) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const { data: featureDetails, isLoading, error } = useQuery<FeatureDetailsResponse>({
    queryKey: ['feature-details', feature?.id],
    queryFn: async () => {
      const res = await apiClient.getFeature(feature!.id, environmentKey);
      return res.data as FeatureDetailsResponse;
    },
    enabled: open && !!feature?.id,
    staleTime: 0, // No caching - always fetch fresh data
    refetchOnWindowFocus: true, // Refetch when window gains focus
  });

  const projectId = featureDetails?.feature.project_id;
  const rbac = useRBAC(projectId || '');
  const canToggleFeature = rbac.canToggleFeature();
  const { data: segments } = useQuery<Segment[]>({
    queryKey: ['project-segments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectSegments(projectId!);
      const resp = res.data as any;
      return Array.isArray(resp?.items) ? (resp.items as Segment[]) : (resp as Segment[]);
    },
    enabled: Boolean(projectId),
    staleTime: 0, // No caching - always fetch fresh data
    refetchOnWindowFocus: true, // Refetch when window gains focus
  });

  const getVariantName = (id: string) => {
    const arr = featureDetails?.variants || [];
    const found = arr.find(v => v.id === id);
    return found ? (found.name || found.id) : id;
  };

  const getSegmentName = (segmentId: string | undefined) => {
    if (!segmentId) return 'user defined rule';
    if (!segments || !Array.isArray(segments)) return 'user defined rule';
    const found = segments.find(s => s.id === segmentId);
    return found ? found.name : 'user defined rule';
  };

  const renderConditionExpression = (expression: any, depth: number = 0): React.ReactNode => {
    if (!expression) return null;

    // Handle legacy array format
    if (Array.isArray(expression)) {
      return (
        <Box sx={{ ml: depth * 2 }}>
          {expression.map((c: any, idx: number) => (
            <Box key={idx} sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr 1fr', md: '1.2fr 0.8fr 1.5fr' }, gap: 1, mb: 0.5, alignItems: 'center' }}>
              <Typography variant="body2">{c.attribute}</Typography>
              <Typography variant="body2" color="text.secondary">{c.operator}</Typography>
              <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                {typeof c.value === 'string' ? c.value : JSON.stringify(c.value)}
              </Typography>
            </Box>
          ))}
        </Box>
      );
    }

    // Handle single condition
    if (expression.condition) {
      const c = expression.condition;
      return (
        <Box sx={{ ml: depth * 2 }}>
          <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr 1fr', md: '1.2fr 0.8fr 1.5fr' }, gap: 1, mb: 0.5, alignItems: 'center' }}>
            <Typography variant="body2">{c.attribute}</Typography>
            <Typography variant="body2" color="text.secondary">{c.operator}</Typography>
            <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
              {typeof c.value === 'string' ? c.value : JSON.stringify(c.value)}
            </Typography>
          </Box>
        </Box>
      );
    }

    // Handle group with logical operator
    if (expression.group) {
      const group = expression.group;
      const operatorText = group.operator === 'and' ? 'AND' : group.operator === 'or' ? 'OR' : 'AND NOT';
      const operatorColor = group.operator === 'and' ? 'success' : group.operator === 'or' ? 'info' : 'error';
      
      return (
        <Box sx={{ ml: depth * 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
            <Typography variant="body2" color="text.secondary">Group:</Typography>
            <Chip size="small" label={operatorText} color={operatorColor} />
          </Box>
          <Box sx={{ borderLeft: '2px solid', borderColor: 'divider', pl: 2 }}>
            {group.children?.map((child: any, idx: number) => (
              <Box key={idx} sx={{ mb: 1 }}>
                {renderConditionExpression(child, depth + 1)}
              </Box>
            ))}
          </Box>
        </Box>
      );
    }

    return null;
  };

  const canToggle = featureDetails ? canToggleFeature : false;
  
  // Check if feature has pending changes
  const hasPendingChanges = useFeatureHasPendingChanges(feature?.id || '', featureDetails?.feature.project_id);

  const toggleMutation = useMutation({
    mutationFn: async (enabled: boolean) => {
      if (!featureDetails) return;
      const response = await apiClient.toggleFeature(featureDetails.feature.id, environmentKey, { enabled });
      return response;
    },
    onSuccess: (response) => {
      if (!featureDetails) return;
      
      // Check if we got a 202 response (pending change created)
      if (response.status === 202 && response.data) {
        setGuardResponse({
          pendingChange: response.data,
        });
        return;
      }
      
      // Check if we got a 409 response (conflict)
      if (response.status === 409) {
        setGuardResponse({
          conflictError: 'Feature is already locked by another pending change',
        });
        return;
      }
      
      // Check if we got a 403 response (forbidden)
      if (response.status === 403) {
        setGuardResponse({
          forbiddenError: 'You don\'t have permission to modify this guarded feature',
        });
        return;
      }
      
      // Normal success - invalidate queries
      queryClient.invalidateQueries({ queryKey: ['feature-details'] });
      queryClient.invalidateQueries({ queryKey: ['project-features'] });
      queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
      queryClient.invalidateQueries({ queryKey: ['feature-details', featureDetails.feature.id] });
      queryClient.invalidateQueries({ queryKey: ['project-features', featureDetails.feature.project_id] });
    },
    onError: (error: any) => {
      // Check for conflict error
      if (error?.response?.status === 409) {
        setGuardResponse({
          conflictError: error.response.data?.error?.message || 'Conflict: Another pending change exists for this feature',
        });
      }
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async () => {
      if (!featureDetails) return;
      const response = await apiClient.deleteFeature(featureDetails.feature.id, environmentKey);
      return response;
    },
    onSuccess: (response) => {
      if (!featureDetails) return;
      
      // Check if we got a 202 response (pending change created)
      if (response.status === 202 && response.data) {
        setGuardResponse({
          pendingChange: response.data,
        });
        return;
      }
      
      // Check if we got a 409 response (conflict)
      if (response.status === 409) {
        setGuardResponse({
          conflictError: 'Feature is already locked by another pending change',
        });
        return;
      }
      
      // Check if we got a 403 response (forbidden)
      if (response.status === 403) {
        setGuardResponse({
          forbiddenError: 'You don\'t have permission to modify this guarded feature',
        });
        return;
      }
      
      // Normal success - invalidate queries and close dialog
      queryClient.invalidateQueries({ queryKey: ['feature-details'] });
      queryClient.invalidateQueries({ queryKey: ['project-features'] });
      queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
      queryClient.invalidateQueries({ queryKey: ['feature-details', featureDetails.feature.id] });
      queryClient.invalidateQueries({ queryKey: ['project-features', featureDetails.feature.project_id] });
      onClose();
    },
    onError: (error: any) => {
      // Check for conflict error
      if (error?.response?.status === 409) {
        setGuardResponse({
          conflictError: error.response.data?.error?.message || 'Conflict: Another pending change exists for this feature',
        });
      }
    },
  });

  const [editOpen, setEditOpen] = useState(false);
  const [expandedSections, setExpandedSections] = useState({
    assign: false,
    include: false,
    exclude: false
  });
  const [guardResponse, setGuardResponse] = useState<{
    pendingChange?: any;
    conflictError?: string;
    forbiddenError?: string;
  }>({});
  const canManage = featureDetails ? Boolean(user?.is_superuser || user?.project_permissions?.[featureDetails.feature.project_id]?.includes('feature.manage')) : false;
  
  const approveMutation = useApprovePendingChange();

  const toggleSection = (section: keyof typeof expandedSections) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  const handleAutoApprove = async (authMethod: AuthCredentialsMethodEnum, credential: string, sessionId?: string) => {
    if (!guardResponse.pendingChange?.id) return;

    try {
      await approveMutation.mutateAsync({
        id: guardResponse.pendingChange.id,
        request: {
          approver_user_id: user?.id || 0,
          approver_name: user?.username || 'Unknown',
          auth: {
            method: authMethod,
            credential: credential,
            ...(sessionId && { session_id: sessionId }),
          },
        },
      });

      // Success - invalidate queries and close guard response
      queryClient.invalidateQueries({ queryKey: ['feature-details'] });
      queryClient.invalidateQueries({ queryKey: ['project-features'] });
      queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
      queryClient.invalidateQueries({ queryKey: ['feature-details', featureDetails?.feature.id] });
      queryClient.invalidateQueries({ queryKey: ['project-features', featureDetails?.feature.project_id] });
      setGuardResponse({});
    } catch (error) {
      console.error('Auto-approve failed:', error);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle sx={{ color: 'primary.main' }}>Feature Details</DialogTitle>
      <DialogContent>
        {!feature || isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Typography color="error">Failed to load feature details.</Typography>
        ) : featureDetails ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="h6">{featureDetails.feature.name}</Typography>
                {canToggle && !hasPendingChanges ? (
                  <Tooltip title={featureDetails.feature.enabled ? 'Disable feature' : 'Enable feature'}>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={featureDetails.feature.enabled}
                          onChange={(e) => toggleMutation.mutate(e.target.checked)}
                          disabled={toggleMutation.isPending}
                          inputProps={{ 'aria-label': 'toggle feature in dialog' }}
                        />
                      }
                      label={"Enable"}
                    />
                  </Tooltip>
                ) : (
                  <Tooltip title={
                    hasPendingChanges 
                      ? "Cannot toggle: feature has pending changes awaiting approval"
                      : "You don't have permission to toggle features in this project"
                  }>
                    <span>
                      <FormControlLabel
                        control={<Switch checked={featureDetails.feature.enabled} disabled />}
                        label={"Enable"}
                      />
                    </span>
                  </Tooltip>
                )}
              </Box>
              <Typography variant="body2" color="text.secondary">Key: {featureDetails.feature.key}</Typography>
              {featureDetails.feature.description && (
                <Typography variant="body2" sx={{ mt: 1 }}>{featureDetails.feature.description}</Typography>
              )}
              <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                <Chip size="small" label={`kind: ${featureDetails.feature.kind}`} />
                {featureDetails.feature.enabled && featureDetails.feature.health_status && (
                  <Chip 
                    size="small" 
                    label={`health: ${featureDetails.feature.health_status}`} 
                    color={getHealthStatusColor(featureDetails.feature.health_status)} 
                    variant={getHealthStatusVariant(featureDetails.feature.health_status)} 
                  />
                )}
                <Chip size="small" label={`default: ${featureDetails.feature.default_value}`} />
                {featureDetails.feature.kind === 'multivariant' && (
                  <Chip size="small" label={`rollout key: ${featureDetails.feature.rollout_key || '-'}`} />
                )}
                {getFirstEnabledAlgorithmSlug(featureDetails.feature.algorithms) && (
                  <Chip 
                    size="small" 
                    label={`Algorithm: ${getFirstEnabledAlgorithmSlug(featureDetails.feature.algorithms)}`}
                    color="info"
                    variant="outlined"
                  />
                )}
                <Chip size="small" label={featureDetails.feature.is_active ? 'active' : 'not active'} color={featureDetails.feature.is_active ? 'success' : 'default'} />
                {featureDetails.feature.next_state !== undefined && featureDetails.feature.next_state_time && (
                  <Chip 
                    size="small" 
                    icon={<ScheduleIcon />}
                    label={getNextStateDescription(featureDetails.feature.next_state, featureDetails.feature.next_state_time) || 'Scheduled'} 
                    color={featureDetails.feature.next_state ? 'info' : 'warning'}
                    variant="outlined"
                  />
                )}
                {hasPendingChanges && (
                  <Tooltip title="This feature has pending changes awaiting approval">
                    <Chip 
                      size="small" 
                      icon={<PendingIcon />}
                      label="Pending" 
                      color="warning"
                      variant="filled"
                    />
                  </Tooltip>
                )}
              </Box>
              
              {/* Tags */}
              {featureDetails.tags && featureDetails.tags.length > 0 && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
                    Tags
                  </Typography>
                  <Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
                    {featureDetails.tags.map((tag) => (
                      <Chip
                        key={tag.id}
                        label={tag.slug}
                        size="small"
                        sx={{ 
                          fontSize: '0.7rem', 
                          height: 20,
                          backgroundColor: tag.color || 'default',
                          color: tag.color ? 'white' : 'inherit',
                          '& .MuiChip-label': {
                            color: tag.color ? 'white' : 'inherit'
                          }
                        }}
                      />
                    ))}
                  </Box>
                </Box>
              )}
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Variants</Typography>
              {featureDetails.variants && featureDetails.variants.length > 0 ? (
                <Box>
                  {featureDetails.variants.map((v) => (
                    <Box key={v.id} sx={{ display: 'grid', gridTemplateColumns: { xs: '2fr 1fr' }, gap: 1, mb: 0.5 }}>
                      <Typography variant="body2">{v.name}</Typography>
                      <Typography variant="body2" color="text.secondary">{v.rollout_percent}%</Typography>
                    </Box>
                  ))}
                </Box>
              ) : (
                <Typography variant="body2" color="text.secondary">No variants</Typography>
              )}
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Rules</Typography>
              {featureDetails.rules && featureDetails.rules.some((r: any) => (r as any).segment_id && (r as any).is_customized) && (
                <Box sx={{ mb: 1 }}>
                  <Typography variant="body2" color="warning.main" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <WarningAmber fontSize="small" />
                    Some rules are customized. These rules differs from segment templates.
                  </Typography>
                </Box>
              )}
              {!featureDetails.rules || featureDetails.rules.length === 0 ? (
                <Typography variant="body2" color="text.secondary">No rules</Typography>
              ) : (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <Box sx={{ border: '2px solid', borderColor: 'success.light', borderRadius: 1, p: 1 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                      <Typography variant="subtitle2">Assign request to variant rules</Typography>
                      <IconButton size="small" onClick={() => toggleSection('assign')}>
                        {expandedSections.assign ? <ExpandLess /> : <ExpandMore />}
                      </IconButton>
                    </Box>
                    <Collapse in={expandedSections.assign}>
                      {featureDetails.rules.filter(r => r.action === 'assign').sort((a,b) => a.priority - b.priority).length === 0 ? (
                        <Typography variant="body2" color="text.secondary">No assign rules</Typography>
                      ) : (
                        featureDetails.rules.filter(r => r.action === 'assign').sort((a,b) => a.priority - b.priority).map(r => (
                        <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1, mb: 1 }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1, flexWrap: 'wrap' }}>
                            <Chip size="small" label={`priority: ${r.priority}`} />
                            <Chip size="small" label={getSegmentName((r as any).segment_id)} color={(r as any).segment_id ? 'primary' : 'default'} />
                            {(r as any).segment_id && (r as any).is_customized && <Chip size="small" label="customized" color="warning" />}
                            {r.flag_variant_id && <Chip size="small" label={`target: ${getVariantName(r.flag_variant_id)}`} />}
                          </Box>
                          {r.conditions && (
                            <>
                              <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>Conditions:</Typography>
                              {renderConditionExpression(r.conditions)}
                            </>
                          )}
                        </Box>
                      ))
                    )}
                    </Collapse>
                  </Box>

                  <Box sx={{ border: '2px solid', borderColor: 'info.light', borderRadius: 1, p: 1 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                      <Typography variant="subtitle2">Include rules</Typography>
                      <IconButton size="small" onClick={() => toggleSection('include')}>
                        {expandedSections.include ? <ExpandLess /> : <ExpandMore />}
                      </IconButton>
                    </Box>
                    <Collapse in={expandedSections.include}>
                      {featureDetails.rules.filter(r => r.action === 'include').sort((a,b) => a.priority - b.priority).length === 0 ? (
                        <Typography variant="body2" color="text.secondary">No include rules</Typography>
                      ) : (
                        featureDetails.rules.filter(r => r.action === 'include').sort((a,b) => a.priority - b.priority).map(r => (
                        <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1, mb: 1 }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1, flexWrap: 'wrap' }}>
                            <Chip size="small" label={`priority: ${r.priority}`} />
                            <Chip size="small" label={getSegmentName((r as any).segment_id)} color={(r as any).segment_id ? 'primary' : 'default'} />
                            {(r as any).segment_id && (r as any).is_customized && <Chip size="small" label="customized" color="warning" />}
                          </Box>
                          {r.conditions && (
                            <>
                              <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>Conditions:</Typography>
                              {renderConditionExpression(r.conditions)}
                            </>
                          )}
                        </Box>
                      ))
                    )}
                    </Collapse>
                  </Box>

                  <Box sx={{ border: '2px solid', borderColor: 'error.light', borderRadius: 1, p: 1 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
                      <Typography variant="subtitle2">Exclude rules</Typography>
                      <IconButton size="small" onClick={() => toggleSection('exclude')}>
                        {expandedSections.exclude ? <ExpandLess /> : <ExpandMore />}
                      </IconButton>
                    </Box>
                    <Collapse in={expandedSections.exclude}>
                      {featureDetails.rules.filter(r => r.action === 'exclude').sort((a,b) => a.priority - b.priority).length === 0 ? (
                        <Typography variant="body2" color="text.secondary">No exclude rules</Typography>
                      ) : (
                        featureDetails.rules.filter(r => r.action === 'exclude').sort((a,b) => a.priority - b.priority).map(r => (
                        <Box key={r.id} sx={{ border: '1px solid', borderColor: 'divider', borderRadius: 1, p: 1, mb: 1 }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1, flexWrap: 'wrap' }}>
                            <Chip size="small" label={`priority: ${r.priority}`} />
                            <Chip size="small" label={getSegmentName((r as any).segment_id)} color={(r as any).segment_id ? 'primary' : 'default'} />
                            {(r as any).segment_id && (r as any).is_customized && <Chip size="small" label="customized" color="warning" />}
                          </Box>
                          {r.conditions && (
                            <>
                              <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>Conditions:</Typography>
                              {renderConditionExpression(r.conditions)}
                            </>
                          )}
                        </Box>
                      ))
                    )}
                    </Collapse>
                  </Box>
                </Box>
              )}
            </Box>
          </Box>
        ) : null}
      </DialogContent>
      <DialogActions>
        {featureDetails && (
          <>
            {canManage && (
              <Tooltip title={hasPendingChanges ? "Cannot delete: feature has pending changes awaiting approval" : ""}>
                <span>
                  <Button
                    onClick={() => {
                      if (deleteMutation.isPending) return;
                      if (window.confirm('Are you sure you want to delete this feature? This action cannot be undone.')) {
                        deleteMutation.mutate();
                      }
                    }}
                    color="error"
                    disabled={deleteMutation.isPending || hasPendingChanges}
                    size="small"
                  >
                    {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
                  </Button>
                </span>
              </Tooltip>
            )}
            {canManage && (
              <Tooltip title={hasPendingChanges ? "Cannot edit: feature has pending changes awaiting approval" : ""}>
                <span>
                  <Button 
                    onClick={() => setEditOpen(true)} 
                    color="secondary" 
                    size="small"
                    disabled={hasPendingChanges}
                  >
                    Edit
                  </Button>
                </span>
              </Tooltip>
            )}
          </>
        )}
        <Button onClick={onClose} size="small">Close</Button>
      </DialogActions>
      {/* Advanced edit dialog */}
      <EditFeatureDialog open={editOpen} onClose={() => setEditOpen(false)} featureDetails={featureDetails ?? null} environmentKey={environmentKey} />

      {/* Guard Response Handler */}
      <GuardResponseHandler
        pendingChange={guardResponse.pendingChange}
        conflictError={guardResponse.conflictError}
        forbiddenError={guardResponse.forbiddenError}
        onClose={() => setGuardResponse({})}
        onParentClose={onClose}
        onApprove={handleAutoApprove}
        approveLoading={approveMutation.isPending}
      />
    </Dialog>
  );
};

export default FeatureDetailsDialog;
