// Types for extension messaging
export interface ExtensionMessage {
  action: string;
  data?: any;
  enabled?: boolean;
  url?: string;
}

export interface MessageResponse {
  success: boolean;
  data?: any;
  message?: string;
  error?: string;
}

export interface LanguageUpdateMessage extends ExtensionMessage {
  action: 'UPDATE_LANGUAGE_SETTINGS';
  data: {
    primaryLanguage: string;
    secondaryLanguage: string;
    dualLanguageEnabled: boolean;
  };
}

export interface TranscriptOverlayMessage extends ExtensionMessage {
  action: 'TOGGLE_TRANSCRIPT_OVERLAY';
  data: {
    enabled: boolean;
    videoId?: string;
  };
}
