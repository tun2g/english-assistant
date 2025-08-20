import React from 'react';
import { Button, Spinner } from '@english/ui';
import { useOAuthQuery } from '../../../hooks';

interface TranscriptButtonProps {
  videoId: string | null;
}

export function TranscriptButton({ videoId }: TranscriptButtonProps) {
  const { isAuthenticated, isLoading } = useOAuthQuery();

  const handleGetTranscript = async () => {
    if (!videoId) return;

    // Send message to content script to trigger transcript request
    const [tab] = await chrome.tabs.query({
      active: true,
      currentWindow: true,
    });
    if (tab.id) {
      await chrome.tabs.sendMessage(tab.id, {
        action: 'GET_TRANSCRIPT_WITH_AUTH',
      });
    }
  };

  if (!videoId) return null;

  return (
    <div className="p-4">
      <Button variant="default" onClick={handleGetTranscript} disabled={isLoading} className="w-full">
        {isLoading && <Spinner size="sm" className="mr-2" />}
        <span className="mr-2">{isAuthenticated ? '📝' : '🔒'}</span>
        {isAuthenticated ? 'Get Real Transcript' : 'Get Transcript (Auth Required)'}
      </Button>
    </div>
  );
}
