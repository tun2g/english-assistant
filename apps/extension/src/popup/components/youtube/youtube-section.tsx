import React from 'react';
import { YouTubeVideoInfo } from './youtube-video-info';
import { AutoTranslateToggle } from './auto-translate-toggle';
import { TranscriptButton } from './transcript-button';
import { OAuthStatus } from '../oauth/oauth-status';
import { OAuthControls } from '../oauth/oauth-controls';
import { LanguageSettings } from '../language/language-settings';

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

export function YouTubeSection({ pageInfo }: YouTubeSectionProps) {
  return (
    <div>
      <YouTubeVideoInfo pageInfo={pageInfo} />
      <OAuthStatus />
      <LanguageSettings />
      <AutoTranslateToggle />
      <OAuthControls />
      <TranscriptButton videoId={pageInfo.videoId} />
    </div>
  );
}
