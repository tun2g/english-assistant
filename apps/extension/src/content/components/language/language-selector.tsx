import { SUPPORTED_LANGUAGES } from '@english/shared';
import {
  Button,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Label,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Switch,
} from '@english/ui';
import { useState } from 'react';
import type { LanguageSettings } from '../../../shared/types/extension-types';

interface LanguageSelectorProps {
  initialSettings: LanguageSettings;
  onSave: (settings: LanguageSettings) => void;
  onClose: () => void;
}

export function LanguageSelector({ initialSettings, onSave, onClose }: LanguageSelectorProps) {
  const [settings, setSettings] = useState<LanguageSettings>(initialSettings);

  const handleSave = () => {
    onSave(settings);
    onClose();
  };

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  // Remove createPortal since we're now rendered within Shadow DOM
  return (
    <div
      className="z-extension-modal font-inter animate-fade-in fixed inset-0 flex items-center justify-center bg-black/80"
      onClick={handleBackdropClick}
      data-english-extension="language-selector-backdrop"
    >
      <Card
        className="bg-card border-border animate-scale-in max-h-[80vh] w-[480px] max-w-[90%] overflow-y-auto rounded-lg border shadow-xl"
        onClick={e => e.stopPropagation()}
        data-english-extension="language-selector-card"
      >
        <CardHeader
          className="flex flex-row items-center justify-between space-y-0 px-6 pb-2 pt-6"
          data-english-extension="language-selector-header"
        >
          <CardTitle className="text-card-foreground text-xl font-semibold">Language Settings</CardTitle>
          <Button
            variant="ghost"
            size="sm"
            onClick={onClose}
            className="text-muted-foreground hover:text-foreground h-6 w-6 p-0"
            data-english-extension="close-button"
          >
            ✕
          </Button>
        </CardHeader>

        <CardContent className="space-y-6 px-6 pb-6 pt-0" data-english-extension="language-selector-content">
          {/* Dual Language Toggle */}
          <div className="flex items-center justify-between" data-english-extension="dual-language-toggle">
            <div className="space-y-0.5">
              <Label className="text-base">Enable Dual Language Display</Label>
              <CardDescription>Show both original and translated text side by side</CardDescription>
            </div>
            <Switch
              checked={settings.dualLanguageEnabled}
              onCheckedChange={checked =>
                setSettings({
                  ...settings,
                  dualLanguageEnabled: checked,
                })
              }
            />
          </div>

          {/* Primary Language */}
          <div className="space-y-2">
            <Label>Primary Language</Label>
            <Select
              value={settings.primaryLanguage}
              onValueChange={(value: string) =>
                setSettings({
                  ...settings,
                  primaryLanguage: value as any,
                })
              }
            >
              <SelectTrigger>
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

          {/* Secondary Language (only if dual language is enabled) */}
          {settings.dualLanguageEnabled && (
            <div className="space-y-2">
              <Label>Secondary Language</Label>
              <Select
                value={settings.secondaryLanguage}
                onValueChange={(value: string) =>
                  setSettings({
                    ...settings,
                    secondaryLanguage: value as any,
                  })
                }
              >
                <SelectTrigger>
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

          {/* Auto Translate Toggle */}
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label className="text-base">Auto-translate Segments</Label>
              <CardDescription>Automatically translate transcript segments as they play</CardDescription>
            </div>
            <Switch
              checked={settings.autoTranslateEnabled}
              onCheckedChange={checked =>
                setSettings({
                  ...settings,
                  autoTranslateEnabled: checked,
                })
              }
            />
          </div>
        </CardContent>

        {/* Footer */}
        <div className="flex justify-end gap-3 px-6 pb-6 pt-0" data-english-extension="language-selector-footer">
          <Button variant="outline" onClick={onClose} className="px-4 py-2" data-english-extension="cancel-button">
            Cancel
          </Button>
          <Button onClick={handleSave} className="px-4 py-2" data-english-extension="save-button">
            Save Changes
          </Button>
        </div>
      </Card>
    </div>
  );
}
