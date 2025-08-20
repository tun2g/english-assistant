// Types for notification system
export type NotificationType = 'success' | 'error' | 'warning' | 'loading';

export interface NotificationConfig {
  type: NotificationType;
  message: string;
  duration?: number;
}
