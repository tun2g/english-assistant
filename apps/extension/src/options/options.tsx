import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';
import { 
  App, 
  View, 
  Page, 
  Navbar, 
  Block, 
  Button, 
  List, 
  ListInput,
  Toggle,
  BlockTitle 
} from 'framework7-react';
import 'framework7/css/bundle';

// Options page component with Framework7
function OptionsApp() {
  const [autoTranslate, setAutoTranslate] = useState(false);
  const [saveHistory, setSaveHistory] = useState(true);
  const [apiKey, setApiKey] = useState('');

  const f7params = {
    name: 'English Learning Assistant Settings',
    theme: 'auto',
  };

  return (
    <App {...f7params}>
      <View main className="safe-areas" url="/">
        <Page>
          <Navbar title="Settings" />
          
          <BlockTitle>General Settings</BlockTitle>
          <List>
            <ListInput
              label="API Key"
              type="password"
              placeholder="Enter your translation API key"
              value={apiKey}
              onInput={(e) => setApiKey(e.target.value)}
            />
            <li>
              <div className="item-content">
                <div className="item-inner">
                  <div className="item-title">Auto-translate</div>
                  <div className="item-after">
                    <Toggle 
                      checked={autoTranslate}
                      onChange={setAutoTranslate}
                    />
                  </div>
                </div>
              </div>
            </li>
            <li>
              <div className="item-content">
                <div className="item-inner">
                  <div className="item-title">Save translation history</div>
                  <div className="item-after">
                    <Toggle 
                      checked={saveHistory}
                      onChange={setSaveHistory}
                    />
                  </div>
                </div>
              </div>
            </li>
          </List>
          
          <BlockTitle>Data Management</BlockTitle>
          <Block>
            <div className="grid grid-cols-1 gap-4">
              <Button fill color="blue">
                Export Data
              </Button>
              <Button fill color="orange">
                Import Data
              </Button>
              <Button fill color="red">
                Clear All Data
              </Button>
            </div>
          </Block>
          
          <Block>
            <p className="text-sm text-gray-600">
              Settings are automatically saved when changed.
            </p>
          </Block>
        </Page>
      </View>
    </App>
  );
}

// Mount the React app
const container = document.getElementById('app');
if (container) {
  const root = ReactDOM.createRoot(container);
  root.render(<OptionsApp />);
}