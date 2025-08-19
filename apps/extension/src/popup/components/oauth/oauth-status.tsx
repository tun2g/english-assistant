import React from 'react';
import { Block } from 'framework7-react';

interface OAuthStatusProps {
  isAuthenticated: boolean;
}

export function OAuthStatus({ isAuthenticated }: OAuthStatusProps) {
  return (
    <Block strong>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        <span>{isAuthenticated ? '✅' : '🔒'}</span>
        <span style={{ fontSize: '14px', color: isAuthenticated ? '#4caf50' : '#ff9800' }}>
          YouTube API {isAuthenticated ? 'Connected' : 'Not Connected'}
        </span>
      </div>
      {!isAuthenticated && (
        <p style={{ fontSize: '12px', color: '#666', margin: '4px 0' }}>
          Connect to access real YouTube transcripts
        </p>
      )}
    </Block>
  );
}