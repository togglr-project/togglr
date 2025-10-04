import React from 'react';

interface LicenseGuardProps {
  children: React.ReactNode;
}

const LicenseGuard: React.FC<LicenseGuardProps> = ({ children }) => {
  // Empty stub - no license checking in open source version
  return <>{children}</>;
};

export default LicenseGuard;
