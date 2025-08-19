import React from 'react';
import { App, View, Page, Navbar, Block } from 'framework7-react';
import 'framework7/css/bundle';

import { usePageInfo } from '../hooks/use-page-info';
import { useOAuth } from '../hooks/use-oauth';
import { useAutoTranslate } from '../hooks/use-auto-translate';

import { YouTubeSection } from './youtube/youtube-section';
import { QuickActions } from './navigation/quick-actions';
import { SettingsButton } from './navigation/settings-button';

export function PopupApp() {
  const { pageInfo, isLoading: isPageLoading } = usePageInfo();
  const { isAuthenticated: isOAuthAuthenticated, isLoading: isOAuthLoading, connect: connectOAuth, disconnect: disconnectOAuth } = useOAuth();
  const { isEnabled: isAutoTranslateEnabled, toggle: toggleAutoTranslate } = useAutoTranslate();

  const f7params = {
    name: 'English Learning Assistant',
    theme: 'auto',
  };

  const handleGetTranscript = async (): Promise<void> => {
    if (!isOAuthAuthenticated) {
      await connectOAuth();
      return;
    }

    // Send message to content script to trigger transcript request
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    if (tab.id) {
      await chrome.tabs.sendMessage(tab.id, {
        action: 'GET_TRANSCRIPT_WITH_AUTH'
      });
    }
  };

  if (isPageLoading) {
    return (
      <App {...f7params}>
        <View main className="safe-areas" url="/">
          <Page>
            <Navbar title="English Learning Assistant" />
            <Block className="text-align-center">
              <p>Loading...</p>
            </Block>
          </Page>
        </View>
      </App>
    );
  }

  return (
    <App {...f7params}>
      <View main className="safe-areas" url="/">
        <Page>
          <Navbar title="English Learning Assistant" />
          
          {pageInfo?.isYouTube ? (
            <YouTubeSection
              pageInfo={pageInfo}
              isOAuthAuthenticated={isOAuthAuthenticated}
              isOAuthLoading={isOAuthLoading}
              isAutoTranslateEnabled={isAutoTranslateEnabled}
              onOAuthConnect={connectOAuth}
              onOAuthDisconnect={disconnectOAuth}
              onAutoTranslateToggle={toggleAutoTranslate}
              onGetTranscript={handleGetTranscript}
            />
          ) : (
            <QuickActions />
          )}
          
          <SettingsButton />
        </Page>
      </View>
    </App>
  );
}