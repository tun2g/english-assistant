import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { ClassValue } from 'clsx';

/**
 * Utility function for concatenating Tailwind CSS classes
 * Following the coding convention rule for class concatenation
 */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}