import { EXTENSION_CLASSES } from '../../../shared/constants';
import { removeElementIfExists } from '../../utils/dom/dom-utils';
import type { VideoTranscript } from '../../../shared/types/extension-types';

export class TranscriptModal {
  // Show transcript in a modal overlay
  static show(transcript: VideoTranscript): void {
    // Remove existing modal if present
    removeElementIfExists(`.${EXTENSION_CLASSES.MODAL}`);

    const modal = document.createElement('div');
    modal.className = EXTENSION_CLASSES.MODAL;
    modal.style.cssText = `
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      z-index: 10000;
      display: flex;
      justify-content: center;
      align-items: center;
      font-family: "YouTube Sans", Roboto, sans-serif;
    `;

    const modalContent = document.createElement('div');
    modalContent.style.cssText = `
      background: white;
      max-width: 80%;
      max-height: 80%;
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    `;

    const header = this.createHeader();
    const content = this.createContent(transcript);

    modalContent.appendChild(header);
    modalContent.appendChild(content);
    modal.appendChild(modalContent);

    this.setupEventListeners(modal);
    document.body.appendChild(modal);
  }

  // Create modal header
  private static createHeader(): HTMLElement {
    const header = document.createElement('div');
    header.style.cssText = `
      padding: 20px;
      background: #f9f9f9;
      border-bottom: 1px solid #eee;
      display: flex;
      justify-content: space-between;
      align-items: center;
    `;

    const title = document.createElement('h2');
    title.textContent = 'Video Transcript';
    title.style.cssText = 'margin: 0; color: #333; font-size: 18px;';

    const closeBtn = document.createElement('button');
    closeBtn.textContent = '×';
    closeBtn.className = 'modal-close-btn';
    closeBtn.style.cssText = `
      background: none;
      border: none;
      font-size: 24px;
      cursor: pointer;
      color: #666;
      padding: 0;
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
    `;

    header.appendChild(title);
    header.appendChild(closeBtn);
    return header;
  }

  // Create modal content
  private static createContent(transcript: VideoTranscript): HTMLElement {
    const content = document.createElement('div');
    content.style.cssText = `
      padding: 20px;
      max-height: 400px;
      overflow-y: auto;
      color: #333;
      line-height: 1.5;
    `;

    // Display transcript segments
    if (transcript.segments && transcript.segments.length > 0) {
      const transcriptText = transcript.segments.map(segment => segment.text).join(' ');

      content.innerHTML = `
        <p style="margin: 0 0 10px 0; color: #666; font-size: 12px;">
          Language: ${transcript.language || 'Unknown'} | 
          Segments: ${transcript.segments.length} | 
          Source: ${transcript.available ? 'Real YouTube Captions' : 'Fallback'}
        </p>
        <div style="font-size: 14px; white-space: pre-wrap;">${transcriptText}</div>
      `;
    } else {
      content.textContent = 'No transcript segments available';
    }

    return content;
  }

  // Setup event listeners for modal
  private static setupEventListeners(modal: HTMLElement): void {
    const closeBtn = modal.querySelector('.modal-close-btn') as HTMLButtonElement;

    closeBtn.addEventListener('click', () => modal.remove());

    // Close modal when clicking outside
    modal.addEventListener('click', e => {
      if (e.target === modal) {
        modal.remove();
      }
    });
  }

  // Hide transcript modal
  static hide(): void {
    removeElementIfExists(`.${EXTENSION_CLASSES.MODAL}`);
  }
}
