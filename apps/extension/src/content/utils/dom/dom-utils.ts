import { YOUTUBE_SELECTORS } from '../../../shared/constants/extension-constants';

// Wait for an element to appear in the DOM
export function waitForElement(
  selector: string, 
  timeout: number = 10000,
  parent: Document | Element = document
): Promise<Element | null> {
  return new Promise((resolve) => {
    const element = parent.querySelector(selector);
    if (element) {
      resolve(element);
      return;
    }

    const observer = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        if (mutation.type === 'childList') {
          const foundElement = parent.querySelector(selector);
          if (foundElement) {
            observer.disconnect();
            resolve(foundElement);
            return;
          }
        }
      }
    });

    observer.observe(parent, {
      childList: true,
      subtree: true
    });

    // Timeout fallback
    setTimeout(() => {
      observer.disconnect();
      resolve(null);
    }, timeout);
  });
}

// Find YouTube player right controls with multiple fallback selectors
export function findYouTubeRightControls(): Element | null {
  for (const selector of YOUTUBE_SELECTORS.RIGHT_CONTROLS) {
    const element = document.querySelector(selector);
    if (element) {
      console.log(`Found YouTube controls using selector: ${selector}`);
      return element;
    }
  }
  return null;
}

// Get video container for positioning overlays
export function getVideoContainer(): Element | null {
  return document.querySelector(YOUTUBE_SELECTORS.MOVIE_PLAYER) || 
         document.querySelector(YOUTUBE_SELECTORS.VIDEO_PLAYER);
}

// Remove element if it exists
export function removeElementIfExists(selector: string): void {
  const element = document.querySelector(selector);
  if (element) {
    element.remove();
  }
}

// Create styled button element
export function createStyledButton(
  className: string,
  title: string,
  innerHTML: string,
  styles: string
): HTMLButtonElement {
  const button = document.createElement('button');
  button.className = className;
  button.setAttribute('title', title);
  button.setAttribute('aria-label', title);
  button.innerHTML = innerHTML;
  button.style.cssText = styles;
  
  return button;
}