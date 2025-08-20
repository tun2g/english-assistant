import React from 'react';
import { Switch, Alert, AlertDescription } from '@english/ui';
import { useAutoTranslateQuery } from '../../../hooks';

export function AutoTranslateToggle() {
  const { isEnabled, toggle, isToggling, error } = useAutoTranslateQuery();

  const handleToggle = async (checked: boolean) => {
    try {
      await toggle(checked);
    } catch (error) {
      console.error('Failed to toggle auto-translate:', error);
    }
  };

  return (
    <div className="border-b p-4">
      <div className="flex items-center justify-between">
        <div className="flex flex-col">
          <span className="text-sm font-medium">Auto-translate Videos</span>
          <span className="text-xs text-gray-500">Automatically translate video content</span>
        </div>
        <Switch checked={isEnabled} onCheckedChange={handleToggle} disabled={isToggling} />
      </div>
      {error && (
        <Alert variant="destructive" className="mt-2">
          <AlertDescription>Failed to update setting. Please try again.</AlertDescription>
        </Alert>
      )}
    </div>
  );
}
