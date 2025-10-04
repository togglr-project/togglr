import React, { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';
import Layout from './Layout';
import ChangePasswordForm from './ChangePasswordForm';
import UserAgreementDialog from './UserAgreementDialog';

interface AuthenticatedLayoutProps {
  children: React.ReactNode;
  showBackButton?: boolean;
  backTo?: string;
}

const AuthenticatedLayout: React.FC<AuthenticatedLayoutProps> = ({ 
  children, 
  showBackButton = false, 
  backTo = '/dashboard' 
}) => {
  const { isAuthenticated, hasTmpPassword, user } = useAuth();
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);
  const [showUserAgreementDialog, setShowUserAgreementDialog] = useState(false);

  // Show password dialog when hasTmpPassword is true
  useEffect(() => {
    if (hasTmpPassword) {
      setShowPasswordDialog(true);
    } else {
      setShowPasswordDialog(false);
    }
  }, [hasTmpPassword]);

  // Show user agreement dialog when license_accepted is false (this is actually user agreement acceptance)
  useEffect(() => {
    if (user && !user.license_accepted) {
      setShowUserAgreementDialog(true);
    } else {
      setShowUserAgreementDialog(false);
    }
  }, [user]);

  // Handle password dialog close
  const handleClosePasswordDialog = () => {
    // The dialog will be closed automatically after successful password change
    // or when hasTmpPassword becomes false
    setShowPasswordDialog(false);
  };

  // Handle user agreement dialog close
  const handleCloseUserAgreementDialog = () => {
    // The dialog will be closed automatically after accepting the user agreement
    // or when the user data is updated with license_accepted=true
    setShowUserAgreementDialog(false);
  };

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <Layout showBackButton={showBackButton} backTo={backTo}>
      {children}
      <ChangePasswordForm open={showPasswordDialog} onClose={handleClosePasswordDialog} />
      <UserAgreementDialog open={showUserAgreementDialog} onClose={handleCloseUserAgreementDialog} />
    </Layout>
  );
};

export default AuthenticatedLayout; 