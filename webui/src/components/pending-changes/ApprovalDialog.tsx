import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Box,
  Typography,
  CircularProgress,
} from '@mui/material';
import { useAuth } from '../../auth/AuthContext';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '../../api/apiClient';
import type { AuthCredentialsMethodEnum } from '../../generated/api/client';

interface ApprovalDialogProps {
  open: boolean;
  onClose: () => void;
  onApprove: (authMethod: AuthCredentialsMethodEnum, credential: string, sessionId?: string) => void;
  loading?: boolean;
  error?: string;
  title?: string;
  description?: string;
  pendingChangeId?: string;
}

const ApprovalDialog: React.FC<ApprovalDialogProps> = ({
  open,
  onClose,
  onApprove,
  loading = false,
  error,
  title = 'Confirm Action',
  description = 'This action requires additional verification. Please provide your credentials.',
  pendingChangeId,
}) => {
  const { user } = useAuth();
  const [authMethod, setAuthMethod] = useState<AuthCredentialsMethodEnum>('password');
  const [credential, setCredential] = useState('');
  const [localError, setLocalError] = useState('');
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [totpInitiated, setTotpInitiated] = useState(false);

  // Mutation for initiating TOTP approval
  const initiateTOTPMutation = useMutation({
    mutationFn: async () => {
      if (!pendingChangeId || !user?.id) {
        throw new Error('Missing pending change ID or user ID');
      }
      const response = await apiClient.initiateTOTPApproval(pendingChangeId, {
        approver_user_id: user.id,
      });
      return response.data;
    },
    onSuccess: (data) => {
      setSessionId(data.session_id);
      setTotpInitiated(true);
      setLocalError('');
    },
    onError: (error: unknown) => {
      const errorMessage = (error as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || 'Failed to initiate TOTP approval';
      setLocalError(errorMessage);
    },
  });

  const handleSubmit = () => {
    if (!credential.trim()) {
      setLocalError('Please enter your credentials');
      return;
    }

    setLocalError('');
    onApprove(authMethod, credential, sessionId || undefined);
  };

  const handleInitiateTOTP = () => {
    setLocalError('');
    initiateTOTPMutation.mutate();
  };

  const handleClose = () => {
    setCredential('');
    setLocalError('');
    setSessionId(null);
    setTotpInitiated(false);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Box sx={{ mb: 2 }}>
          <Typography variant="body2" color="text.secondary">
            {description}
          </Typography>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {localError && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {localError}
          </Alert>
        )}

        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel>Verification Method</InputLabel>
          <Select
            value={authMethod}
            label="Verification Method"
            onChange={(e) => setAuthMethod(e.target.value as AuthCredentialsMethodEnum)}
          >
            <MenuItem value="password">Password</MenuItem>
            <MenuItem value="totp">TOTP Code</MenuItem>
          </Select>
        </FormControl>

        <TextField
          fullWidth
          type={authMethod === 'password' ? 'password' : 'text'}
          label={authMethod === 'password' ? 'Password' : 'TOTP Code'}
          value={credential}
          onChange={(e) => setCredential(e.target.value)}
          disabled={loading}
          placeholder={
            authMethod === 'password' 
              ? 'Enter your password' 
              : 'Enter your 6-digit TOTP code'
          }
          inputProps={{
            maxLength: authMethod === 'totp' ? 6 : undefined,
          }}
        />

        {authMethod === 'totp' && (
          <Box sx={{ mt: 1 }}>
            {!totpInitiated ? (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Button
                  variant="outlined"
                  size="small"
                  onClick={handleInitiateTOTP}
                  disabled={initiateTOTPMutation.isPending}
                  startIcon={initiateTOTPMutation.isPending ? <CircularProgress size={16} /> : null}
                >
                  {initiateTOTPMutation.isPending ? 'Sending...' : 'Initiate'}
                </Button>
                <Typography variant="caption" color="text.secondary">
                  Click to initiate TOTP approval session
                </Typography>
              </Box>
            ) : (
              <Typography variant="caption" color="success.main" sx={{ display: 'block' }}>
                ✓ TOTP session initiated. Enter the 6-digit code from your authenticator app
              </Typography>
            )}
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Cancel
        </Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained" 
          disabled={loading || !credential.trim() || (authMethod === 'totp' && !totpInitiated)}
          startIcon={loading ? <CircularProgress size={20} /> : null}
        >
          {loading ? 'Verifying...' : 'Confirm'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ApprovalDialog;
