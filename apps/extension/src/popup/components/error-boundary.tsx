import React, { Component, ErrorInfo, ReactNode } from 'react';
import { Card, CardContent, CardHeader, CardTitle, Button } from '@english/ui';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: ErrorInfo;
  debugInfo: {
    timestamp: string;
    userAgent: string;
    chromeVersion?: string;
    extensionContext: string;
  };
}

export class PopupErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);

    // Collect debug information
    const debugInfo = {
      timestamp: new Date().toISOString(),
      userAgent: navigator.userAgent,
      chromeVersion: navigator.userAgent.match(/Chrome\/(\d+\.\d+\.\d+\.\d+)/)?.[1],
      extensionContext: 'popup',
    };

    this.state = {
      hasError: false,
      debugInfo,
    };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Log to extension storage for debugging
    this.logErrorToStorage(error, errorInfo);

    this.setState({
      error,
      errorInfo,
    });
  }

  private async logErrorToStorage(error: Error, errorInfo: ErrorInfo) {
    try {
      const errorLog = {
        timestamp: new Date().toISOString(),
        error: {
          message: error.message,
          stack: error.stack,
          name: error.name,
        },
        errorInfo: {
          componentStack: errorInfo.componentStack,
        },
        debugInfo: this.state.debugInfo,
      };

      // Store in Chrome extension storage
      if (typeof chrome !== 'undefined' && chrome.storage) {
        const { extensionErrorLogs = [] } = await new Promise<any>(resolve => {
          chrome.storage.local.get(['extensionErrorLogs'], resolve);
        });

        // Keep only last 10 error logs
        const updatedLogs = [errorLog, ...extensionErrorLogs].slice(0, 10);

        chrome.storage.local.set({ extensionErrorLogs: updatedLogs });
      }
    } catch (storageError) {
      // Silent error handling
    }
  }

  private handleRetry = () => {
    this.setState({
      hasError: false,
      error: undefined,
      errorInfo: undefined,
    });
  };

  private handleOpenDevTools = () => {
    // Instructions for user
    alert(
      'Error details logged to extension storage. To inspect:\n' +
        '1. Right-click extension icon\n' +
        '2. Select "Inspect popup"\n' +
        '3. Check Console tab for error details'
    );
  };

  render() {
    if (this.state.hasError) {
      return (
        <div className="popup-container w-80 p-4">
          <Card>
            <CardHeader>
              <CardTitle className="text-red-600">Extension Error</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="text-sm">
                <p className="mb-2 font-medium text-red-700">Something went wrong with the popup</p>
                <p className="mb-3 text-gray-600">{this.state.error?.message || 'Unknown error occurred'}</p>

                <div className="rounded bg-gray-50 p-3 text-xs">
                  <p>
                    <strong>Time:</strong> {this.state.debugInfo.timestamp}
                  </p>
                  <p>
                    <strong>Chrome:</strong> {this.state.debugInfo.chromeVersion}
                  </p>
                  <p>
                    <strong>Context:</strong> {this.state.debugInfo.extensionContext}
                  </p>
                </div>
              </div>

              <div className="flex gap-2 pt-2">
                <Button onClick={this.handleRetry} size="sm" className="flex-1">
                  Try Again
                </Button>
                <Button onClick={this.handleOpenDevTools} variant="outline" size="sm" className="flex-1">
                  Debug Info
                </Button>
              </div>

              <div className="border-t pt-2 text-xs text-gray-500">
                <p>Extension is still working on YouTube pages.</p>
                <p>This only affects the popup interface.</p>
              </div>
            </CardContent>
          </Card>
        </div>
      );
    }

    return this.props.children;
  }
}
