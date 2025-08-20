import { createQuery } from 'react-query-kit';
import { QUERY_KEY } from '../../constants';

interface VideoInfo {
  id: string;
  title: string;
  url: string;
  duration?: number;
  thumbnail?: string;
}

export const useVideoInfoQuery = createQuery({
  queryKey: [QUERY_KEY.VIDEO.INFO],
  fetcher: async (videoId: string): Promise<VideoInfo> => {
    // This would typically call a YouTube API or backend service
    // For now, return basic info that can be extracted from the page
    return {
      id: videoId,
      title: document.title,
      url: window.location.href,
    };
  },
  staleTime: 1000 * 60 * 30, // 30 minutes
});
