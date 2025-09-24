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
import type { AuthCredentialsMethodEnum } from '../../generated/api/client';

interface ApprovalDialogProps {
  open: boolean;
  onClose: () => void;
  onApprove: (authMethod: AuthCredentialsMethodEnum, credential: string) => void;
  loading?: boolean;
  error?: string;
  title?: string;
  description?: string;
}

const ApprovalDialog: React.FC<ApprovalDialogProps> = ({
  open,
  onClose,
  onApprove,
  loading = false,
  error,
  title = 'Confirm Action',
  description = 'This action requires additional verification. Please provide your credentials.',
}) => {
  const { user } = useAuth();
  const [authMethod, setAuthMethod] = useState<AuthCredentialsMethodEnum>('password');
  const [credential, setCredential] = useState('');
  const [localError, setLocalError] = useState('');

  const handleSubmit = () => {
    if (!credential.trim()) {
      setLocalError('Please enter your credentials');
      return;
    }

    setLocalError('');
    onApprove(authMethod, credential);
  };

  const handleClose = () => {
    setCredential('');
    setLocalError('');
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
          <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
            Enter the 6-digit code from your authenticator app
          </Typography>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Cancel
        </Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained" 
          disabled={loading || !credential.trim()}
          startIcon={loading ? <CircularProgress size={20} /> : null}
        >
          {loading ? 'Verifying...' : 'Confirm'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ApprovalDialog;
