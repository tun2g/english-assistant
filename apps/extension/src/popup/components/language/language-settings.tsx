import { SUPPORTED_LANGUAGES } from '@english/shared';
import {
  Alert,
  AlertDescription,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Switch,
} from '@english/ui';
import { useEffect, useState } from 'react';
import { EXTENSION_STORAGE_KEYS } from '../../../shared/constants';
import type { LanguageSettingParams } from '../../../shared/types/language-types';

export function LanguageSettings() {
  const [settings, setSettings] = useState<LanguageSettingParams>({
    primaryLanguage: 'en',
    secondaryLanguage: 'vi',
    dualLanguageEnabled: false,
    autoTranslateEnabled: false,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Load settings from storage
  useEffect(() => {
    const loadSettings = async () => {
      try {
        const result = await chrome.storage.local.get([
          EXTENSION_STORAGE_KEYS.DUAL_LANGUAGE_ENABLED,
          EXTENSION_STORAGE_KEYS.PRIMARY_LANGUAGE,
          EXTENSION_STORAGE_KEYS.SECONDARY_LANGUAGE,
          EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED,
        ]);

        setSettings({
          primaryLanguage: result[EXTENSION_STORAGE_KEYS.PRIMARY_LANGUAGE] || 'en',
          secondaryLanguage: result[EXTENSION_STORAGE_KEYS.SECONDARY_LANGUAGE] || 'es',
          dualLanguageEnabled: result[EXTENSION_STORAGE_KEYS.DUAL_LANGUAGE_ENABLED] || false,
          autoTranslateEnabled: result[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED] || false,
        });
      } catch (err) {
        console.error('Failed to load language settings:', err);
        setError('Failed to load settings');
      } finally {
        setIsLoading(false);
      }
    };

    loadSettings();
  }, []);

  // Save settings to storage
  const saveSettings = async (newSettings: Partial<LanguageSettingParams>) => {
    try {
      setError(null);
      const updatedSettings = { ...settings, ...newSettings };

      await chrome.storage.local.set({
        [EXTENSION_STORAGE_KEYS.PRIMARY_LANGUAGE]: updatedSettings.primaryLanguage,
        [EXTENSION_STORAGE_KEYS.SECONDARY_LANGUAGE]: updatedSettings.secondaryLanguage,
        [EXTENSION_STORAGE_KEYS.DUAL_LANGUAGE_ENABLED]: updatedSettings.dualLanguageEnabled,
        [EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]: updatedSettings.autoTranslateEnabled,
      });

      setSettings(updatedSettings);

      // Notify content script of settings change
      const [tab] = await chrome.tabs.query({
        active: true,
        currentWindow: true,
      });
      if (tab.id) {
        chrome.tabs
          .sendMessage(tab.id, {
            action: 'UPDATE_LANGUAGE_SETTINGS',
            settings: updatedSettings,
          })
          .catch(() => {
            // Content script might not be available, that's okay
          });
      }
    } catch (err) {
      console.error('Failed to save language settings:', err);
      setError('Failed to save settings');
    }
  };

  if (isLoading) {
    return (
      <div className="border-b p-4">
        <div className="text-sm text-gray-500">Loading language settings...</div>
      </div>
    );
  }

  return (
    <div className="space-y-4 border-b p-4">
      <div className="flex flex-col space-y-2">
        <h3 className="text-sm font-medium">Language Settings</h3>

        {/* Primary Language */}
        <div className="space-y-1">
          <Label htmlFor="primary-language" className="text-xs">
            Primary Language
          </Label>
          <Select
            value={settings.primaryLanguage}
            onValueChange={value => saveSettings({ primaryLanguage: value as any })}
          >
            <SelectTrigger id="primary-language" className="h-8">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {SUPPORTED_LANGUAGES.map(language => (
                <SelectItem key={language.code} value={language.code}>
                  {language.flag} {language.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Dual Language Toggle */}
        <div className="flex items-center justify-between">
          <div className="space-y-0">
            <Label htmlFor="dual-language" className="text-xs">
              Dual Language Mode
            </Label>
            <p className="text-xs text-gray-500">Show original and translated text</p>
          </div>
          <Switch
            id="dual-language"
            checked={settings.dualLanguageEnabled}
            onCheckedChange={checked => saveSettings({ dualLanguageEnabled: checked })}
          />
        </div>

        {/* Secondary Language (only if dual language is enabled) */}
        {settings.dualLanguageEnabled && (
          <div className="space-y-1">
            <Label htmlFor="secondary-language" className="text-xs">
              Secondary Language
            </Label>
            <Select
              value={settings.secondaryLanguage}
              onValueChange={value => saveSettings({ secondaryLanguage: value as any })}
            >
              <SelectTrigger id="secondary-language" className="h-8">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {SUPPORTED_LANGUAGES.map(language => (
                  <SelectItem key={language.code} value={language.code}>
                    {language.flag} {language.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}

        {error && (
          <Alert variant="destructive">
            <AlertDescription className="text-xs">{error}</AlertDescription>
          </Alert>
        )}
      </div>
    </div>
  );
}
