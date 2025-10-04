import React from 'react';

interface LicenseDialogProps {
  open: boolean;
  onClose: () => void;
}

const LicenseDialog: React.FC<LicenseDialogProps> = () => {
  // Empty stub - no license checking in open source version
  return null;
};

export default LicenseDialog;
