import { Button, Card, CardContent, CardHeader, CardTitle, Badge } from '@english/ui';
import type { TranscriptSegment, TranslatedSegment, LanguageSettings } from '../../../shared/types/extension-types';

interface DualLanguageDisplayProps {
  currentSegment: TranscriptSegment | null;
  translatedSegment: TranslatedSegment | null;
  languageSettings: LanguageSettings;
  segmentIndex: number;
  onSegmentClick: (segmentIndex: number) => void;
  onSettingsClick: () => void;
  onClose: () => void;
}

export function DualLanguageDisplay({
  currentSegment,
  translatedSegment,
  languageSettings,
  segmentIndex,
  onSegmentClick,
  onSettingsClick,
  onClose,
}: DualLanguageDisplayProps) {
  const handleSegmentClick = () => {
    if (segmentIndex >= 0) {
      onSegmentClick(segmentIndex);
    }
  };

  // Remove createPortal since we're now rendered within Shadow DOM
  return (
    <div
      className="z-extension-overlay animate-slide-up fixed bottom-[100px] right-5 max-h-[350px] w-[420px]"
      data-english-extension="dual-language-display"
    >
      <Card
        className="overflow-hidden rounded-xl border border-gray-200 bg-white text-gray-900 shadow-2xl"
        data-english-extension="dual-language-card"
      >
        <CardHeader
          className="flex flex-row items-center justify-between space-y-0 border-b border-gray-200 bg-gray-50 px-5 pb-3 pt-4"
          data-english-extension="dual-language-header"
        >
          <div className="flex items-center gap-2">
            <div className="bg-primary h-2 w-2 animate-pulse rounded-full"></div>
            <CardTitle className="text-base font-semibold text-gray-900" data-english-extension="dual-language-title">
              🎯 Live Transcript
            </CardTitle>
          </div>
          <div className="flex gap-2" data-english-extension="dual-language-controls">
            <Button
              variant="ghost"
              size="sm"
              onClick={onSettingsClick}
              className="flex h-7 w-7 items-center justify-center rounded-md p-0 text-gray-500 transition-colors hover:bg-blue-50 hover:text-blue-600"
              title="Language Settings"
              data-english-extension="settings-button"
            >
              ⚙️
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={onClose}
              className="flex h-7 w-7 items-center justify-center rounded-md p-0 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600"
              title="Close Transcript"
              data-english-extension="close-button"
            >
              ✕
            </Button>
          </div>
        </CardHeader>
        <CardContent
          className="custom-scrollbar max-h-[220px] min-h-[80px] overflow-y-auto px-5 pb-4 pt-3"
          data-english-extension="dual-language-content"
        >
          {currentSegment ? (
            <div className="flex flex-col gap-4" data-english-extension="transcript-segments">
              {/* Primary language (original) */}
              <div
                className="group cursor-pointer rounded-lg border border-blue-200 bg-blue-50 p-4 transition-all hover:border-blue-300 hover:bg-blue-100 hover:shadow-md"
                onClick={handleSegmentClick}
                data-english-extension="original-segment"
              >
                <div className="mb-2 flex items-center justify-between">
                  <Badge
                    variant="secondary"
                    className="border-blue-300 bg-blue-100 text-xs font-medium text-blue-700"
                    data-english-extension="original-badge"
                  >
                    🎤 {languageSettings.primaryLanguage.toUpperCase()}
                  </Badge>
                  <div className="text-xs text-gray-500 opacity-0 transition-opacity group-hover:opacity-100">
                    Click to seek
                  </div>
                </div>
                <div
                  className="text-sm font-medium leading-relaxed text-gray-900"
                  data-english-extension="original-text"
                >
                  {currentSegment.text}
                </div>
              </div>

              {/* Secondary language (translated) */}
              {languageSettings.dualLanguageEnabled && (
                <div
                  className={`group cursor-pointer rounded-lg border p-4 transition-all hover:shadow-md ${
                    translatedSegment?.isTranslated
                      ? 'border-green-200 bg-green-50 hover:border-green-300 hover:bg-green-100'
                      : 'border-yellow-200 bg-yellow-50 hover:border-yellow-300 hover:bg-yellow-100'
                  }`}
                  onClick={handleSegmentClick}
                  data-english-extension="translated-segment"
                >
                  <div
                    className="mb-2 flex items-center justify-between"
                    data-english-extension="translation-badge-container"
                  >
                    {translatedSegment?.isTranslated ? (
                      <>
                        <Badge
                          variant="secondary"
                          className="border-green-300 bg-green-100 text-xs font-medium text-green-700"
                          data-english-extension="translated-badge"
                        >
                          ✨ {languageSettings.secondaryLanguage.toUpperCase()}
                        </Badge>
                      </>
                    ) : (
                      <>
                        <span
                          className="inline-block animate-spin text-xs text-yellow-600"
                          data-english-extension="translation-spinner"
                        >
                          ⟳
                        </span>
                        <Badge
                          variant="secondary"
                          className="border-yellow-300 bg-yellow-100 text-xs font-medium text-yellow-700"
                          data-english-extension="translating-badge"
                        >
                          🔄 Translating...
                        </Badge>
                      </>
                    )}
                    <div className="text-xs text-gray-500 opacity-0 transition-opacity group-hover:opacity-100">
                      Click to seek
                    </div>
                  </div>
                  <div
                    className={`text-sm leading-relaxed ${
                      !translatedSegment?.isTranslated ? 'italic text-gray-600' : 'font-medium text-gray-900'
                    }`}
                    data-english-extension="translated-text"
                  >
                    {translatedSegment?.translatedText || 'Translation in progress...'}
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="py-5 text-center italic text-gray-600" data-english-extension="no-transcript-message">
              Play video to see transcript
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
