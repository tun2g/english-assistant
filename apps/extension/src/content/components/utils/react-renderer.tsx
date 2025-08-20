import * as React from 'react';
import { createRoot, Root } from 'react-dom/client';
import { EXTENSION_CLASSES } from '../../../shared/constants';
import { captureError, logDebug, logWarning } from '../../utils/error-handler';

// Global registry of React roots for cleanup
const reactRoots = new Map<string, { container: HTMLElement; root: Root; shadowRoot?: ShadowRoot }>();

/**
 * Update global debug information
 */
function updateGlobalDebugInfo(): void {
  try {
    // @ts-ignore
    window.__englishExtensionReactRoots = reactRoots.size;

    // Expose component debug info if debug object exists
    if (window.__englishExtensionDebug) {
      window.__englishExtensionDebug.getComponentDebugInfo = getComponentDebugInfo;
    }
  } catch (error) {
    // Ignore debug info update errors
  }
}

/**
 * Wait for DOM to be ready for manipulation
 */
function waitForDOM(): Promise<void> {
  return new Promise(resolve => {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => resolve());
    } else {
      resolve();
    }
  });
}

/**
 * Create shadow root with full Tailwind CSS injection
 */
async function createShadowContainer(
  className: string,
  useShadowDOM: boolean = true
): Promise<{
  container: HTMLElement;
  renderTarget: HTMLElement | ShadowRoot;
}> {
  const container = document.createElement('div');
  container.className = className;

  // Add extension-specific attributes for debugging
  container.setAttribute('data-english-extension', 'true');
  container.setAttribute('data-component-id', className);

  if (useShadowDOM) {
    try {
      const shadowRoot = container.attachShadow({ mode: 'open' });

      // Inject Tailwind CSS into Shadow DOM
      await injectTailwindCSS(shadowRoot);

      return { container, renderTarget: shadowRoot };
    } catch (error) {
      logWarning('Shadow DOM creation failed, falling back to regular DOM', error, 'ReactRenderer');
    }
  }

  // Fallback to regular DOM with style isolation
  container.style.cssText = `
    all: initial;
    display: block !important;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif !important;
    line-height: 1.5 !important;
    -webkit-text-size-adjust: 100% !important;
    box-sizing: border-box !important;
    z-index: 2147483647 !important;
  `;

  return { container, renderTarget: container };
}

/**
 * Inject CSS content into Shadow DOM with proper reset and scoping
 */
function injectCSSIntoShadowDOM(shadowRoot: ShadowRoot, cssContent: string, source = 'CSS'): void {
  try {
    // Create style element for Shadow DOM
    const style = document.createElement('style');
    style.textContent = `
      /* Reset styles for Shadow DOM */
      :host {
        all: initial;
        display: block;
        font-family: "Inter", system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        line-height: 1.5;
        box-sizing: border-box;
        color: #1E293B;
        background: transparent;
      }
      
      /* Box sizing for all elements */
      *, *::before, *::after {
        box-sizing: border-box;
      }
      
      /* Extension ${source} */
      ${cssContent}
    `;

    shadowRoot.appendChild(style);
    logDebug(`${source} injected into Shadow DOM`, undefined, 'ReactRenderer');
  } catch (error) {
    captureError(`Failed to inject ${source} into Shadow DOM`, error, 'css-injection', 'ReactRenderer');
  }
}

/**
 * Inject multiple CSS sources into Shadow DOM (supports inline imports)
 */
async function injectMultipleCSSources(
  shadowRoot: ShadowRoot,
  cssSources: (string | Promise<string>)[]
): Promise<void> {
  try {
    for (let i = 0; i < cssSources.length; i++) {
      const source = cssSources[i];
      let cssContent: string;

      if (typeof source === 'string') {
        cssContent = source;
      } else {
        cssContent = await source;
      }

      injectCSSIntoShadowDOM(shadowRoot, cssContent, `CSS Source ${i + 1}`);
    }
  } catch (error) {
    captureError('Failed to inject multiple CSS sources into Shadow DOM', error, 'css-injection', 'ReactRenderer');
  }
}

/**
 * Inject Tailwind CSS into Shadow DOM (legacy method for backward compatibility)
 */
async function injectTailwindCSS(shadowRoot: ShadowRoot): Promise<void> {
  try {
    // Get the extension's CSS file URL
    const cssUrl = chrome.runtime.getURL('style.css');

    // Fetch the CSS content
    const response = await fetch(cssUrl);
    const cssContent = await response.text();

    injectCSSIntoShadowDOM(shadowRoot, cssContent, 'Tailwind CSS');
  } catch (error) {
    captureError('Failed to inject Tailwind CSS into Shadow DOM', error, 'css-injection', 'ReactRenderer');
  }
}

/**
 * Render a React component into the DOM for content scripts
 */
export async function renderReactComponent(
  component: React.ReactElement,
  className: string,
  containerId?: string,
  options: {
    useShadowDOM?: boolean;
    appendToElement?: Element;
    inlineCSS?: (string | Promise<string>)[];
  } = {}
): Promise<void> {
  const id = containerId || className;
  const { useShadowDOM = true, appendToElement, inlineCSS } = options; // Default to true for Shadow DOM

  try {
    // Wait for DOM to be ready
    await waitForDOM();

    logDebug(`Starting render for ${id}`, undefined, 'ReactRenderer');

    // Clean up existing container if it exists
    const existing = reactRoots.get(id);
    if (existing) {
      logDebug(`Cleaning up existing component ${id}`, undefined, 'ReactRenderer');
      try {
        existing.root.unmount();
        existing.container.remove();
      } catch (error) {
        captureError(`Error cleaning up existing component ${id}`, error, 'cleanup', 'ReactRenderer');
      }
      reactRoots.delete(id);
    }

    // Remove any existing container with same class
    const domContainer = document.querySelector(`.${className}`);
    if (domContainer) {
      domContainer.remove();
    }

    // Create new container with shadow DOM and CSS injection
    const { container, renderTarget } = await createShadowContainer(className, useShadowDOM);
    const shadowRoot = renderTarget !== container ? (renderTarget as ShadowRoot) : undefined;

    // Inject additional inline CSS if provided
    if (shadowRoot && inlineCSS && inlineCSS.length > 0) {
      await injectMultipleCSSources(shadowRoot, inlineCSS);
    }

    // Append to specified element or body
    const parentElement = appendToElement || document.body;
    if (!parentElement) {
      throw new Error('ReactRenderer: No parent element found to append container');
    }

    parentElement.appendChild(container);
    logDebug(`Container created and appended to ${parentElement.tagName}`, undefined, 'ReactRenderer');

    // Create React root and render component
    const root = createRoot(renderTarget as Element);

    // Wrap component rendering in error boundary
    const ComponentWithErrorBoundary = () => {
      try {
        return component;
      } catch (error) {
        captureError(`Component render error for ${id}`, error, 'component-render', 'ReactRenderer');
        return React.createElement(
          'div',
          {
            style: {
              padding: '10px',
              backgroundColor: 'rgba(255, 0, 0, 0.8)',
              color: 'white',
              fontSize: '12px',
              fontFamily: 'monospace',
              borderRadius: '4px',
              border: '2px solid red',
              zIndex: '9999999',
              position: 'fixed',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%, -50%)',
            },
          },
          `Extension Error: Failed to render ${id}`
        );
      }
    };

    root.render(React.createElement(ComponentWithErrorBoundary));

    // Store root and container for cleanup
    reactRoots.set(id, { container, root, shadowRoot });

    // Update global debug info
    updateGlobalDebugInfo();

    logDebug(`Successfully rendered component ${id}`, undefined, 'ReactRenderer');
  } catch (error) {
    captureError(`Failed to render component ${id}`, error, 'render', 'ReactRenderer');
    throw error;
  }
}

/**
 * Unmount a React component from the DOM
 */
export function unmountReactComponent(className: string, containerId?: string): void {
  const id = containerId || className;

  try {
    logDebug(`Unmounting component ${id}`, undefined, 'ReactRenderer');

    // Unmount React component
    const existing = reactRoots.get(id);
    if (existing) {
      try {
        existing.root.unmount();
      } catch (error) {
        captureError(`Error unmounting React root for ${id}`, error, 'unmount', 'ReactRenderer');
      }

      try {
        existing.container.remove();
      } catch (error) {
        captureError(`Error removing container for ${id}`, error, 'dom-cleanup', 'ReactRenderer');
      }

      reactRoots.delete(id);
      updateGlobalDebugInfo();
      logDebug(`Successfully unmounted ${id}`, undefined, 'ReactRenderer');
    }

    // Also remove any DOM element with this class
    const domContainer = document.querySelector(`.${className}`);
    if (domContainer) {
      domContainer.remove();
      logDebug(`Removed DOM container with class ${className}`, undefined, 'ReactRenderer');
    }
  } catch (error) {
    captureError(`Failed to unmount component ${id}`, error, 'unmount', 'ReactRenderer');
  }
}

/**
 * Check if a React component is currently mounted
 */
export function isReactComponentMounted(className: string, containerId?: string): boolean {
  const id = containerId || className;
  const hasReactRoot = reactRoots.has(id);
  const hasDOMElement = document.querySelector(`.${className}`) !== null;

  // Component is considered mounted only if both React root exists and DOM element is present
  return hasReactRoot && hasDOMElement;
}

/**
 * Clean up all React components (useful for extension cleanup)
 */
export function cleanupAllReactComponents(): void {
  logDebug('Cleaning up all React components', undefined, 'ReactRenderer');

  reactRoots.forEach(({ container, root, shadowRoot }, id) => {
    try {
      // Unmount React root
      root.unmount();

      // Clean up shadow DOM if it exists
      if (shadowRoot) {
        // Clear shadow root content
        while (shadowRoot.firstChild) {
          shadowRoot.removeChild(shadowRoot.firstChild);
        }
      }

      // Remove container
      container.remove();

      logDebug(`Unmounted ${id}`, undefined, 'ReactRenderer');
    } catch (error) {
      captureError(`Failed to unmount ${id}`, error, 'cleanup-all', 'ReactRenderer');
    }
  });

  reactRoots.clear();
  updateGlobalDebugInfo();

  // Also remove any leftover extension elements
  Object.values(EXTENSION_CLASSES).forEach(className => {
    const elements = document.querySelectorAll(`.${className}`);
    elements.forEach(element => {
      try {
        element.remove();
        logDebug(`Removed leftover element with class ${className}`, undefined, 'ReactRenderer');
      } catch (error) {
        captureError(`Failed to remove element ${className}`, error, 'dom-cleanup', 'ReactRenderer');
      }
    });
  });

  // Also remove any elements with our extension data attribute
  const extensionElements = document.querySelectorAll('[data-english-extension="true"]');
  extensionElements.forEach(element => {
    try {
      element.remove();
      logDebug('Removed leftover extension element', undefined, 'ReactRenderer');
    } catch (error) {
      captureError('Failed to remove extension element', error, 'dom-cleanup', 'ReactRenderer');
    }
  });
}

/**
 * Export utility functions for CSS injection (for advanced use cases)
 */
export { injectCSSIntoShadowDOM, injectMultipleCSSources, injectTailwindCSS };

/**
 * Get debug information about currently mounted components
 */
export function getComponentDebugInfo(): Array<{
  id: string;
  className: string;
  hasDOMElement: boolean;
  hasReactRoot: boolean;
  hasShadowRoot: boolean;
}> {
  const info: Array<{
    id: string;
    className: string;
    hasDOMElement: boolean;
    hasReactRoot: boolean;
    hasShadowRoot: boolean;
  }> = [];

  reactRoots.forEach(({ container, shadowRoot }, id) => {
    info.push({
      id,
      className: container.className,
      hasDOMElement: document.contains(container),
      hasReactRoot: true,
      hasShadowRoot: !!shadowRoot,
    });
  });

  return info;
}
