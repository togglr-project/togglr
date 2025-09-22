import { useState, useEffect, useCallback } from 'react';

interface UseFormChangesOptions {
  initialData?: any;
  onDiscard?: () => void;
}

export const useFormChanges = ({ initialData, onDiscard }: UseFormChangesOptions = {}) => {
  const [hasChanges, setHasChanges] = useState(false);
  const [showDiscardDialog, setShowDiscardDialog] = useState(false);
  const [originalData, setOriginalData] = useState(initialData);

  // Update original data when initialData changes
  useEffect(() => {
    if (initialData) {
      setOriginalData(initialData);
      setHasChanges(false);
    }
  }, [initialData]);

  // Check if current data differs from original
  const checkForChanges = useCallback((currentData: any) => {
    if (!originalData || !currentData) {
      setHasChanges(false);
      return;
    }

    const hasChanges = JSON.stringify(currentData) !== JSON.stringify(originalData);
    setHasChanges(hasChanges);
  }, [originalData]);

  // Mark form as having changes
  const markAsChanged = useCallback(() => {
    setHasChanges(true);
  }, []);

  // Reset changes
  const resetChanges = useCallback(() => {
    setHasChanges(false);
    setShowDiscardDialog(false);
  }, []);

  // Handle discard confirmation
  const handleDiscard = useCallback(() => {
    if (onDiscard) {
      onDiscard();
    }
    resetChanges();
  }, [onDiscard, resetChanges]);

  // Show discard dialog
  const showDiscardConfirmation = useCallback(() => {
    if (hasChanges) {
      setShowDiscardDialog(true);
      return true; // Indicates that dialog was shown
    }
    return false; // No changes, safe to proceed
  }, [hasChanges]);

  // Close discard dialog
  const closeDiscardDialog = useCallback(() => {
    setShowDiscardDialog(false);
  }, []);

  return {
    hasChanges,
    showDiscardDialog,
    checkForChanges,
    markAsChanged,
    resetChanges,
    handleDiscard,
    showDiscardConfirmation,
    closeDiscardDialog,
  };
};
