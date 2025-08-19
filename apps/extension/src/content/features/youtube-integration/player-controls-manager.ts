import { 
  EXTENSION_CLASSES, 
  YOUTUBE_SELECTORS 
} from '../../../shared/constants/extension-constants';
import { findYouTubeRightControls, createStyledButton } from '../../utils/dom/dom-utils';

export interface PlayerControlsCallbacks {
  onButtonClick: () => Promise<void>;
}

export class PlayerControlsManager {
  private callbacks: PlayerControlsCallbacks;
  private playerObserver: MutationObserver | null = null;

  constructor(callbacks: PlayerControlsCallbacks) {
    this.callbacks = callbacks;
  }

  // Inject extension button into YouTube player controls
  injectControls(): void {
    this.waitForPlayerControls();
    this.setupPlayerObserver();
  }

  // Remove extension controls
  removeControls(): void {
    const button = document.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`);
    if (button) {
      button.remove();
    }

    if (this.playerObserver) {
      this.playerObserver.disconnect();
      this.playerObserver = null;
    }
  }

  // Wait for player controls to appear
  private waitForPlayerControls(): void {
    const attemptInject = () => {
      console.log('Searching for YouTube player controls...');
      
      const rightControls = findYouTubeRightControls();

      if (rightControls && !rightControls.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`)) {
        console.log('Adding extension icon to controls');
        this.addExtensionIconToControls(rightControls);
      } else if (!rightControls) {
        console.log('Player controls not found, retrying...');
        setTimeout(attemptInject, 1000);
      } else {
        console.log('Extension icon already exists');
      }
    };

    attemptInject();
    setTimeout(attemptInject, 2000);
    setTimeout(attemptInject, 4000);
  }

  // Setup observer for dynamic player controls
  private setupPlayerObserver(): void {
    this.playerObserver = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        if (mutation.type === 'childList') {
          const rightControls = findYouTubeRightControls();
          if (rightControls && !rightControls.querySelector(`.${EXTENSION_CLASSES.EXTENSION_BTN}`)) {
            console.log('Detected right controls via MutationObserver');
            this.addExtensionIconToControls(rightControls);
            break;
          }
        }
      }
    });

    const moviePlayer = document.querySelector(YOUTUBE_SELECTORS.MOVIE_PLAYER);
    if (moviePlayer) {
      this.playerObserver.observe(moviePlayer, {
        childList: true,
        subtree: true
      });
    }
  }

  // Add extension icon to player controls
  private addExtensionIconToControls(controlsContainer: Element): void {
    console.log('Creating extension button...');
    
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
    extensionBtn.addEventListener('click', async (e) => {
      e.preventDefault();
      e.stopPropagation();
      console.log('Extension button clicked!');
      await this.callbacks.onButtonClick();
    });

    // Insert before settings button or append to end
    const settingsButton = controlsContainer.querySelector(YOUTUBE_SELECTORS.SETTINGS_BUTTON);
    if (settingsButton) {
      console.log('Inserting extension button before settings button');
      controlsContainer.insertBefore(extensionBtn, settingsButton);
    } else {
      console.log('Settings button not found, appending to end');
      controlsContainer.appendChild(extensionBtn);
    }

    console.log('Extension button added successfully!');
  }
}