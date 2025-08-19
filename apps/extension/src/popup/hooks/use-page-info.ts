import { useState, useEffect } from 'react';
import { EXTENSION_MESSAGES } from '../../shared/constants/extension-constants';
import type { PageInfo } from '../../shared/types/extension-types';

export function usePageInfo() {
  const [pageInfo, setPageInfo] = useState<PageInfo | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchPageInfo = async () => {
      try {
        const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
        if (tab.id) {
          const response = await chrome.tabs.sendMessage(tab.id, { 
            action: EXTENSION_MESSAGES.GET_PAGE_INFO 
          });
          if (response?.success) {
            setPageInfo(response.data);
          }
        }
      } catch (error) {
        console.error('Failed to get page info:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchPageInfo();
  }, []);

  return { pageInfo, isLoading };
}