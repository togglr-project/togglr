import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
} from '@mui/material';

interface ConfirmDiscardDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title?: string;
  message?: string;
}

const ConfirmDiscardDialog: React.FC<ConfirmDiscardDialogProps> = ({
  open,
  onClose,
  onConfirm,
  title = 'Discard Changes?',
  message = 'You have unsaved changes. Are you sure you want to discard them?',
}) => {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      aria-labelledby="discard-dialog-title"
      aria-describedby="discard-dialog-description"
    >
      <DialogTitle id="discard-dialog-title">
        {title}
      </DialogTitle>
      <DialogContent>
        <DialogContentText id="discard-dialog-description">
          {message}
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} size="small">
          Cancel
        </Button>
        <Button onClick={onConfirm} color="warning" variant="contained" size="small">
          Discard
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ConfirmDiscardDialog;
