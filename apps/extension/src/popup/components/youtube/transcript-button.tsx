import React from 'react';
import { Button, ListItem } from 'framework7-react';

interface TranscriptButtonProps {
  isAuthenticated: boolean;
  isLoading: boolean;
  videoId: string | null;
  onGetTranscript: () => Promise<void>;
}

export function TranscriptButton({ 
  isAuthenticated, 
  isLoading, 
  videoId, 
  onGetTranscript 
}: TranscriptButtonProps) {
  if (!videoId) return null;

  return (
    <ListItem>
      <Button 
        fill
        color="green"
        onClick={onGetTranscript}
        disabled={isLoading}
        style={{ width: '100%' }}
      >
        {isAuthenticated 
          ? '📝 Get Real Transcript' 
          : '🔒 Get Transcript (Auth Required)'
        }
      </Button>
    </ListItem>
  );
}