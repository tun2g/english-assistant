import { YOUTUBE_SELECTORS } from '../../../shared/constants';
import type { VideoPlayerInfo } from '../../../shared/types/extension-types';
import { captureError } from '../../utils/error-handler';

export interface VideoTimeTrackerCallbacks {
  onTimeUpdate: (playerInfo: VideoPlayerInfo) => void;
}

export class VideoTimeTracker {
  private callbacks: VideoTimeTrackerCallbacks;
  private videoElement: HTMLVideoElement | null = null;
  private isTracking: boolean = false;
  private rafId: number | null = null;
  private lastCurrentTime: number = -1;

  constructor(callbacks: VideoTimeTrackerCallbacks) {
    this.callbacks = callbacks;
  }

  // Start tracking video time
  start(): void {
    if (this.isTracking) {
      return;
    }

    this.findVideoElement();
    if (!this.videoElement) {
      setTimeout(() => this.start(), 1000);
      return;
    }

    this.isTracking = true;
    this.startTimeTracking();
  }

  // Stop tracking video time
  stop(): void {
    if (!this.isTracking) return;

    this.isTracking = false;

    if (this.rafId) {
      cancelAnimationFrame(this.rafId);
      this.rafId = null;
    }

    this.videoElement = null;
    this.lastCurrentTime = -1;
  }

  // Find the YouTube video element
  private findVideoElement(): void {
    // Try different selectors to find the video element
    const selectors = [
      `${YOUTUBE_SELECTORS.MOVIE_PLAYER} ${YOUTUBE_SELECTORS.VIDEO_ELEMENT}`,
      `${YOUTUBE_SELECTORS.VIDEO_PLAYER} ${YOUTUBE_SELECTORS.VIDEO_ELEMENT}`,
      YOUTUBE_SELECTORS.VIDEO_ELEMENT,
    ];

    for (const selector of selectors) {
      const video = document.querySelector(selector) as HTMLVideoElement;
      if (video && video.tagName === 'VIDEO') {
        this.videoElement = video;
        break;
      }
    }
  }

  // Start the time tracking loop
  private startTimeTracking(): void {
    const trackTime = () => {
      if (!this.isTracking || !this.videoElement) return;

      const currentTime = this.videoElement.currentTime;

      // Debug: Log video state periodically
      if (Math.floor(currentTime) % 5 === 0 && Math.abs(currentTime - this.lastCurrentTime) > 0.2) {
        console.log('VideoTimeTracker: Video state -', {
          currentTime: currentTime.toFixed(2),
          duration: (this.videoElement.duration || 0).toFixed(2),
          isPlaying: !this.videoElement.paused,
          playbackRate: this.videoElement.playbackRate,
        });
      }

      // Detect significant time jumps (seeking) or regular time updates
      const timeDiff = Math.abs(currentTime - this.lastCurrentTime);
      const isSeek = timeDiff > 2.0; // Consider jumps > 2 seconds as seeks
      const isRegularUpdate = timeDiff > 0.1; // Regular time progression

      if (isSeek || isRegularUpdate) {
        const playerInfo: VideoPlayerInfo = {
          currentTime,
          duration: this.videoElement.duration || 0,
          isPlaying: !this.videoElement.paused,
          playbackRate: this.videoElement.playbackRate || 1,
          isSeek: isSeek,
        };

        if (isSeek) {
          console.log('VideoTimeTracker: Seek detected', {
            from: this.lastCurrentTime.toFixed(2),
            to: currentTime.toFixed(2),
            jump: timeDiff.toFixed(2),
          });
        }

        try {
          this.callbacks.onTimeUpdate(playerInfo);
          this.lastCurrentTime = currentTime;
        } catch (error) {
          captureError('Error in time update callback', error, 'callback', 'VideoTimeTracker');
        }
      }

      // Schedule next check
      this.rafId = requestAnimationFrame(trackTime);
    };

    // Start the tracking loop
    this.rafId = requestAnimationFrame(trackTime);
  }

  // Get current video information (synchronous)
  getCurrentVideoInfo(): VideoPlayerInfo | null {
    if (!this.videoElement) return null;

    return {
      currentTime: this.videoElement.currentTime,
      duration: this.videoElement.duration || 0,
      isPlaying: !this.videoElement.paused,
      playbackRate: this.videoElement.playbackRate || 1,
    };
  }

  // Seek to specific time in video
  seekTo(time: number): void {
    if (!this.videoElement) return;

    this.videoElement.currentTime = Math.max(0, Math.min(time, this.videoElement.duration));
  }

  // Force an immediate time update (useful for initialization)
  forceTimeUpdate(): void {
    if (!this.isTracking || !this.videoElement) return;

    const currentTime = this.videoElement.currentTime;
    const playerInfo: VideoPlayerInfo = {
      currentTime,
      duration: this.videoElement.duration || 0,
      isPlaying: !this.videoElement.paused,
      playbackRate: this.videoElement.playbackRate || 1,
    };

    try {
      this.callbacks.onTimeUpdate(playerInfo);
      this.lastCurrentTime = currentTime;
    } catch (error) {
      captureError('Error in forced time update callback', error, 'force-callback', 'VideoTimeTracker');
    }
  }

  // Check if video element is available
  get hasVideoElement(): boolean {
    return this.videoElement !== null;
  }
}
