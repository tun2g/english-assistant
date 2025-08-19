// Background Service Worker Template for English Learning Assistant
// TODO: Implement features as needed

console.log('English Learning Assistant: Background script loaded');

// Extension lifecycle events
chrome.runtime.onStartup.addListener(() => {
  console.log('Extension started');
  // TODO: Initialize extension
});

chrome.runtime.onInstalled.addListener((details: chrome.runtime.InstalledDetails) => {
  console.log('Extension installed/updated:', details.reason);
  // TODO: Handle installation/updates
});

// Message handling
chrome.runtime.onMessage.addListener((request: any, _sender: chrome.runtime.MessageSender, sendResponse: (response?: any) => void) => {
  console.log('Message received:', request);
  
  switch (request.action) {
    case 'OPEN_TAB':
      // Handle OAuth tab opening from content script
      if (request.url) {
        chrome.tabs.create({
          url: request.url,
          active: true
        });
        sendResponse({ success: true, message: 'Tab opened' });
      } else {
        sendResponse({ success: false, message: 'No URL provided' });
      }
      break;
      
    case 'OPEN_POPUP':
      // Handle popup opening request
      chrome.action.openPopup();
      sendResponse({ success: true, message: 'Popup opened' });
      break;
      
    default:
      sendResponse({ success: true, message: 'Unknown action' });
  }
  
  return true; // Indicates async response
});

// Context menu template
chrome.runtime.onInstalled.addListener(() => {
  // TODO: Create context menu items if needed
});

// Tab management template
chrome.tabs.onUpdated.addListener((_tabId: number, _changeInfo: chrome.tabs.TabChangeInfo, _tab: chrome.tabs.Tab) => {
  // TODO: Handle tab updates if needed
});

// Keyboard shortcuts template
chrome.commands.onCommand.addListener((command: string) => {
  console.log('Command received:', command);
  // TODO: Handle keyboard shortcuts
});