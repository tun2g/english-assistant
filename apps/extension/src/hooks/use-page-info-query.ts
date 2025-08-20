import { useQuery } from '@tanstack/react-query';
import { QUERY_KEY } from '@english/shared';

interface PageInfo {
  url: string;
  title: string;
  isYouTube: boolean;
  videoId: string | null;
}

async function getCurrentPageInfo(): Promise<PageInfo> {
  return new Promise(resolve => {
    chrome.tabs.query({ active: true, currentWindow: true }, ([tab]) => {
      if (!tab) {
        resolve({
          url: '',
          title: '',
          isYouTube: false,
          videoId: null,
        });
        return;
      }

      const isYouTube = tab.url?.includes('youtube.com/watch') || false;
      const videoId = isYouTube ? new URL(tab.url!).searchParams.get('v') || null : null;

      resolve({
        url: tab.url || '',
        title: tab.title || '',
        isYouTube,
        videoId,
      });
    });
  });
}

export function usePageInfoQuery(options: { enabled?: boolean } = {}) {
  const {
    data: pageInfo,
    isLoading,
    error,
  } = useQuery({
    queryKey: [QUERY_KEY.PAGE_INFO],
    queryFn: getCurrentPageInfo,
    refetchOnWindowFocus: true,
    refetchInterval: 5000, // Refetch every 5 seconds to catch navigation
    enabled: options.enabled !== false, // Default to enabled unless explicitly disabled
  });

  return {
    pageInfo: pageInfo || {
      url: '',
      title: '',
      isYouTube: false,
      videoId: null,
    },
    isLoading,
    error,
  };
}
