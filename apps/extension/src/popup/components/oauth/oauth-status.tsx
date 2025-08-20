import React from 'react';
import { Badge, Alert, AlertDescription } from '@english/ui';
import { useOAuthQuery } from '../../../hooks/use-oauth-query';

export function OAuthStatus() {
  const { isAuthenticated } = useOAuthQuery();

  return (
    <div className="border-b p-4">
      <div className="mb-2 flex items-center gap-2">
        <span>{isAuthenticated ? '✅' : '🔒'}</span>
        <Badge variant={isAuthenticated ? 'success' : 'warning'}>
          YouTube API {isAuthenticated ? 'Connected' : 'Not Connected'}
        </Badge>
      </div>
      {!isAuthenticated && (
        <Alert variant="info">
          <AlertDescription>Connect to access real YouTube transcripts</AlertDescription>
        </Alert>
      )}
    </div>
  );
}
