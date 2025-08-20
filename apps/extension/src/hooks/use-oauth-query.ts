import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { QUERY_KEY } from '@english/shared';
import { oauthService } from '../services/oauth-service';

export function useOAuthQuery() {
  const queryClient = useQueryClient();

  const {
    data: authStatus,
    isLoading: isCheckingAuth,
    error: authError,
  } = useQuery({
    queryKey: [QUERY_KEY.OAUTH.STATUS],
    queryFn: () => oauthService.checkAuthStatus(),
    refetchInterval: 30000, // Check every 30 seconds
    retry: (failureCount, error) => {
      // Don't retry as aggressively for network errors
      console.log('OAuth query retry attempt:', failureCount, error);
      return failureCount < 2;
    },
  });

  const connectMutation = useMutation({
    mutationFn: () => oauthService.initiateOAuth(),
    onSuccess: () => {
      // Invalidate auth status to trigger refetch
      queryClient.invalidateQueries({ queryKey: [QUERY_KEY.OAUTH.STATUS] });
    },
  });

  const disconnectMutation = useMutation({
    mutationFn: () => oauthService.revokeAuth(),
    onSuccess: () => {
      // Invalidate auth status to trigger refetch
      queryClient.invalidateQueries({ queryKey: [QUERY_KEY.OAUTH.STATUS] });
    },
  });

  return {
    isAuthenticated: authStatus?.authenticated || false,
    isLoading: isCheckingAuth || connectMutation.isPending || disconnectMutation.isPending,
    connect: connectMutation.mutateAsync,
    disconnect: disconnectMutation.mutateAsync,
    error: authError || connectMutation.error || disconnectMutation.error,
  };
}
