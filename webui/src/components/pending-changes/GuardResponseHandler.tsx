import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Alert,
  Typography,
  Box,
} from '@mui/material';
import { Assignment as ChangesIcon } from '@mui/icons-material';
import ApprovalDialog from './ApprovalDialog';
import { useFeatureNames, getEntityDisplayName } from '../../hooks/useFeatureNames';
import type { PendingChangeResponse, AuthCredentialsMethodEnum } from '../../generated/api/client';

interface GuardResponseHandlerProps {
  // For 202 Accepted - pending change created
  pendingChange?: PendingChangeResponse;
  // For 409 Conflict - entity locked by another pending change
  conflictError?: string;
  // For 403 Forbidden - user doesn't have permission to modify guarded feature
  forbiddenError?: string;
  // General error handling
  error?: string;
  onClose: () => void;
  // Optional: close parent dialog (e.g., EditFeatureDialog)
  onParentClose?: () => void;
  // Optional: allow immediate approval for single-user projects
  onApprove?: (authMethod: AuthCredentialsMethodEnum, credential: string, sessionId?: string) => void;
  approveLoading?: boolean;
}

const GuardResponseHandler: React.FC<GuardResponseHandlerProps> = ({
  pendingChange,
  conflictError,
  forbiddenError,
  error,
  onClose,
  onParentClose,
  onApprove,
  approveLoading = false,
}) => {
  const [showApprovalDialog, setShowApprovalDialog] = useState(false);

  // Get feature names for display
  const featureIds = pendingChange?.change.entities
    ?.filter(entity => entity.entity === 'feature')
    ?.map(entity => entity.entity_id) || [];
  
  const { data: featureNames } = useFeatureNames(featureIds, pendingChange?.project_id);

  const handleApprove = (authMethod: AuthCredentialsMethodEnum, credential: string, sessionId?: string) => {
    if (onApprove) {
      onApprove(authMethod, credential, sessionId);
    }
  };

  // 202 Accepted - Change is pending approval
  if (pendingChange) {
    const isSingleUserProject = pendingChange.change.meta?.single_user_project === true;
    
    return (
      <>
        <Dialog open={true} onClose={onClose} maxWidth="md" fullWidth>
          <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <ChangesIcon color="warning" />
            Change Pending Approval
          </DialogTitle>
          <DialogContent>
            <Alert severity="info" sx={{ mb: 2 }}>
              {isSingleUserProject 
                ? "Your change has been submitted. Since you are the only active user in this project, you can approve it immediately by verifying your credentials."
                : "Your change has been submitted and is pending approval from an authorized user."
              }
            </Alert>
            
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6" gutterBottom>
                Pending Change Details
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Change ID: {pendingChange.id}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Requested by: {pendingChange.requested_by}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Entities: {pendingChange.change.entities.length} item(s)
              </Typography>
            </Box>

            <Box>
              <Typography variant="body1">
                The following changes will be applied once approved:
              </Typography>
              {pendingChange.change.entities.map((entity, index) => (
                <Box key={index} sx={{ mt: 1, p: 2, bgcolor: 'background.paper', border: 1, borderColor: 'divider', borderRadius: 1 }}>
                  <Typography variant="subtitle2">
                    {getEntityDisplayName(entity, featureNames)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Action: {entity.action}
                  </Typography>
                  {Object.keys(entity.changes).length > 0 && (
                    <Box sx={{ mt: 1 }}>
                      <Typography variant="caption" color="text.secondary">
                        Fields to be changed: {Object.keys(entity.changes).join(', ')}
                      </Typography>
                    </Box>
                  )}
                </Box>
              ))}
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => {
              onClose();
              onParentClose?.();
            }}>
              Close
            </Button>
            {onApprove && isSingleUserProject && (
              <Button 
                variant="contained" 
                onClick={() => setShowApprovalDialog(true)}
                disabled={approveLoading}
              >
                Approve Now
              </Button>
            )}
          </DialogActions>
        </Dialog>

        {onApprove && isSingleUserProject && (
          <ApprovalDialog
            open={showApprovalDialog}
            onClose={() => setShowApprovalDialog(false)}
            onApprove={handleApprove}
            loading={approveLoading}
            title="Auto-Approve Change"
            description="Since you are the only active user in this project, you can approve this change immediately by verifying your credentials."
            pendingChangeId={pendingChange.id}
          />
        )}
      </>
    );
  }

  // 409 Conflict - Entity locked by another pending change
  if (conflictError) {
    return (
      <Dialog open={true} onClose={onClose} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <ChangesIcon color="error" />
          Change Conflict
        </DialogTitle>
        <DialogContent>
          <Alert severity="error" sx={{ mb: 2 }}>
            Cannot apply changes - entity is locked by another pending change.
          </Alert>
          <Typography variant="body1">
            {conflictError}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
            Please wait for the existing pending change to be processed or contact an administrator.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => {
            onClose();
            onParentClose?.();
          }} variant="contained">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  // 403 Forbidden - User doesn't have permission to modify guarded feature
  if (forbiddenError) {
    return (
      <Dialog open={true} onClose={onClose} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <ChangesIcon color="error" />
          Access Denied
        </DialogTitle>
        <DialogContent>
          <Alert severity="error" sx={{ mb: 2 }}>
            You don't have permission to modify this guarded feature.
          </Alert>
          <Typography variant="body1">
            {forbiddenError}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
            Contact your project administrator to request access or assign the appropriate role.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => {
            onClose();
            onParentClose?.();
          }} variant="contained">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  // General error
  if (error) {
    return (
      <Dialog open={true} onClose={onClose} maxWidth="sm" fullWidth>
        <DialogTitle>Error</DialogTitle>
        <DialogContent>
          <Alert severity="error">
            {error}
          </Alert>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => {
            onClose();
            onParentClose?.();
          }} variant="contained">
            Close
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  return null;
};

export default GuardResponseHandler;
