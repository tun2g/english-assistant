import React from 'react';
import { List } from 'framework7-react';

import { YouTubeVideoInfo } from './youtube-video-info';
import { AutoTranslateToggle } from './auto-translate-toggle';
import { TranscriptButton } from './transcript-button';
import { OAuthStatus } from '../oauth/oauth-status';
import { OAuthControls } from '../oauth/oauth-controls';

import type { PageInfo } from '../../../shared/types/extension-types';

interface YouTubeSectionProps {
  pageInfo: PageInfo;
  isOAuthAuthenticated: boolean;
  isOAuthLoading: boolean;
  isAutoTranslateEnabled: boolean;
  onOAuthConnect: () => Promise<void>;
  onOAuthDisconnect: () => Promise<void>;
  onAutoTranslateToggle: (enabled: boolean) => Promise<void>;
  onGetTranscript: () => Promise<void>;
}

export function YouTubeSection({
  pageInfo,
  isOAuthAuthenticated,
  isOAuthLoading,
  isAutoTranslateEnabled,
  onOAuthConnect,
  onOAuthDisconnect,
  onAutoTranslateToggle,
  onGetTranscript
}: YouTubeSectionProps) {
  return (
    <>
      <YouTubeVideoInfo pageInfo={pageInfo} />
      
      <OAuthStatus isAuthenticated={isOAuthAuthenticated} />
      
      <List>
        <AutoTranslateToggle 
          isEnabled={isAutoTranslateEnabled}
          onToggle={onAutoTranslateToggle}
        />
        
        <OAuthControls
          isAuthenticated={isOAuthAuthenticated}
          isLoading={isOAuthLoading}
          onConnect={onOAuthConnect}
          onDisconnect={onOAuthDisconnect}
        />

        <TranscriptButton
          isAuthenticated={isOAuthAuthenticated}
          isLoading={isOAuthLoading}
          videoId={pageInfo.videoId}
          onGetTranscript={onGetTranscript}
        />
      </List>
    </>
  );
}