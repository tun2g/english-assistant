import { extractVideoIdFromCurrentUrl } from '../../utils/video/video-utils';

export interface VideoMonitorCallbacks {
  onVideoChange: (videoId: string | null) => Promise<void>;
}

export class VideoMonitor {
  private callbacks: VideoMonitorCallbacks;
  private observer: MutationObserver | null = null;
  private isMonitoring: boolean = false;
  private lastUrl: string = '';

  constructor(callbacks: VideoMonitorCallbacks) {
    this.callbacks = callbacks;
  }

  // Start monitoring for video changes
  start(): void {
    if (this.isMonitoring) {
      console.log('VideoMonitor: Already monitoring');
      return;
    }

    console.log('VideoMonitor: Starting video monitoring...');
    this.isMonitoring = true;
    this.lastUrl = location.href;

    this.setupUrlChangeMonitoring();
    this.setupPopstateListener();
  }

  // Stop monitoring
  stop(): void {
    if (!this.isMonitoring) return;

    console.log('VideoMonitor: Stopping video monitoring...');
    this.isMonitoring = false;

    if (this.observer) {
      this.observer.disconnect();
      this.observer = null;
    }

    // Remove popstate listener
    window.removeEventListener('popstate', this.handlePopstate);
  }

  // Setup URL change monitoring using MutationObserver
  private setupUrlChangeMonitoring(): void {
    // Use a much more specific observer that only watches for actual navigation changes
    this.observer = new MutationObserver(() => {
      const currentUrl = location.href;
      if (currentUrl !== this.lastUrl) {
        console.log('VideoMonitor: URL changed from', this.lastUrl, 'to', currentUrl);
        this.lastUrl = currentUrl;
        this.handleUrlChange();
      }
    });

    // Only observe the main app element instead of entire body
    const appElement = document.querySelector('#content, #page-manager, ytd-app') || document.body;
    this.observer.observe(appElement, {
      childList: true,
      subtree: false, // Don't watch all children, just direct children
    });
  }

  // Setup popstate event listener
  private setupPopstateListener(): void {
    window.addEventListener('popstate', this.handlePopstate);
  }

  // Handle popstate events (back/forward navigation)
  private handlePopstate = (): void => {
    this.handleUrlChange();
  };

  // Handle URL changes
  private handleUrlChange(): void {
    setTimeout(() => this.checkCurrentVideo(), 500);
  }

  // Check current video and notify of changes
  private async checkCurrentVideo(): Promise<void> {
    const videoId = extractVideoIdFromCurrentUrl();
    console.log('VideoMonitor: Checking video ID:', videoId);

    await this.callbacks.onVideoChange(videoId);
  }
}
