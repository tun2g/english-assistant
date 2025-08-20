import { createMutation } from 'react-query-kit';
import { getQueryClient } from '../../lib/query-client';
import { QUERY_KEY } from '../../constants';
import { saveVideoProgress } from '../../api/video-api';
import type { VideoProgressRequest } from '../../types';

export const useVideoProgressMutation = createMutation({
  mutationFn: (data: VideoProgressRequest): Promise<void> => saveVideoProgress(data),
  onSuccess: () => {
    const queryClient = getQueryClient();
    // Invalidate related queries
    queryClient.invalidateQueries({
      queryKey: [QUERY_KEY.LEARNING.PROGRESS],
    });
  },
});
