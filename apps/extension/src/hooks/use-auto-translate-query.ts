import { useQuery, useMutation } from '@tanstack/react-query';
import { QUERY_KEY, getQueryClient } from '@english/shared';
import { EXTENSION_STORAGE_KEYS } from '../shared/constants/storage-constants';

async function getAutoTranslateStatus(): Promise<boolean> {
  return new Promise(resolve => {
    chrome.storage.local.get([EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED], result => {
      resolve(result[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED] || false);
    });
  });
}

async function setAutoTranslateStatus(enabled: boolean): Promise<boolean> {
  return new Promise(resolve => {
    chrome.storage.local.set({ [EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]: enabled }, () => {
      resolve(enabled);
    });
  });
}

export function useAutoTranslateQuery(options: { enabled?: boolean } = {}) {
  const queryClient = getQueryClient();

  const {
    data: isEnabled,
    isLoading,
    error: queryError,
  } = useQuery({
    queryKey: [QUERY_KEY.AUTO_TRANSLATE],
    queryFn: getAutoTranslateStatus,
    staleTime: 1000 * 60 * 5, // 5 minutes
    enabled: options.enabled !== false, // Default to enabled unless explicitly disabled
  });

  const toggleMutation = useMutation({
    mutationFn: (enabled: boolean) => setAutoTranslateStatus(enabled),
    onSuccess: newStatus => {
      // Update cache immediately for optimistic UI
      queryClient.setQueryData([QUERY_KEY.AUTO_TRANSLATE], newStatus);
    },
  });

  return {
    isEnabled: isEnabled || false,
    isLoading,
    toggle: (enabled: boolean) => toggleMutation.mutateAsync(enabled),
    isToggling: toggleMutation.isPending,
    error: toggleMutation.error,
  };
}
