import { EXTENSION_CLASSES, YOUTUBE_SELECTORS } from '../../../shared/constants';
import { createStyledButton, findYouTubeRightControls } from '../../utils/dom/dom-utils';

export interface PlayerControlsCallbacks {
  onButtonClick: () => Promise<void>;
}

export class PlayerControlsManager {
  private callbacks: PlayerControlsCallbacks;
  private playerObserver: MutationObserver | null = null;
  private injectionAttempts: number = 0;
  private isInjected: boolean = false;
  private maxAttempts: number = 10;

  constructor(callbacks: PlayerControlsCallbacks) {
    this.callbacks = callbacks;
  }

  // Inject extension button into YouTube player controls
  injectControls(): void {
    // Reset state for new injection attempt
    this.injectionAttempts = 0;
    this.isInjected = false;
    this.waitForPlayerControls();
    this.setupPlayerObserver();
  }

  // Remove extension controls
  removeControls(): void {
    const button = document.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`);
    if (button) {
      button.remove();
    }

    // Reset state
    this.isInjected = false;
    this.injectionAttempts = 0;

    if (this.playerObserver) {
      this.playerObserver.disconnect();
      this.playerObserver = null;
    }
  }

  // Wait for player controls to appear with smart retry logic
  private waitForPlayerControls(): void {
    const attemptInject = () => {
      // Skip if already injected or too many attempts
      if (this.isInjected || this.injectionAttempts >= this.maxAttempts) {
        return;
      }

      this.injectionAttempts++;

      const rightControls = findYouTubeRightControls(true); // Silent mode to prevent spam

      if (rightControls) {
        const existingButton = rightControls.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`);
        if (!existingButton) {
          this.addExtensionIconToControls(rightControls);
        } else {
          this.isInjected = true;
        }
      } else {
        const delay = Math.min(1000 * this.injectionAttempts, 5000); // Progressive delay, max 5s
        setTimeout(attemptInject, delay);
      }
    };

    attemptInject();
  }

  // Setup observer for dynamic player controls
  private setupPlayerObserver(): void {
    if (this.isInjected) {
      return; // Don't setup observer if already injected
    }

    this.playerObserver = new MutationObserver(mutations => {
      // Skip if already injected or too many attempts
      if (this.isInjected || this.injectionAttempts >= this.maxAttempts) {
        return;
      }

      // Throttle observer calls - only check on significant changes
      let hasSignificantChange = false;
      for (const mutation of mutations) {
        if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
          // Check if any added nodes are control-related
          for (const node of mutation.addedNodes) {
            if (node.nodeType === Node.ELEMENT_NODE) {
              const element = node as Element;
              if (
                element.matches('.ytp-right-controls, .ytp-chrome-controls') ||
                element.querySelector('.ytp-right-controls')
              ) {
                hasSignificantChange = true;
                break;
              }
            }
          }
          if (hasSignificantChange) break;
        }
      }

      if (hasSignificantChange) {
        const rightControls = findYouTubeRightControls(true);
        if (rightControls && !rightControls.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`)) {
          this.addExtensionIconToControls(rightControls);
        }
      }
    });

    const moviePlayer = document.querySelector(YOUTUBE_SELECTORS.MOVIE_PLAYER);
    if (moviePlayer) {
      this.playerObserver.observe(moviePlayer, {
        childList: true,
        subtree: true,
      });
    }
  }

  // Add extension icon to player controls
  private addExtensionIconToControls(controlsContainer: Element): void {
    // Prevent duplicate injection
    if (this.isInjected || controlsContainer.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`)) {
      this.isInjected = true;
      return;
    }

    const extensionBtn = createStyledButton(
      EXTENSION_CLASSES.EXTENSION_BTN,
      'English Learning Assistant',
      `<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" style="pointer-events: none;">
        <path d="M12 3L1 9L12 15L21 10.09V17H23V9M5 13.18V17.18L12 21L19 17.18V13.18L12 17L5 13.18Z"/>
      </svg>`,
      `
        width: 48px !important;
        height: 48px !important;
        padding: 8px !important;
        margin: 0 !important;
        border: none !important;
        background: none !important;
        color: white !important;
        cursor: pointer !important;
        opacity: 0.9 !important;
        transition: opacity 0.1s cubic-bezier(0.05,0,0,1) !important;
        display: inline-block !important;
        position: relative !important;
        text-align: center !important;
        vertical-align: top !important;
        outline: none !important;
      `
    );

    // Add hover effects
    extensionBtn.addEventListener('mouseenter', () => {
      extensionBtn.style.opacity = '1';
    });

    extensionBtn.addEventListener('mouseleave', () => {
      extensionBtn.style.opacity = '0.9';
    });

    // Add click handler
    extensionBtn.addEventListener('click', async e => {
      e.preventDefault();
      e.stopPropagation();
      await this.callbacks.onButtonClick();
    });

    // Try to find a reference button for insertion
    let insertionSuccess = false;

    // First, try to insert in the right controls left section (before settings)
    const rightControlsLeft = controlsContainer.querySelector(YOUTUBE_SELECTORS.RIGHT_CONTROLS_LEFT);
    if (rightControlsLeft) {
      // Try settings button first in the left section
      for (const selector of YOUTUBE_SELECTORS.SETTINGS_BUTTON) {
        const settingsButton = rightControlsLeft.querySelector(selector);
        if (settingsButton) {
          try {
            rightControlsLeft.insertBefore(extensionBtn, settingsButton);
            insertionSuccess = true;
            break;
          } catch (error) {
            console.warn(`Failed to insert before settings button in left section with selector ${selector}:`, error);
          }
        }
      }

      // If no settings button in left section, append to end of left section
      if (!insertionSuccess) {
        try {
          rightControlsLeft.appendChild(extensionBtn);
          insertionSuccess = true;
        } catch (error) {
          console.warn('Failed to append to right controls left section:', error);
        }
      }
    }

    // If left section insertion failed, try settings button in main container
    if (!insertionSuccess) {
      for (const selector of YOUTUBE_SELECTORS.SETTINGS_BUTTON) {
        const settingsButton = controlsContainer.querySelector(selector);
        if (settingsButton && settingsButton.parentNode === controlsContainer) {
          try {
            controlsContainer.insertBefore(extensionBtn, settingsButton);
            insertionSuccess = true;
            break;
          } catch (error) {
            console.warn(`Failed to insert before settings button with selector ${selector}:`, error);
          }
        }
      }
    }

    // If settings button insertion failed, try fullscreen button
    if (!insertionSuccess) {
      for (const selector of YOUTUBE_SELECTORS.FULLSCREEN_BUTTON) {
        const fullscreenButton = controlsContainer.querySelector(selector);
        if (fullscreenButton && fullscreenButton.parentNode === controlsContainer) {
          try {
            controlsContainer.insertBefore(extensionBtn, fullscreenButton);
            insertionSuccess = true;
            break;
          } catch (error) {
            console.warn(`Failed to insert before fullscreen button with selector ${selector}:`, error);
          }
        }
      }
    }

    // Fallback: append to end
    if (!insertionSuccess) {
      try {
        controlsContainer.appendChild(extensionBtn);
        insertionSuccess = true;
      } catch (error) {
        console.error('Failed to append extension button:', error);
        return;
      }
    }

    // Verify injection was successful
    const verifyButton = controlsContainer.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`);
    if (verifyButton && insertionSuccess) {
      this.isInjected = true;

      // Disconnect observer since we succeeded
      if (this.playerObserver) {
        this.playerObserver.disconnect();
        this.playerObserver = null;
      }
    }
  }
}
