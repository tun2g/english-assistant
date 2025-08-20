import React from 'react';
import { Button } from '@english/ui';

export function SettingsButton() {
  const handleSettingsClick = () => {
    chrome.runtime.openOptionsPage();
  };

  return (
    <Button className="w-full" onClick={handleSettingsClick}>
      Settings
    </Button>
  );
}
