import { TRANSCRIPT_CONFIG } from '@/shared/constants';
import { TranscriptSegment, VideoPlayerInfo, VideoTranscript } from '@/shared/types';

export interface TranscriptSyncCallbacks {
  onActiveSegmentChange: (segment: TranscriptSegment | null, index: number) => void;
  onSegmentEnter: (segment: TranscriptSegment, index: number) => void;
  onSegmentExit: (segment: TranscriptSegment, index: number) => void;
}

export class TranscriptSyncManager {
  private callbacks: TranscriptSyncCallbacks;
  private transcript: VideoTranscript | null = null;
  private currentSegmentIndex: number = -1;
  private isActive: boolean = false;

  constructor(callbacks: TranscriptSyncCallbacks) {
    this.callbacks = callbacks;
  }

  // Set the transcript to sync with
  setTranscript(transcript: VideoTranscript | null): void {
    this.transcript = transcript;
    this.currentSegmentIndex = -1;
  }

  // Start synchronization
  start(): void {
    if (this.isActive) return;

    this.isActive = true;
  }

  // Stop synchronization
  stop(): void {
    if (!this.isActive) return;

    this.isActive = false;
    this.currentSegmentIndex = -1;
    this.callbacks.onActiveSegmentChange(null, -1);
  }

  // Update with current video time
  updateVideoTime(playerInfo: VideoPlayerInfo): void {
    if (!this.isActive || !this.transcript || !this.transcript.segments.length) {
      return;
    }

    const currentTime = playerInfo.currentTime;
    const newSegmentIndex = this.findActiveSegmentIndex(currentTime);

    // Force segment change detection on seek to resync immediately
    const forceUpdate = playerInfo.isSeek || false;

    // Check if we've moved to a different segment or force update due to seek
    if (newSegmentIndex !== this.currentSegmentIndex || forceUpdate) {
      // Exit previous segment if there was one
      if (this.currentSegmentIndex >= 0 && this.currentSegmentIndex < this.transcript.segments.length) {
        const prevSegment = this.transcript.segments[this.currentSegmentIndex];
        this.callbacks.onSegmentExit(prevSegment, this.currentSegmentIndex);
      }

      // Update current segment
      this.currentSegmentIndex = newSegmentIndex;

      // Enter new segment if there is one
      if (this.currentSegmentIndex >= 0) {
        const currentSegment = this.transcript.segments[this.currentSegmentIndex];
        this.callbacks.onSegmentEnter(currentSegment, this.currentSegmentIndex);
        this.callbacks.onActiveSegmentChange(currentSegment, this.currentSegmentIndex);
      } else {
        this.callbacks.onActiveSegmentChange(null, -1);
      }
    }
  }

  // Find the active segment index for given time
  private findActiveSegmentIndex(currentTime: number): number {
    if (!this.transcript || !this.transcript.segments.length) return -1;

    // Add offset to highlight segment slightly before it speaks
    const timeWithOffset = currentTime + TRANSCRIPT_CONFIG.SEGMENT_HIGHLIGHT_OFFSET / 1000;

    // If using estimated timing (when all segments have start/end times), use proportional matching
    const hasRealTiming = this.transcript.segments.some(s => s.start > 0 && s.end > s.start + 0.5);

    if (!hasRealTiming) {
      // Use proportional matching for estimated timing
      const totalDuration = this.transcript.segments[this.transcript.segments.length - 1]?.end || 0;
      if (totalDuration > 0) {
        const videoElement = document.querySelector('video') as HTMLVideoElement;
        const videoDuration = videoElement?.duration || totalDuration;

        // Calculate proportional time in our estimated timeline
        const proportionalTime = (timeWithOffset / videoDuration) * totalDuration;

        // Find segment using proportional time
        for (let i = 0; i < this.transcript.segments.length; i++) {
          const segment = this.transcript.segments[i];
          if (proportionalTime >= segment.start && proportionalTime <= segment.end) {
            return i;
          }
        }
      }
    }

    // Binary search for efficiency with large transcripts (for real timing)
    let left = 0;
    let right = this.transcript.segments.length - 1;

    while (left <= right) {
      const mid = Math.floor((left + right) / 2);
      const segment = this.transcript.segments[mid];

      if (timeWithOffset >= segment.start && timeWithOffset <= segment.end) {
        return mid;
      } else if (timeWithOffset < segment.start) {
        right = mid - 1;
      } else {
        left = mid + 1;
      }
    }

    return -1;
  }

  // Get segments around current time (for preloading/translation)
  getSegmentsAround(currentTime: number, range: number = 5): TranscriptSegment[] {
    if (!this.transcript || !this.transcript.segments.length) return [];

    const currentIndex = this.findActiveSegmentIndex(currentTime);
    if (currentIndex === -1) return [];

    const start = Math.max(0, currentIndex - range);
    const end = Math.min(this.transcript.segments.length - 1, currentIndex + range);

    return this.transcript.segments.slice(start, end + 1);
  }

  // Get upcoming segments for translation queue
  getUpcomingSegments(currentTime: number, count: number = 10): TranscriptSegment[] {
    if (!this.transcript || !this.transcript.segments.length) return [];

    const segments = this.transcript.segments.filter(segment => segment.start > currentTime);

    return segments.slice(0, count);
  }

  // Get current segment
  getCurrentSegment(): TranscriptSegment | null {
    if (
      !this.transcript ||
      this.currentSegmentIndex < 0 ||
      this.currentSegmentIndex >= this.transcript.segments.length
    ) {
      return null;
    }

    return this.transcript.segments[this.currentSegmentIndex];
  }

  // Get current segment index
  getCurrentSegmentIndex(): number {
    return this.currentSegmentIndex;
  }

  // Get total segments count
  getTotalSegments(): number {
    return this.transcript?.segments.length || 0;
  }

  // Seek to specific segment
  seekToSegment(index: number): number {
    if (!this.transcript || index < 0 || index >= this.transcript.segments.length) {
      return -1;
    }

    const segment = this.transcript.segments[index];
    return segment.start;
  }

  // Get segment at specific index
  getSegment(index: number): TranscriptSegment | null {
    if (!this.transcript || index < 0 || index >= this.transcript.segments.length) {
      return null;
    }

    return this.transcript.segments[index];
  }
}
