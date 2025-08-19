import React from 'react';
import { Button, ListItem } from 'framework7-react';

interface OAuthControlsProps {
  isAuthenticated: boolean;
  isLoading: boolean;
  onConnect: () => Promise<void>;
  onDisconnect: () => Promise<void>;
}

export function OAuthControls({ 
  isAuthenticated, 
  isLoading, 
  onConnect, 
  onDisconnect 
}: OAuthControlsProps) {
  const handleClick = async () => {
    try {
      if (isAuthenticated) {
        await onDisconnect();
      } else {
        await onConnect();
      }
    } catch (error) {
      console.error('OAuth action failed:', error);
    }
  };

  return (
    <ListItem>
      <Button 
        fill={!isAuthenticated}
        color={isAuthenticated ? "red" : "blue"}
        onClick={handleClick}
        disabled={isLoading}
        style={{ width: '100%' }}
      >
        {isLoading 
          ? 'Connecting...' 
          : isAuthenticated 
            ? 'Disconnect YouTube' 
            : 'Connect YouTube Account'
        }
      </Button>
    </ListItem>
  );
}