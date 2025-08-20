import React from 'react';
import { Button, Alert, AlertDescription, Spinner } from '@english/ui';
import { useOAuthQuery } from '../../../hooks';

export function OAuthControls() {
  const { isAuthenticated, isLoading, connect, disconnect, error } = useOAuthQuery();

  const handleClick = async () => {
    try {
      if (isAuthenticated) {
        await disconnect();
      } else {
        await connect();
      }
    } catch (error) {
      console.error('OAuth action failed:', error);
    }
  };

  return (
    <div className="p-4">
      <Button
        variant={isAuthenticated ? 'destructive' : 'default'}
        onClick={handleClick}
        disabled={isLoading}
        className="w-full"
      >
        {isLoading && <Spinner size="sm" className="mr-2" />}
        {isLoading ? 'Connecting...' : isAuthenticated ? 'Disconnect YouTube' : 'Connect YouTube Account'}
      </Button>
      {error && (
        <Alert variant="destructive" className="mt-2">
          <AlertDescription>Failed to {isAuthenticated ? 'disconnect' : 'connect'}. Please try again.</AlertDescription>
        </Alert>
      )}
    </div>
  );
}
