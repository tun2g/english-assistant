import { EXTENSION_CLASSES } from '../../../shared/constants';

export interface AuthModalCallbacks {
  onConnect: () => Promise<void>;
  onCancel: () => void;
}

export class AuthModal {
  private modal: HTMLElement | null = null;
  private callbacks: AuthModalCallbacks;

  constructor(callbacks: AuthModalCallbacks) {
    this.callbacks = callbacks;
  }

  // Show authentication modal
  show(): void {
    // Remove existing modal
    this.hide();

    this.modal = document.createElement('div');
    this.modal.className = EXTENSION_CLASSES.AUTH_MODAL;
    this.modal.style.cssText = `
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      z-index: 10001;
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: "YouTube Sans", Roboto, sans-serif;
    `;

    this.modal.innerHTML = `
      <div class="auth-modal-content" style="
        background: #1c1c1c;
        padding: 32px;
        border-radius: 12px;
        max-width: 400px;
        width: 90%;
        text-align: center;
        color: white;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6);
      ">
        <div style="font-size: 48px; margin-bottom: 16px;">🔒</div>
        <h2 style="margin: 0 0 16px 0; font-size: 20px; font-weight: 600;">YouTube Authentication Required</h2>
        <p style="margin: 0 0 24px 0; color: #ccc; line-height: 1.5;">
          To access real YouTube transcripts, you need to connect your YouTube account. This is secure and uses Google OAuth.
        </p>
        <div style="display: flex; gap: 12px; justify-content: center;">
          <button class="auth-cancel-btn" style="
            padding: 12px 24px;
            border: 1px solid #555;
            background: transparent;
            color: #ccc;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
          ">Cancel</button>
          <button class="auth-connect-btn" style="
            padding: 12px 24px;
            background: #1976d2;
            border: none;
            color: white;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 500;
            font-size: 14px;
          ">Connect YouTube Account</button>
        </div>
      </div>
    `;

    this.setupEventListeners();
    document.body.appendChild(this.modal);
  }

  // Hide authentication modal
  hide(): void {
    const existingModal = document.querySelector(`.${EXTENSION_CLASSES.AUTH_MODAL}`);
    if (existingModal) {
      existingModal.remove();
    }
    this.modal = null;
  }

  // Setup event listeners for modal buttons
  private setupEventListeners(): void {
    if (!this.modal) return;

    const cancelBtn = this.modal.querySelector('.auth-cancel-btn') as HTMLButtonElement;
    const connectBtn = this.modal.querySelector('.auth-connect-btn') as HTMLButtonElement;

    cancelBtn.addEventListener('click', () => {
      this.hide();
      this.callbacks.onCancel();
    });

    connectBtn.addEventListener('click', async () => {
      try {
        connectBtn.textContent = 'Connecting...';
        connectBtn.disabled = true;

        await this.callbacks.onConnect();
        this.hide();
      } catch (error) {
        console.error('OAuth connection failed:', error);
        connectBtn.textContent = 'Connect YouTube Account';
        connectBtn.disabled = false;
      }
    });

    // Close on outside click
    this.modal.addEventListener('click', e => {
      if (e.target === this.modal) {
        this.hide();
        this.callbacks.onCancel();
      }
    });
  }
}
