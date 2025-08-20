// Types for transcript overlay functionality
export interface OverlayPosition {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface OverlaySettings {
  enabled: boolean;
  position: OverlayPosition;
  opacity: number;
  autoHide: boolean;
  hideDelay: number;
}

export interface VideoPlayerInfo {
  currentTime: number;
  duration: number;
  isPlaying: boolean;
  playbackRate: number;
  isSeek?: boolean; // Indicates if this update was caused by seeking
}

export interface SegmentDisplayInfo {
  segment: import('./transcript-types').TranscriptSegment;
  isActive: boolean;
  isVisible: boolean;
  displayText: string;
  translatedText?: string;
}
