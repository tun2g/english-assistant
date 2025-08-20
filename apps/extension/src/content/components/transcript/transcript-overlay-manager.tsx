import { useCallback, useEffect, useRef, useState } from 'react';
import { EXTENSION_STORAGE_KEYS } from '../../../shared/constants';
import type {
  LanguageSettingParams,
  TranscriptSegment,
  TranslatedSegment,
  VideoPlayerInfo,
  VideoTranscript,
} from '../../../shared/types/extension-types';

import {
  TranscriptSyncManager,
  type TranscriptSyncCallbacks,
} from '../../features/transcript-sync/transcript-sync-manager';
import {
  TranslationQueueManager,
  type TranslationQueueCallbacks,
} from '../../features/translation/translation-queue-manager';
import { VideoTimeTracker, type VideoTimeTrackerCallbacks } from '../../features/video-tracking/video-time-tracker';

import { LanguageSelector } from '../language/language-selector';
import { DualLanguageDisplay } from './dual-language-display';

interface TranscriptOverlayManagerProps {
  videoId: string;
  transcript: VideoTranscript;
  onReady: () => void;
  onError: (error: Error) => void;
  onSegmentSeek: (seekTime: number) => void;
  showLanguageSelectorInitially?: boolean;
}

export function TranscriptOverlayManager({
  videoId,
  transcript,
  onReady,
  onError,
  onSegmentSeek,
  showLanguageSelectorInitially = false,
}: TranscriptOverlayManagerProps) {
  // State
  const [isActive, setIsActive] = useState(false);
  const [showDisplay, setShowDisplay] = useState(false);
  const [showLanguageSelector, setShowLanguageSelector] = useState(false);
  const [currentSegment, setCurrentSegment] = useState<TranscriptSegment | null>(null);
  const [currentSegmentIndex, setCurrentSegmentIndex] = useState(-1);
  const [translatedSegments, setTranslatedSegments] = useState<Map<number, TranslatedSegment>>(new Map());

  // Keep refs in sync with state
  useEffect(() => {
    translatedSegmentsRef.current = translatedSegments;
  }, [translatedSegments]);

  const [languageSettingParams, setLanguageSettingParams] = useState<LanguageSettingParams>({
    primaryLanguage: 'en',
    secondaryLanguage: 'vi',
    dualLanguageEnabled: true,
    autoTranslateEnabled: true,
  });

  useEffect(() => {
    languageSettingParamsRef.current = languageSettingParams;

    // Update translation context when language settings change
    if (translationQueueManagerRef.current) {
      translationQueueManagerRef.current.setContext(videoId, languageSettingParams.secondaryLanguage);
    }
  }, [languageSettingParams, videoId]);

  // Refs for managers (avoid re-creating on every render)
  const videoTimeTrackerRef = useRef<VideoTimeTracker | null>(null);
  const transcriptSyncManagerRef = useRef<TranscriptSyncManager | null>(null);
  const translationQueueManagerRef = useRef<TranslationQueueManager | null>(null);

  // Ref for translated segments to avoid dependency loops
  const translatedSegmentsRef = useRef<Map<number, TranslatedSegment>>(new Map());

  // Ref for language settings to avoid dependency loops
  const languageSettingParamsRef = useRef<LanguageSettingParams>(languageSettingParams);

  // Load language settings from storage
  const loadLanguageSettingParams = useCallback(async () => {
    try {
      const result = await new Promise<any>(resolve => {
        if (!chrome?.storage?.local) {
          resolve({});
          return;
        }

        chrome.storage.local.get(
          [
            EXTENSION_STORAGE_KEYS.DUAL_LANGUAGE_ENABLED,
            EXTENSION_STORAGE_KEYS.PRIMARY_LANGUAGE,
            EXTENSION_STORAGE_KEYS.SECONDARY_LANGUAGE,
            EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED,
          ],
          result => {
            if (chrome.runtime.lastError) {
              resolve({});
              return;
            }
            resolve(result);
          }
        );
      });

      const newSettings: LanguageSettingParams = {
        dualLanguageEnabled: result[EXTENSION_STORAGE_KEYS.DUAL_LANGUAGE_ENABLED] ?? true,
        primaryLanguage: result[EXTENSION_STORAGE_KEYS.PRIMARY_LANGUAGE] ?? 'en',
        secondaryLanguage: result[EXTENSION_STORAGE_KEYS.SECONDARY_LANGUAGE] ?? 'es',
        autoTranslateEnabled: result[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED] ?? true,
      };

      setLanguageSettingParams(newSettings);
    } catch (error) {
      // Use default settings if storage fails
      const defaultSettings: LanguageSettingParams = {
        dualLanguageEnabled: true,
        primaryLanguage: 'en',
        secondaryLanguage: 'vi',
        autoTranslateEnabled: true,
      };
      setLanguageSettingParams(defaultSettings);
    }
  }, []);

  // Save language settings to storage
  const saveLanguageSettingParams = useCallback(async (settings: LanguageSettingParams) => {
    await new Promise<void>(resolve => {
      chrome.storage.local.set(
        {
          [EXTENSION_STORAGE_KEYS.DUAL_LANGUAGE_ENABLED]: settings.dualLanguageEnabled,
          [EXTENSION_STORAGE_KEYS.PRIMARY_LANGUAGE]: settings.primaryLanguage,
          [EXTENSION_STORAGE_KEYS.SECONDARY_LANGUAGE]: settings.secondaryLanguage,
          [EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]: settings.autoTranslateEnabled,
        },
        resolve
      );
    });
  }, []);

  // Initialize managers
  const initializeManagers = useCallback(() => {
    // Video time tracker callbacks
    const videoTimeCallbacks: VideoTimeTrackerCallbacks = {
      onTimeUpdate: (playerInfo: VideoPlayerInfo) => {
        transcriptSyncManagerRef.current?.updateVideoTime(playerInfo);

        // Auto-translate logic is handled in the language change handler
        // We don't need to queue individual segments during time updates
      },
    };

    // Transcript sync callbacks
    const syncCallbacks: TranscriptSyncCallbacks = {
      onActiveSegmentChange: (segment, index) => {
        setCurrentSegment(segment);
        setCurrentSegmentIndex(index);

        // Translation is handled in batch, so individual segment translation is not needed here
      },
      onSegmentEnter: (_segment, index) => {},
      onSegmentExit: (_segment, index) => {},
    };

    // Translation callbacks
    const translationCallbacks: TranslationQueueCallbacks = {
      onTranslationComplete: (translatedSegment, segmentIndex) => {
        setTranslatedSegments(prev => new Map(prev).set(segmentIndex, translatedSegment));
      },
      onTranslationError: (segmentIndex, error) => {},
      onBatchTranslationComplete: (translatedSegments: TranslatedSegment[]) => {
        setTranslatedSegments(prev => {
          const newMap = new Map(prev);
          translatedSegments.forEach((segment, index) => {
            const segmentIndex = segment.index ?? index;
            newMap.set(segmentIndex, segment);
          });
          return newMap;
        });
      },
    };

    // Create managers
    videoTimeTrackerRef.current = new VideoTimeTracker(videoTimeCallbacks);
    transcriptSyncManagerRef.current = new TranscriptSyncManager(syncCallbacks);
    translationQueueManagerRef.current = new TranslationQueueManager(translationCallbacks);

    // Configure managers
    transcriptSyncManagerRef.current.setTranscript(transcript);
    translationQueueManagerRef.current.setContext(videoId, languageSettingParamsRef.current.secondaryLanguage);
  }, [videoId, transcript]); // Removed languageSettingParams dependency

  // Start overlay
  const start = useCallback(async () => {
    if (isActive) return;

    setIsActive(true);

    // Start all managers with debug logging

    videoTimeTrackerRef.current?.start();

    transcriptSyncManagerRef.current?.start();

    translationQueueManagerRef.current?.start();

    // Small delay to ensure managers are fully initialized
    await new Promise(resolve => setTimeout(resolve, 100));
  }, [isActive]);

  // Stop overlay
  const stop = useCallback(() => {
    if (!isActive) return;

    setIsActive(false);

    // Stop all managers
    videoTimeTrackerRef.current?.stop();
    transcriptSyncManagerRef.current?.stop();
    translationQueueManagerRef.current?.stop();

    // Hide displays
    setShowDisplay(false);
    setShowLanguageSelector(false);
  }, [isActive]);

  // Handle segment click (for seeking)
  const handleSegmentClick = useCallback(
    (segmentIndex: number) => {
      const seekTime = transcriptSyncManagerRef.current?.seekToSegment(segmentIndex);
      if (seekTime && seekTime >= 0) {
        onSegmentSeek(seekTime);
      }
    },
    [onSegmentSeek]
  );

  // Auto-start overlay with existing settings (without showing language selector)
  const autoStartOverlayWithSettings = useCallback(
    async (settings: LanguageSettingParams) => {
      // Update translation context
      translationQueueManagerRef.current?.setContext(videoId, settings.secondaryLanguage);

      // Start the overlay and all managers if dual language is enabled
      if (settings.dualLanguageEnabled) {
        // Start all managers
        await start();

        // Show the dual language display immediately
        setShowDisplay(true);

        // If auto-translate is enabled, start queuing segments for translation
        if (settings.autoTranslateEnabled) {
          // Give managers a moment to fully initialize, then queue segments
          setTimeout(() => {
            const videoElement = document.querySelector('video') as HTMLVideoElement;
            if (videoElement) {
              const currentTime = videoElement.currentTime;

              // Force an initial time update to get current segment
              if (videoTimeTrackerRef.current) {
                videoTimeTrackerRef.current.forceTimeUpdate();
              }

              // Translate all transcript segments at once
              translationQueueManagerRef.current
                ?.translateAllSegments(transcript.segments)
                .then(translatedSegments => {})
                .catch(error => {});
            }
          }, 100);
        }
      } else {
        setShowDisplay(false);
      }
    },
    [videoId, start]
  );

  // Handle language settings change
  const handleLanguageChange = useCallback(
    async (newSettings: LanguageSettingParams) => {
      console.log('TranscriptOverlayManager: handleLanguageChange called with:', newSettings);

      try {
        // Save settings first
        setLanguageSettingParams(newSettings);
        await saveLanguageSettingParams(newSettings);
        console.log('TranscriptOverlayManager: Settings saved successfully');

        // Update translation context
        translationQueueManagerRef.current?.setContext(videoId, newSettings.secondaryLanguage);

        // Close language selector
        setShowLanguageSelector(false);
        console.log('TranscriptOverlayManager: Language selector closed');

        // Start the overlay and all managers after saving settings
        if (newSettings.dualLanguageEnabled) {
          console.log('TranscriptOverlayManager: Dual language enabled, starting overlay...');

          // Start all managers
          await start();
          console.log('TranscriptOverlayManager: Managers started');

          // Show the dual language display
          setShowDisplay(true);
          console.log('TranscriptOverlayManager: Display set to show');

          // If auto-translate is enabled, start queuing segments for translation
          if (newSettings.autoTranslateEnabled) {
            // Give managers a moment to fully initialize, then queue segments
            setTimeout(() => {
              const videoElement = document.querySelector('video') as HTMLVideoElement;
              if (videoElement) {
                const currentTime = videoElement.currentTime;

                // Force an initial time update to get current segment
                if (videoTimeTrackerRef.current) {
                  videoTimeTrackerRef.current.forceTimeUpdate();
                }

                // Translate all transcript segments at once
                translationQueueManagerRef.current
                  ?.translateAllSegments(transcript.segments)
                  .then(translatedSegments => {})
                  .catch(error => {});
              }
            }, 100);
          }
        } else {
          // Single language mode
          setShowDisplay(false);
        }
      } catch (error) {
        // Show error but don't break the flow
        setShowLanguageSelector(false);
      }
    },
    [videoId, saveLanguageSettingParams, start]
  );

  // Initialize on mount - run only once per component instance
  useEffect(() => {
    const init = async () => {
      try {
        await loadLanguageSettingParams();

        onReady();
      } catch (error) {
        onError(error as Error);
      }
    };

    init();
  }, [videoId]); // Only depend on videoId to avoid infinite loops

  // Initialize managers once when component mounts
  useEffect(() => {
    initializeManagers();
  }, [initializeManagers]); // Only run when initializeManagers changes

  // Handle language settings and auto-start separately
  useEffect(() => {
    if (languageSettingParams) {
      // Show language selector initially if requested
      if (showLanguageSelectorInitially) {
        setShowLanguageSelector(true);

        // Also auto-start if dual language is already enabled in settings
        if (languageSettingParams.dualLanguageEnabled) {
          autoStartOverlayWithSettings(languageSettingParams);
        }
      }
    }
  }, [languageSettingParams, showLanguageSelectorInitially]); // Removed autoStartOverlayWithSettings dependency to prevent loops

  // Don't auto-start overlay - let user click extension button to show language selector
  // This prevents interference with the language selector popup workflow

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      stop();
    };
  }, [stop]);

  // Get current translated segment
  const currentTranslatedSegment =
    currentSegmentIndex >= 0 ? translatedSegments.get(currentSegmentIndex) || null : null;

  return (
    <>
      {/* Dual Language Display - Only show if active and settings allow */}
      {showDisplay && isActive && languageSettingParams.dualLanguageEnabled && (
        <DualLanguageDisplay
          currentSegment={currentSegment}
          translatedSegment={currentTranslatedSegment}
          languageSettings={languageSettingParams}
          segmentIndex={currentSegmentIndex}
          onSegmentClick={handleSegmentClick}
          onSettingsClick={() => setShowLanguageSelector(true)}
          onClose={() => {
            setShowDisplay(false);
            stop(); // Stop all managers when closing display
          }}
        />
      )}

      {/* Language Selector Modal */}
      {showLanguageSelector && (
        <LanguageSelector
          initialSettings={languageSettingParams}
          onSave={handleLanguageChange}
          onClose={() => setShowLanguageSelector(false)}
        />
      )}
    </>
  );
}
