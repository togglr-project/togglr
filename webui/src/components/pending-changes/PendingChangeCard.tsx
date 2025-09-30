import React, { useState } from 'react';
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Box,
  Chip,
  Divider,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Tooltip,
  Collapse,
} from '@mui/material';
import {
  Check as CheckIcon,
  Close as CloseIcon,
  Cancel as CancelIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Person as PersonIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material';
import { format } from 'date-fns';
import { useAuth } from '../../auth/AuthContext';
import { useApprovePendingChange, useRejectPendingChange, useCancelPendingChange } from '../../hooks/usePendingChanges';
import { useFeatureNames, getEntityDisplayName } from '../../hooks/useFeatureNames';
import { useQueryClient } from '@tanstack/react-query';
import ApprovalDialog from './ApprovalDialog';
import type { PendingChangeResponse, AuthCredentialsMethodEnum } from '../../generated/api/client';

interface PendingChangeCardProps {
  pendingChange: PendingChangeResponse;
  onStatusChange?: () => void;
}

const PendingChangeCard: React.FC<PendingChangeCardProps> = ({
  pendingChange,
  onStatusChange,
}) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [showDetails, setShowDetails] = useState(false);
  const [showApprovalDialog, setShowApprovalDialog] = useState(false);
  const [showRejectDialog, setShowRejectDialog] = useState(false);
  const [rejectReason, setRejectReason] = useState('');
  const [authMethod, setAuthMethod] = useState<AuthCredentialsMethodEnum>('password');

  // Get feature names for display
  const featureIds = pendingChange.change.entities
    ?.filter(entity => entity.entity === 'feature')
    ?.map(entity => entity.entity_id) || [];
  
  const { data: featureNames } = useFeatureNames(featureIds, pendingChange.project_id, 'prod');

  const approveMutation = useApprovePendingChange();
  const rejectMutation = useRejectPendingChange();
  const cancelMutation = useCancelPendingChange();

  const handleApprove = (method: AuthCredentialsMethodEnum, credential: string, sessionId?: string) => {
    if (!user) return;

    console.log('Approving pending change:', {
      id: pendingChange.id,
      approver_user_id: user.id,
      approver_name: user.username,
      method,
      sessionId
    });

    approveMutation.mutate(
      {
        id: pendingChange.id,
        request: {
          approver_user_id: user.id,
          approver_name: user.username,
          auth: {
            method,
            credential,
            ...(sessionId && { session_id: sessionId }),
          },
        },
      },
      {
        onSuccess: () => {
          setShowApprovalDialog(false);
          
          // Invalidate all related caches
          queryClient.invalidateQueries({ queryKey: ['feature-details'] });
          queryClient.invalidateQueries({ queryKey: ['project-features'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes', pendingChange.project_id] });
          
          onStatusChange?.();
        },
      }
    );
  };

  const handleReject = () => {
    if (!user) return;

    rejectMutation.mutate(
      {
        id: pendingChange.id,
        request: {
          rejected_by: user.username,
          reason: rejectReason,
        },
      },
      {
        onSuccess: () => {
          setShowRejectDialog(false);
          setRejectReason('');
          
          // Invalidate all related caches
          queryClient.invalidateQueries({ queryKey: ['feature-details'] });
          queryClient.invalidateQueries({ queryKey: ['project-features'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes', pendingChange.project_id] });
          
          onStatusChange?.();
        },
      }
    );
  };

  const handleCancel = () => {
    if (!user) return;

    cancelMutation.mutate(
      {
        id: pendingChange.id,
        request: {
          cancelled_by: user.username,
        },
      },
      {
        onSuccess: () => {
          // Invalidate all related caches
          queryClient.invalidateQueries({ queryKey: ['feature-details'] });
          queryClient.invalidateQueries({ queryKey: ['project-features'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes', pendingChange.project_id] });
          
          onStatusChange?.();
        },
      }
    );
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'warning';
      case 'approved':
        return 'success';
      case 'rejected':
        return 'error';
      case 'cancelled':
        return 'default';
      default:
        return 'default';
    }
  };

  const formatChanges = (changes: Record<string, { old: unknown; new: unknown }>) => {
    return Object.entries(changes).map(([field, change]) => (
      <Box key={field} sx={{ mb: 1 }}>
        <Typography variant="body2" fontWeight="medium">
          {field}:
        </Typography>
        <Box sx={{ ml: 2 }}>
          <Typography variant="body2" color="error">
            - {JSON.stringify(change.old)}
          </Typography>
          <Typography variant="body2" color="success.main">
            + {JSON.stringify(change.new)}
          </Typography>
        </Box>
      </Box>
    ));
  };

  // Check if user can approve (not the requester, and has appropriate permissions)
  // Temporarily simplified for testing - allow all users except the requester
  const canApprove = user && user.id !== pendingChange.request_user_id;
  
  // Check if user can cancel (is the requester)
  const canCancel = user && user.id === pendingChange.request_user_id;

  // Debug info (remove in production)
  console.log('PendingChangeCard debug:', {
    userId: user?.id,
    requestUserId: pendingChange.request_user_id,
    isSuperuser: user?.is_superuser,
    projectPermissions: user?.project_permissions?.[pendingChange.project_id],
    canApprove,
    canCancel,
    status: pendingChange.status
  });

  return (
    <>
      <Card sx={{ mb: 2 }}>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
            <Box>
              <Typography variant="h6" component="div">
                Change Request #{pendingChange.id.slice(0, 8)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Requested by {pendingChange.requested_by}
              </Typography>
            </Box>
            <Chip 
              label={pendingChange.status} 
              color={getStatusColor(pendingChange.status) as any}
              size="small"
            />
          </Box>

          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <PersonIcon fontSize="small" color="action" />
              <Typography variant="body2" color="text.secondary">
                {pendingChange.requested_by}
              </Typography>
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <ScheduleIcon fontSize="small" color="action" />
              <Typography variant="body2" color="text.secondary">
                {format(new Date(pendingChange.created_at), 'MMM dd, yyyy HH:mm')}
              </Typography>
            </Box>
          </Box>

          <Typography variant="body2" sx={{ mb: 2 }}>
            {pendingChange.change.entities.length} entity(ies) to be modified
          </Typography>

          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Button
              size="small"
              onClick={() => setShowDetails(!showDetails)}
              endIcon={showDetails ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            >
              {showDetails ? 'Hide Details' : 'Show Details'}
            </Button>
          </Box>

          <Collapse in={showDetails}>
            <Box sx={{ mt: 2 }}>
              <Divider sx={{ mb: 2 }} />
              {pendingChange.change.entities.map((entity, index) => (
                <Box key={index} sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" gutterBottom>
                    {getEntityDisplayName(entity, featureNames)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                    Action: {entity.action}
                  </Typography>
                  {Object.keys(entity.changes).length > 0 && (
                    <Box>
                      <Typography variant="body2" fontWeight="medium" sx={{ mb: 1 }}>
                        Changes:
                      </Typography>
                      {formatChanges(entity.changes)}
                    </Box>
                  )}
                </Box>
              ))}
            </Box>
          </Collapse>
        </CardContent>

        {pendingChange.status === 'pending' && (
          <CardActions>
            {canApprove && (
              <Button
                variant="contained"
                color="success"
                startIcon={<CheckIcon />}
                onClick={() => setShowApprovalDialog(true)}
                disabled={approveMutation.isPending}
              >
                Approve
              </Button>
            )}
            {canApprove && (
              <Button
                variant="outlined"
                color="error"
                startIcon={<CloseIcon />}
                onClick={() => setShowRejectDialog(true)}
                disabled={rejectMutation.isPending}
              >
                Reject
              </Button>
            )}
            {canCancel && (
              <Button
                variant="outlined"
                startIcon={<CancelIcon />}
                onClick={handleCancel}
                disabled={cancelMutation.isPending}
              >
                Cancel
              </Button>
            )}
          </CardActions>
        )}
      </Card>

      <ApprovalDialog
        open={showApprovalDialog}
        onClose={() => setShowApprovalDialog(false)}
        onApprove={handleApprove}
        loading={approveMutation.isPending}
        error={approveMutation.error?.message}
        title="Approve Change Request"
        description="Please verify your identity to approve this change request."
        pendingChangeId={pendingChange.id}
      />

      <Dialog open={showRejectDialog} onClose={() => setShowRejectDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Reject Change Request</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            multiline
            rows={3}
            label="Rejection Reason"
            value={rejectReason}
            onChange={(e) => setRejectReason(e.target.value)}
            placeholder="Please provide a reason for rejecting this change request..."
            sx={{ mt: 1 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowRejectDialog(false)}>Cancel</Button>
          <Button
            onClick={handleReject}
            variant="contained"
            color="error"
            disabled={!rejectReason.trim() || rejectMutation.isPending}
          >
            {rejectMutation.isPending ? 'Rejecting...' : 'Reject'}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default PendingChangeCard;
