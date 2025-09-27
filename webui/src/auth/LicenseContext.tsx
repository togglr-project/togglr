import React, { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import { useAuth } from './AuthContext';
import apiClient from '../api/apiClient';
import type { LicenseStatusResponse } from '../generated/api/client';

interface LicenseContextType {
  licenseStatus: LicenseStatusResponse | null;
  isLicenseValid: boolean;
  isLoading: boolean;
  error: string | null;
  checkLicenseStatus: () => Promise<void>;
}

const LicenseContext = createContext<LicenseContextType | undefined>(undefined);

interface LicenseProviderProps {
  children: ReactNode;
}

export const LicenseProvider: React.FC<LicenseProviderProps> = ({ children }) => {
  const [licenseStatus, setLicenseStatus] = useState<LicenseStatusResponse | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const { isAuthenticated } = useAuth();

  const checkLicenseStatus = useCallback(async () => {
    if (!isAuthenticated) {
      return;
    }

    try {
      setIsLoading(true);
      setError(null);
      const response = await apiClient.getLicenseStatus();
      setLicenseStatus(response.data);
    } catch (err) {
      console.error('Failed to check license status:', err);
      setError('Failed to check license status');
      setLicenseStatus(null);
    } finally {
      setIsLoading(false);
    }
  }, [isAuthenticated]);

  // Check license status when user becomes authenticated
  useEffect(() => {
    if (isAuthenticated) {
      checkLicenseStatus();
    } else {
      setLicenseStatus(null);
      setError(null);
    }
  }, [isAuthenticated, checkLicenseStatus]);

  // Determine if license is valid
  const isLicenseValid = licenseStatus?.license?.is_valid ?? false;

  const value: LicenseContextType = {
    licenseStatus,
    isLicenseValid,
    isLoading,
    error,
    checkLicenseStatus,
  };

  return (
    <LicenseContext.Provider value={value}>
      {children}
    </LicenseContext.Provider>
  );
};

export const useLicense = (): LicenseContextType => {
  const context = useContext(LicenseContext);
  if (context === undefined) {
    throw new Error('useLicense must be used within a LicenseProvider');
  }
  return context;
}; 