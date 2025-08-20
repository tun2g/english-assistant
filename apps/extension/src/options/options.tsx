import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import {
  Button,
  Input,
  Switch,
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Label,
  Separator,
} from '@english/ui';
import '../styles/globals.scss';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000,
    },
  },
});

// Options page component with @english/ui
function OptionsApp() {
  const [autoTranslate, setAutoTranslate] = useState(false);
  const [saveHistory, setSaveHistory] = useState(true);
  const [apiKey, setApiKey] = useState('');

  return (
    <QueryClientProvider client={queryClient}>
      <div className="bg-background min-h-screen p-6">
        <div className="mx-auto max-w-2xl space-y-6">
          <div className="space-y-2">
            <h1 className="text-2xl font-bold">Settings</h1>
            <p className="text-muted-foreground">Configure your English Learning Assistant extension.</p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>General Settings</CardTitle>
              <CardDescription>Configure basic extension behavior</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="api-key">API Key</Label>
                <Input
                  id="api-key"
                  type="password"
                  placeholder="Enter your translation API key"
                  value={apiKey}
                  onChange={e => setApiKey(e.target.value)}
                />
              </div>

              <Separator />

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Auto-translate</Label>
                  <p className="text-muted-foreground text-sm">Automatically translate text when detected</p>
                </div>
                <Switch checked={autoTranslate} onCheckedChange={setAutoTranslate} />
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Save translation history</Label>
                  <p className="text-muted-foreground text-sm">Keep a record of your translations</p>
                </div>
                <Switch checked={saveHistory} onCheckedChange={setSaveHistory} />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Data Management</CardTitle>
              <CardDescription>Manage your extension data and settings</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 gap-4">
                <Button variant="outline">Export Data</Button>
                <Button variant="outline">Import Data</Button>
                <Button variant="destructive">Clear All Data</Button>
              </div>

              <p className="text-muted-foreground text-sm">Settings are automatically saved when changed.</p>
            </CardContent>
          </Card>
        </div>
      </div>
    </QueryClientProvider>
  );
}

// Mount the React app
const container = document.getElementById('app');
if (container) {
  const root = ReactDOM.createRoot(container);
  root.render(<OptionsApp />);
}
