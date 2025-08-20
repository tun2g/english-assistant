// Types for language settings and preferences
import { SUPPORTED_LANGUAGES } from '@english/shared';

type SupportedLanguageCode = (typeof SUPPORTED_LANGUAGES)[number]['code'];

export interface LanguageSettingParams {
  primaryLanguage: SupportedLanguageCode;
  secondaryLanguage: SupportedLanguageCode;
  dualLanguageEnabled: boolean;
  autoTranslateEnabled: boolean;
}

export interface LanguageDisplayOptions {
  showOriginal: boolean;
  showTranslated: boolean;
  highlightActive: boolean;
  fontSize: 'small' | 'medium' | 'large';
  position: 'bottom' | 'side' | 'overlay';
}
