import React from 'react';
import { ListItem, Toggle } from 'framework7-react';

interface AutoTranslateToggleProps {
  isEnabled: boolean;
  onToggle: (enabled: boolean) => Promise<void>;
}

export function AutoTranslateToggle({ isEnabled, onToggle }: AutoTranslateToggleProps) {
  const handleToggle = async (e: React.ChangeEvent<HTMLInputElement>) => {
    try {
      await onToggle(e.target.checked);
    } catch (error) {
      console.error('Failed to toggle auto-translate:', error);
    }
  };

  return (
    <ListItem title="Auto-translate Videos">
      <Toggle
        slot="after"
        checked={isEnabled}
        onChange={handleToggle}
      />
    </ListItem>
  );
}