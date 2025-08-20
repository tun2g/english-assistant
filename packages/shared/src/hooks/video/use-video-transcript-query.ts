import { createQuery } from 'react-query-kit';
import { QUERY_KEY } from '../../constants';
import { getVideoTranscript } from '../../api/video-api';
import type { VideoTranscript } from '../../types';

export const useVideoTranscriptQuery = createQuery({
  queryKey: [QUERY_KEY.VIDEO.TRANSCRIPT],
  fetcher: (videoId: string): Promise<VideoTranscript> => {
    return getVideoTranscript({ videoUrl: videoId });
  },
  staleTime: 1000 * 60 * 10, // 10 minutes - transcripts don't change often
  retry: 2,
});
