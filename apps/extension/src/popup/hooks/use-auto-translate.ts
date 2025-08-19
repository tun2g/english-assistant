import { useState, useEffect } from 'react';
import { EXTENSION_STORAGE_KEYS, EXTENSION_MESSAGES } from '../../shared/constants/extension-constants';

export function useAutoTranslate() {
  const [isEnabled, setIsEnabled] = useState(false);

  useEffect(() => {
    const loadSetting = async () => {
      try {
        const result = await chrome.storage.local.get([EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]);
        setIsEnabled(result[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED] || false);
      } catch (error) {
        console.error('Failed to load auto-translate setting:', error);
      }
    };

    loadSetting();
  }, []);

  const toggle = async (enabled: boolean): Promise<void> => {
    try {
      // Save setting
      await chrome.storage.local.set({
        [EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]: enabled
      });
      
      setIsEnabled(enabled);

      // Notify content script
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      if (tab.id) {
        await chrome.tabs.sendMessage(tab.id, {
          action: EXTENSION_MESSAGES.TOGGLE_TRANSLATION,
          enabled
        });
      }
    } catch (error) {
      console.error('Failed to toggle auto-translate:', error);
      throw error;
    }
  };

  return {
    isEnabled,
    toggle
  };
}