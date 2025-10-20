import React, { useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Typography,
  CircularProgress,
  Alert,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
} from '@mui/material';
import {
  // Add as AddIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { ProjectSetting } from '../../generated/api/client';
import ProjectSettingFormDialog from './ProjectSettingFormDialog';
import ProjectSettingCard from './ProjectSettingCard';
import { useNotification } from '../../App';
import { isPredefinedSetting } from '../../constants/projectSettings';

interface ProjectSettingsListProps {
  projectId: string;
}

const ProjectSettingsList: React.FC<ProjectSettingsListProps> = ({ projectId }) => {
  const { showNotification } = useNotification();
  const queryClient = useQueryClient();
  // const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedSetting, setSelectedSetting] = useState<ProjectSetting | null>(null);

  const { data: settingsData, isLoading, error } = useQuery({
    queryKey: ['project-settings', projectId],
    queryFn: async () => {
      const response = await apiClient.listProjectSettings(projectId);
      return response.data;
    },
    enabled: !!projectId,
  });

  // const deleteMutation = useMutation({
    // mutationFn: async (settingName: string) => {
    //   await apiClient.updateProjectSetting(projectId, settingName, { value: '' });
    // },
    // onSuccess: () => {
    //   queryClient.invalidateQueries({ queryKey: ['project-settings', projectId] });
    //   showNotification('Setting deleted', 'success');
    //   setDeleteDialogOpen(false);
    //   setSelectedSetting(null);
    // },
    // onError: (error: any) => {
    //   showNotification(`Error deleting setting: ${error.message}`, 'error');
    // },
  // });

  const settings = settingsData?.data || [];

  // const handleCreateSetting = () => {
  //   setSelectedSetting(null);
  //   setCreateDialogOpen(true);
  // };

  const handleEditSetting = (setting: ProjectSetting) => {
    setSelectedSetting(setting);
    setEditDialogOpen(true);
  };

  const handleDeleteSetting = (setting: ProjectSetting) => {
    if (isPredefinedSetting(setting.name)) {
      showNotification('Cannot delete predefined settings', 'error');
      return;
    }
    setSelectedSetting(setting);
    setDeleteDialogOpen(true);
  };

  // const handleDeleteConfirm = () => {
  //   if (selectedSetting) {
  //     deleteMutation.mutate(selectedSetting.name);
  //   }
  // };

  // const handleCreateSuccess = () => {
  //   setCreateDialogOpen(false);
  //   queryClient.invalidateQueries({ queryKey: ['project-settings', projectId] });
  //   showNotification('Setting created', 'success');
  // };

  const handleEditSuccess = () => {
    setEditDialogOpen(false);
    setSelectedSetting(null);
    queryClient.invalidateQueries({ queryKey: ['project-settings', projectId] });
    showNotification('Setting updated', 'success');
  };

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Error loading project settings
      </Alert>
    );
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h6" component="h2">
          Project Settings
        </Typography>
        {/*<Button*/}
        {/*  variant="contained"*/}
        {/*  startIcon={<AddIcon />}*/}
        {/*  onClick={handleCreateSetting}*/}
        {/*>*/}
        {/*  Add Setting*/}
        {/*</Button>*/}
      </Box>

      {settings.length === 0 ? (
        <Card>
          <CardContent sx={{ textAlign: 'center', py: 4 }}>
            <SettingsIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" gutterBottom>
              No settings found
            </Typography>
            {/*<Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>*/}
            {/*  Create the first setting for this project*/}
            {/*</Typography>*/}
            {/*<Button*/}
            {/*  variant="outlined"*/}
            {/*  startIcon={<AddIcon />}*/}
            {/*  onClick={handleCreateSetting}*/}
            {/*>*/}
            {/*  Add Setting*/}
            {/*</Button>*/}
          </CardContent>
        </Card>
      ) : (
        <Grid container spacing={2}>
          {settings.map((setting) => (
            <Grid item xs={12} md={6} lg={4} key={setting.id}>
              <ProjectSettingCard
                setting={setting}
                onEdit={() => handleEditSetting(setting)}
                onDelete={() => handleDeleteSetting(setting)}
              />
            </Grid>
          ))}
        </Grid>
      )}

      {/*<ProjectSettingFormDialog*/}
      {/*  open={createDialogOpen}*/}
      {/*  onClose={() => setCreateDialogOpen(false)}*/}
      {/*  onSubmit={handleCreateSuccess}*/}
      {/*  projectId={projectId}*/}
      {/*  mode="create"*/}
      {/*/>*/}

      <ProjectSettingFormDialog
        open={editDialogOpen}
        onClose={() => setEditDialogOpen(false)}
        onSubmit={handleEditSuccess}
        projectId={projectId}
        mode="edit"
        initialData={selectedSetting}
      />

      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Setting</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete setting "{selectedSetting?.name}"? 
            This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>
            Cancel
          </Button>
          {/*<Button*/}
          {/*  onClick={handleDeleteConfirm}*/}
          {/*  color="error"*/}
          {/*  disabled={deleteMutation.isPending}*/}
          {/*>*/}
          {/*  {deleteMutation.isPending ? <CircularProgress size={20} /> : 'Delete'}*/}
          {/*</Button>*/}
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default ProjectSettingsList;
