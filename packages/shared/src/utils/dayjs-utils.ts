import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import relativeTime from 'dayjs/plugin/relativeTime';
import duration from 'dayjs/plugin/duration';

// Extend dayjs with plugins
dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(relativeTime);
dayjs.extend(duration);

/**
 * Customized dayjs instance following coding convention
 * This should be used instead of primitive dayjs or Date()
 */
export const dayjsLib = dayjs;

export function formatDate(date: string | Date, format = 'YYYY-MM-DD'): string {
  return dayjsLib(date).format(format);
}

export function formatDateTime(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjsLib(date).format(format);
}

export function formatRelativeTime(date: string | Date): string {
  return dayjsLib(date).fromNow();
}

export function isValidDate(date: string | Date): boolean {
  return dayjsLib(date).isValid();
}

export function getCurrentTimestamp(): string {
  return dayjsLib().toISOString();
}

export function addDays(date: string | Date, days: number): string {
  return dayjsLib(date).add(days, 'day').toISOString();
}

export function subtractDays(date: string | Date, days: number): string {
  return dayjsLib(date).subtract(days, 'day').toISOString();
}