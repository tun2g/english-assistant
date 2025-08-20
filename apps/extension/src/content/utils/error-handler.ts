// Error handling and debugging utilities for Chrome extension content scripts

export interface ExtensionError {
  message: string;
  error?: Error;
  context?: string;
  timestamp: number;
  userAgent: string;
  url: string;
  componentId?: string;
}

class ExtensionErrorHandler {
  private static instance: ExtensionErrorHandler;
  private errors: ExtensionError[] = [];
  private maxErrors = 100;
  private enableConsoleLogging = true;

  private constructor() {
    // Listen for unhandled errors
    window.addEventListener('error', event => {
      this.captureError('Uncaught Error', event.error, 'window.error');
    });

    // Listen for unhandled promise rejections
    window.addEventListener('unhandledrejection', event => {
      this.captureError('Unhandled Promise Rejection', event.reason, 'promise.rejection');
    });
  }

  static getInstance(): ExtensionErrorHandler {
    if (!ExtensionErrorHandler.instance) {
      ExtensionErrorHandler.instance = new ExtensionErrorHandler();
    }
    return ExtensionErrorHandler.instance;
  }

  // Capture and log errors
  captureError(message: string, error?: Error | any, context?: string, componentId?: string): void {
    const extensionError: ExtensionError = {
      message,
      error: error instanceof Error ? error : new Error(String(error)),
      context,
      componentId,
      timestamp: Date.now(),
      userAgent: navigator.userAgent,
      url: window.location.href,
    };

    this.errors.push(extensionError);

    // Keep only the most recent errors
    if (this.errors.length > this.maxErrors) {
      this.errors = this.errors.slice(-this.maxErrors);
    }

    // Log to console if enabled
    if (this.enableConsoleLogging) {
      const prefix = `[English Extension${componentId ? ` - ${componentId}` : ''}]`;
      const contextMsg = context ? ` (${context})` : '';

      console.error(`${prefix} ${message}${contextMsg}:`, error);
      console.error('Error details:', {
        message,
        context,
        componentId,
        timestamp: new Date(extensionError.timestamp).toISOString(),
        url: extensionError.url,
        stack: error?.stack,
      });
    }

    // Send error to background script for potential reporting
    try {
      chrome.runtime
        .sendMessage({
          action: 'ERROR_REPORT',
          data: extensionError,
        })
        .catch(() => {
          // Ignore if background script is not available
        });
    } catch (error) {
      // Ignore chrome.runtime errors
    }
  }

  // Get recent errors for debugging
  getRecentErrors(count: number = 10): ExtensionError[] {
    return this.errors.slice(-count);
  }

  // Clear error history
  clearErrors(): void {
    this.errors = [];
  }

  // Enable/disable console logging
  setConsoleLogging(enabled: boolean): void {
    this.enableConsoleLogging = enabled;
  }

  // Get error statistics
  getErrorStats(): {
    totalErrors: number;
    recentErrors: number;
    errorsByContext: Record<string, number>;
    errorsByComponent: Record<string, number>;
  } {
    const now = Date.now();
    const recentThreshold = now - 5 * 60 * 1000; // Last 5 minutes

    const errorsByContext: Record<string, number> = {};
    const errorsByComponent: Record<string, number> = {};
    let recentErrors = 0;

    this.errors.forEach(error => {
      if (error.timestamp > recentThreshold) {
        recentErrors++;
      }

      const context = error.context || 'unknown';
      errorsByContext[context] = (errorsByContext[context] || 0) + 1;

      const component = error.componentId || 'unknown';
      errorsByComponent[component] = (errorsByComponent[component] || 0) + 1;
    });

    return {
      totalErrors: this.errors.length,
      recentErrors,
      errorsByContext,
      errorsByComponent,
    };
  }
}

// Export singleton instance
export const errorHandler = ExtensionErrorHandler.getInstance();

// Convenience functions
export function captureError(message: string, error?: Error | any, context?: string, componentId?: string): void {
  errorHandler.captureError(message, error, context, componentId);
}

export function logDebug(message: string, data?: any, componentId?: string): void {
  const prefix = `[English Extension${componentId ? ` - ${componentId}` : ''}]`;
  console.log(`${prefix} ${message}`, data || '');
}

export function logWarning(message: string, data?: any, componentId?: string): void {
  const prefix = `[English Extension${componentId ? ` - ${componentId}` : ''}]`;
  console.warn(`${prefix} ${message}`, data || '');
}

// Add to window for debugging
declare global {
  interface Window {
    __englishExtensionDebug?: {
      errorHandler: ExtensionErrorHandler;
      getComponentDebugInfo?: () => any[];
    };
  }
}

// Expose debug functions to window for easier debugging
window.__englishExtensionDebug = {
  errorHandler,
};
