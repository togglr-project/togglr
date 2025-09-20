import React, { useState } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
  Button,
  Box,
  Typography,
  TableContainer,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  IconButton,
  Paper,
  CircularProgress
} from '@mui/material';

interface CreateProjectDialogProps {
  open: boolean;
  onClose: () => void;
  onCreateProject: (name: string, description: string) => void;
  isLoadingTeams: boolean;
}

const CreateProjectDialog: React.FC<CreateProjectDialogProps> = ({
  open,
  onClose,
  onCreateProject,
}) => {
  const [projectName, setProjectName] = useState('');
  const [projectDescription, setProjectDescription] = useState('');
  const [descriptionError, setDescriptionError] = useState<string>('');

  // Validate description length
  const validateDescription = (description: string) => {
    if (description.trim().length < 10) {
      setDescriptionError('Description must be at least 10 characters long');
      return false;
    } else {
      setDescriptionError('');
      return true;
    }
  };

  const handleDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setProjectDescription(value);
    validateDescription(value);
  };

  // Reset error when dialog opens
  React.useEffect(() => {
    if (open) {
      setDescriptionError('');
    }
  }, [open]);

  const handleCreate = () => {
    // Validate description before creating
    if (!validateDescription(projectDescription)) {
      return;
    }
    
    onCreateProject(projectName, projectDescription);
    // Reset form
    setProjectName('');
    setProjectDescription('');
    setDescriptionError('');
  };

  const handleCancel = () => {
    // Reset form
    setProjectName('');
    setProjectDescription('');
    setDescriptionError('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleCancel}>
      <DialogTitle sx={{ color: 'primary.main' }}>Create New Project</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Enter the name for the new project and select a team (optional).
        </DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Project Name"
          type="text"
          fullWidth
          variant="outlined"
          value={projectName}
          onChange={(e) => setProjectName(e.target.value)}
        />
        <TextField
          margin="dense"
          id="description"
          label="Project Description"
          type="text"
          fullWidth
          variant="outlined"
          value={projectDescription}
          onChange={handleDescriptionChange}
          multiline
          minRows={2}
          required
          error={!!descriptionError}
          helperText={descriptionError || 'Enter a detailed description of the project (minimum 10 characters)'}
        />
      </DialogContent>
      <DialogActions>
        <Button 
          onClick={handleCancel}
          color="primary"
        >
          Cancel
        </Button>
        <Button 
          onClick={handleCreate} 
          variant="contained"
          color="primary"
          disabled={!projectName.trim() || !projectDescription.trim() || !!descriptionError}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateProjectDialog;
