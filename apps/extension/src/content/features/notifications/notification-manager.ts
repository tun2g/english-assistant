import { 
  EXTENSION_CLASSES, 
  NOTIFICATION_TYPES, 
  NOTIFICATION_DURATIONS 
} from '../../../shared/constants/extension-constants';
import type { NotificationType } from '../../../shared/types/extension-types';

export class NotificationManager {
  // Show notification with different types
  static show(
    message: string, 
    type: NotificationType = NOTIFICATION_TYPES.SUCCESS
  ): void {
    // Remove existing notification
    const existing = document.querySelector(`.${EXTENSION_CLASSES.NOTIFICATION}`);
    if (existing) {
      existing.remove();
    }

    const notification = document.createElement('div');
    notification.className = EXTENSION_CLASSES.NOTIFICATION;

    const colors = {
      [NOTIFICATION_TYPES.SUCCESS]: '#4caf50',
      [NOTIFICATION_TYPES.ERROR]: '#f44336',
      [NOTIFICATION_TYPES.WARNING]: '#ff9800',
      [NOTIFICATION_TYPES.LOADING]: '#2196f3'
    };

    const icons = {
      [NOTIFICATION_TYPES.SUCCESS]: '✓',
      [NOTIFICATION_TYPES.ERROR]: '✗',
      [NOTIFICATION_TYPES.WARNING]: '⚠',
      [NOTIFICATION_TYPES.LOADING]: '⟳'
    };

    notification.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: ${colors[type]};
      color: white;
      padding: 12px 16px;
      border-radius: 8px;
      z-index: 10001;
      font-family: "YouTube Sans", Roboto, sans-serif;
      font-size: 14px;
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
      display: flex;
      align-items: center;
      gap: 8px;
      max-width: 300px;
    `;

    notification.innerHTML = `
      <span style="font-size: 16px;">${icons[type]}</span>
      <span>${message}</span>
    `;

    document.body.appendChild(notification);

    // Auto-hide notification (except loading)
    if (type !== NOTIFICATION_TYPES.LOADING) {
      const duration = NOTIFICATION_DURATIONS[type.toUpperCase() as keyof typeof NOTIFICATION_DURATIONS];
      setTimeout(() => {
        notification.remove();
      }, duration);
    }
  }

  // Show success notification
  static showSuccess(message: string): void {
    this.show(message, NOTIFICATION_TYPES.SUCCESS);
  }

  // Show error notification
  static showError(message: string): void {
    this.show(message, NOTIFICATION_TYPES.ERROR);
  }

  // Show warning notification
  static showWarning(message: string): void {
    this.show(message, NOTIFICATION_TYPES.WARNING);
  }

  // Show loading notification
  static showLoading(message: string): void {
    this.show(message, NOTIFICATION_TYPES.LOADING);
  }

  // Hide current notification
  static hide(): void {
    const notification = document.querySelector(`.${EXTENSION_CLASSES.NOTIFICATION}`);
    if (notification) {
      notification.remove();
    }
  }
}