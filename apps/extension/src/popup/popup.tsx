import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '../lib/query-client';
import { Card, CardContent, CardHeader, CardTitle, Alert, AlertDescription } from '@english/ui';
import { PopupErrorBoundary } from './components/error-boundary';

import { usePageInfoQuery } from '../hooks/use-page-info-query';
import { useOAuthQuery } from '../hooks/use-oauth-query';
import { useAutoTranslateQuery } from '../hooks/use-auto-translate-query';

import { YouTubeSection } from './components/youtube/youtube-section';
import { QuickActions } from './components/navigation/quick-actions';
import { SettingsButton } from './components/navigation/settings-button';
import '../styles/globals.scss';
import './popup.scss';

function PopupApp() {
  const { pageInfo, isLoading: isPageLoading, error: pageError } = usePageInfoQuery();
  const {
    isAuthenticated: isOAuthAuthenticated,
    isLoading: isOAuthLoading,
    connect: connectOAuth,
    disconnect: disconnectOAuth,
    error: oauthError,
  } = useOAuthQuery();
  const {
    isEnabled: isAutoTranslateEnabled,
    toggle: toggleAutoTranslate,
    error: autoTranslateError,
  } = useAutoTranslateQuery();

  const handleGetTranscript = async (): Promise<void> => {
    if (!isOAuthAuthenticated) {
      await connectOAuth();
      return;
    }

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

  if (isPageLoading) {
    return (
      <div className="popup-container w-80 p-4">
        <Card>
          <CardHeader>
            <CardTitle>English Learning Assistant</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-center py-8">
              <p className="text-sm text-gray-600">Loading...</p>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="popup-container w-80 p-4">
      <Card>
        <CardHeader>
          <CardTitle>English Learning Assistant</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {pageInfo?.isYouTube ? (
            <YouTubeSection
              pageInfo={pageInfo}
              isOAuthAuthenticated={isOAuthAuthenticated}
              isOAuthLoading={isOAuthLoading}
              isAutoTranslateEnabled={isAutoTranslateEnabled}
              onOAuthConnect={async () => {
                await connectOAuth();
              }}
              onOAuthDisconnect={async () => {
                await disconnectOAuth();
              }}
              onAutoTranslateToggle={async (enabled: boolean) => {
                await toggleAutoTranslate(enabled);
              }}
              onGetTranscript={handleGetTranscript}
            />
          ) : (
            <QuickActions />
          )}

          <div className="border-t p-4">
            <SettingsButton />
            {(pageError || oauthError || autoTranslateError) && (
              <Alert variant="destructive" className="mt-2">
                <AlertDescription className="text-xs">
                  {String(pageError || oauthError || autoTranslateError)}
                </AlertDescription>
              </Alert>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// Enhanced popup initialization
function initializePopup() {
  const container = document.getElementById('english-learning-extension-popup');
  if (!container) {
    throw new Error('English Learning Extension popup container not found in popup HTML');
  }

  const root = ReactDOM.createRoot(container);

  root.render(
    <PopupErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <PopupApp />
      </QueryClientProvider>
    </PopupErrorBoundary>
  );
}

// Initialize when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initializePopup);
} else {
  initializePopup();
}
