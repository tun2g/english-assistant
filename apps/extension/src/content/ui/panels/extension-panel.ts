import { EXTENSION_CLASSES } from '../../../shared/constants/extension-constants';
import { getVideoContainer, removeElementIfExists } from '../../utils/dom/dom-utils';

export interface ExtensionPanelCallbacks {
  onOAuthConnect: () => Promise<void>;
  onTranscriptRequest: () => Promise<void>;
  onClose: () => void;
}

export class ExtensionPanel {
  private panel: HTMLElement | null = null;
  private callbacks: ExtensionPanelCallbacks;
  private isOAuthAuthenticated: boolean;

  constructor(callbacks: ExtensionPanelCallbacks, isOAuthAuthenticated: boolean) {
    this.callbacks = callbacks;
    this.isOAuthAuthenticated = isOAuthAuthenticated;
  }

  // Show extension panel
  show(): void {
    // Toggle panel - close if already open
    const existingPanel = document.querySelector(`.${EXTENSION_CLASSES.PANEL}`);
    if (existingPanel) {
      existingPanel.remove();
      return;
    }

    this.createPanel();
    this.setupEventListeners();
    this.attachToDOM();
    this.setupOutsideClickHandler();
  }

  // Hide extension panel
  hide(): void {
    removeElementIfExists(`.${EXTENSION_CLASSES.PANEL}`);
    this.panel = null;
  }

  // Create panel element
  private createPanel(): void {
    this.panel = document.createElement('div');
    this.panel.className = EXTENSION_CLASSES.PANEL;
    this.panel.style.cssText = `
      position: absolute;
      top: 60px;
      right: 20px;
      width: 350px;
      background: rgba(28, 28, 28, 0.95);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 12px;
      padding: 20px;
      z-index: 9999;
      color: white;
      font-family: "YouTube Sans", Roboto, sans-serif;
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    `;

    const oauthStatusHtml = this.createOAuthStatusHtml();
    
    this.panel.innerHTML = `
      <div class="panel-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
        <h3 style="margin: 0; font-size: 16px; font-weight: 600;">English Learning Assistant</h3>
        <button class="close-btn" style="background: none; border: none; color: white; font-size: 20px; cursor: pointer; padding: 4px;">×</button>
      </div>
      
      ${oauthStatusHtml}
      
      <div class="panel-content">
        <div class="transcript-section" style="margin-bottom: 16px;">
          <button class="transcript-btn" style="width: 100%; padding: 12px; background: #065fd4; border: none; border-radius: 8px; color: white; font-weight: 500; cursor: pointer; font-size: 14px;">
            📝 ${this.isOAuthAuthenticated ? 'Get Real Transcript' : 'Get Transcript (Auth Required)'}
          </button>
        </div>
      </div>
    `;
  }

  // Create OAuth status HTML
  private createOAuthStatusHtml(): string {
    if (this.isOAuthAuthenticated) {
      return `<div class="oauth-status" style="display: flex; align-items: center; gap: 8px; margin-bottom: 16px; padding: 8px; background: rgba(15, 157, 88, 0.2); border: 1px solid rgba(15, 157, 88, 0.4); border-radius: 6px; font-size: 12px;">
        <span style="color: #0f9d58;">✓</span>
        <span style="color: #ccc;">Connected to YouTube API</span>
      </div>`;
    } else {
      return `<div class="oauth-status" style="display: flex; align-items: center; gap: 8px; margin-bottom: 16px; padding: 8px; background: rgba(244, 67, 54, 0.2); border: 1px solid rgba(244, 67, 54, 0.4); border-radius: 6px; font-size: 12px;">
        <span style="color: #f44336;">🔒</span>
        <span style="color: #ccc;">Authentication required for real transcripts</span>
        <button class="oauth-connect-btn" style="margin-left: auto; background: #1976d2; border: none; color: white; padding: 4px 8px; border-radius: 4px; cursor: pointer; font-size: 11px;">Connect</button>
      </div>`;
    }
  }

  // Setup event listeners
  private setupEventListeners(): void {
    if (!this.panel) return;

    const closeBtn = this.panel.querySelector('.close-btn') as HTMLButtonElement;
    closeBtn.addEventListener('click', () => {
      this.hide();
      this.callbacks.onClose();
    });

    const oauthConnectBtn = this.panel.querySelector('.oauth-connect-btn') as HTMLButtonElement;
    if (oauthConnectBtn) {
      oauthConnectBtn.addEventListener('click', () => this.callbacks.onOAuthConnect());
    }

    const transcriptBtn = this.panel.querySelector('.transcript-btn') as HTMLButtonElement;
    transcriptBtn.addEventListener('click', () => this.callbacks.onTranscriptRequest());
  }

  // Attach panel to DOM
  private attachToDOM(): void {
    if (!this.panel) return;

    const videoContainer = getVideoContainer();
    if (videoContainer) {
      videoContainer.appendChild(this.panel);
    } else {
      document.body.appendChild(this.panel);
    }
  }

  // Setup outside click handler
  private setupOutsideClickHandler(): void {
    const closeOnOutsideClick = (e: Event) => {
      if (this.panel && !this.panel.contains(e.target as Node)) {
        this.hide();
        this.callbacks.onClose();
        document.removeEventListener('click', closeOnOutsideClick);
      }
    };

    setTimeout(() => {
      document.addEventListener('click', closeOnOutsideClick);
    }, 100);
  }
}