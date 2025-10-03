import React, { useState, useEffect } from 'react';
import { 
  Button, 
  CircularProgress,
  Alert,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Paper,
  Tabs,
  Tab,
  Box
} from '@mui/material';
import { useAuth } from '../auth/AuthContext';
import apiClient from '../api/apiClient';
import { fetchTextFile } from '../utils/fileUtils';

interface UserAgreementDialogProps {
  open: boolean;
  onClose: () => void;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`agreement-tabpanel-${index}`}
      aria-labelledby={`agreement-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const UserAgreementDialog: React.FC<UserAgreementDialogProps> = ({ open, onClose }) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { logout, updateUserData } = useAuth();
  const [tabValue, setTabValue] = useState(0);
  const [englishText, setEnglishText] = useState<string>('Loading...');
  const [russianText, setRussianText] = useState<string>('Загрузка...');
  const [isLoadingTexts, setIsLoadingTexts] = useState(true);

  useEffect(() => {
    const loadAgreementTexts = async () => {
      setIsLoadingTexts(true);
      try {
        const [enText, ruText] = await Promise.all([
          fetchTextFile('/USER_AGREEMENT_EN.txt'),
          fetchTextFile('/USER_AGREEMENT_RU.txt')
        ]);
        setEnglishText(enText);
        setRussianText(ruText);
      } catch (err) {
        console.error('Failed to load user agreement texts:', err);
        setError('Failed to load user agreement texts. Please try again.');
      } finally {
        setIsLoadingTexts(false);
      }
    };

    if (open) {
      loadAgreementTexts();
    }
  }, [open]);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleAccept = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      await apiClient.updateLicenseAcceptance({ accepted: true });
      // Update user data in AuthContext to reflect the user agreement acceptance
      await updateUserData();
      onClose();
    } catch (err) {
      console.error('Failed to accept user agreement:', err);
      setError('Failed to update user agreement acceptance status. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDecline = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      await apiClient.updateLicenseAcceptance({ accepted: false });
      // Logout the user after declining the user agreement
      logout();
    } catch (err) {
      console.error('Failed to decline user agreement:', err);
      setError('Failed to update user agreement acceptance status. Please try again.');
      setIsLoading(false);
    }
  };

  return (
    <Dialog 
      open={open} 
      onClose={() => {}} // Empty function to prevent closing by clicking outside
      maxWidth="md" 
      fullWidth
      disableEscapeKeyDown // Prevent closing with Escape key
    >
      <DialogTitle>
        <Typography variant="h5" component="div" align="center">
          User Agreement
        </Typography>
      </DialogTitle>
      <DialogContent>
        <Typography variant="body1" sx={{ mb: 2 }}>
          Please read and accept the following user agreement to continue using the application.
        </Typography>
        
        {error && (
          <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={tabValue} onChange={handleTabChange} aria-label="user agreement language tabs">
            <Tab label="English" id="agreement-tab-0" aria-controls="agreement-tabpanel-0" />
            <Tab label="Русский" id="agreement-tab-1" aria-controls="agreement-tabpanel-1" />
          </Tabs>
        </Box>
        
        <Paper 
          elevation={0} 
          variant="outlined" 
          sx={{ 
            maxHeight: '400px', 
            overflow: 'auto',
            mb: 2,
            backgroundColor: 'background.paper',
            color: 'text.primary'
          }}
        >
          {isLoadingTexts ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
              <CircularProgress />
            </Box>
          ) : (
            <>
              <TabPanel value={tabValue} index={0}>
                <Typography variant="body2" component="div" sx={{ whiteSpace: 'pre-line', color: 'text.primary' }}>
                  {englishText}
                </Typography>
              </TabPanel>
              <TabPanel value={tabValue} index={1}>
                <Typography variant="body2" component="div" sx={{ whiteSpace: 'pre-line', color: 'text.primary' }}>
                  {russianText}
                </Typography>
              </TabPanel>
            </>
          )}
        </Paper>
      </DialogContent>
      <DialogActions sx={{ justifyContent: 'space-between', px: 3, pb: 3 }}>
        <Button 
          onClick={handleDecline} 
          disabled={isLoading || isLoadingTexts}
          variant="outlined"
          color="error"
        >
          {isLoading ? <CircularProgress size={24} /> : 'Decline'}
        </Button>
        <Button 
          onClick={handleAccept} 
          disabled={isLoading || isLoadingTexts}
          variant="contained"
          color="primary"
        >
          {isLoading ? <CircularProgress size={24} /> : 'Accept'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default UserAgreementDialog;