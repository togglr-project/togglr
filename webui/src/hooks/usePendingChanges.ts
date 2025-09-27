import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import apiClient from '../api/apiClient';
import type { 
  PendingChangeResponse, 
  PendingChangesListResponse,
  ApprovePendingChangeRequest,
  RejectPendingChangeRequest,
  CancelPendingChangeRequest
} from '../generated/api/client';

// Query keys
export const pendingChangesKeys = {
  all: ['pendingChanges'] as const,
  lists: () => [...pendingChangesKeys.all, 'list'] as const,
  list: (projectId: string, status?: string, environmentId?: number) => [...pendingChangesKeys.lists(), { projectId, status, environmentId }] as const,
  details: () => [...pendingChangesKeys.all, 'detail'] as const,
  detail: (id: string) => [...pendingChangesKeys.details(), id] as const,
};

// Hook for listing pending changes
export const usePendingChanges = (projectId: string, status?: string, environmentId?: number) => {
  return useQuery({
    queryKey: pendingChangesKeys.list(projectId, status, environmentId),
    queryFn: async (): Promise<PendingChangesListResponse> => {
      const response = await apiClient.listPendingChanges(
        environmentId,
        projectId,
        status as any, // Cast to the correct enum type
        undefined, // userId
        undefined, // page
        undefined, // perPage
        undefined, // sortBy
        undefined  // sortDesc
      );
      return response.data;
    },
    enabled: !!projectId,
  });
};

// Hook for getting a single pending change
export const usePendingChange = (id: string) => {
  return useQuery({
    queryKey: pendingChangesKeys.detail(id),
    queryFn: async (): Promise<PendingChangeResponse> => {
      const response = await apiClient.getPendingChange(id);
      return response.data;
    },
    enabled: !!id,
  });
};

// Hook for approving a pending change
export const useApprovePendingChange = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ 
      id, 
      request 
    }: { 
      id: string; 
      request: ApprovePendingChangeRequest 
    }) => {
      const response = await apiClient.approvePendingChange(id, request);
      return response.data;
    },
    onSuccess: (data, variables) => {
      // Invalidate and refetch pending changes lists
      queryClient.invalidateQueries({ queryKey: pendingChangesKeys.lists() });
      // Invalidate the specific pending change
      queryClient.invalidateQueries({ queryKey: pendingChangesKeys.detail(variables.id) });
    },
  });
};

// Hook for rejecting a pending change
export const useRejectPendingChange = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ 
      id, 
      request 
    }: { 
      id: string; 
      request: RejectPendingChangeRequest 
    }) => {
      const response = await apiClient.rejectPendingChange(id, request);
      return response.data;
    },
    onSuccess: (data, variables) => {
      // Invalidate and refetch pending changes lists
      queryClient.invalidateQueries({ queryKey: pendingChangesKeys.lists() });
      // Invalidate the specific pending change
      queryClient.invalidateQueries({ queryKey: pendingChangesKeys.detail(variables.id) });
    },
  });
};

// Hook for cancelling a pending change
export const useCancelPendingChange = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ 
      id, 
      request 
    }: { 
      id: string; 
      request: CancelPendingChangeRequest 
    }) => {
      const response = await apiClient.cancelPendingChange(id, request);
      return response.data;
    },
    onSuccess: (data, variables) => {
      // Invalidate and refetch pending changes lists
      queryClient.invalidateQueries({ queryKey: pendingChangesKeys.lists() });
      // Invalidate the specific pending change
      queryClient.invalidateQueries({ queryKey: pendingChangesKeys.detail(variables.id) });
    },
  });
};

// Hook for getting pending changes count (for menu indicator)
export const usePendingChangesCount = (projectId: string) => {
  return useQuery({
    queryKey: [...pendingChangesKeys.list(projectId, 'pending'), 'count'],
    queryFn: async (): Promise<number> => {
      const response = await apiClient.listPendingChanges(
        projectId,
        'pending' as any, // Cast to the correct enum type
        undefined, // userId
        undefined, // page
        undefined, // perPage
        undefined, // sortBy
        undefined  // sortDesc
      );
      return response.data.data.length;
    },
    enabled: !!projectId,
    refetchInterval: 30000, // Refetch every 30 seconds
  });
};
