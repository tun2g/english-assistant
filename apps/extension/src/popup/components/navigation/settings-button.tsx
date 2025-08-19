import React from 'react';
import { Block, Button } from 'framework7-react';

export function SettingsButton() {
  const handleSettingsClick = () => {
    chrome.runtime.openOptionsPage();
  };

  return (
    <Block>
      <Button fill onClick={handleSettingsClick}>
        Settings
      </Button>
    </Block>
  );
}