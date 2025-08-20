import { cva, type VariantProps } from 'class-variance-authority';
import { forwardRef, type HTMLAttributes } from 'react';
import { cn } from '../lib/utils';

const spinnerVariants = cva('animate-spin rounded-full border-2 border-solid border-current border-r-transparent', {
  variants: {
    size: {
      sm: 'h-4 w-4',
      md: 'h-6 w-6',
      lg: 'h-8 w-8',
      xl: 'h-12 w-12',
    },
  },
  defaultVariants: {
    size: 'md',
  },
});

export interface SpinnerProps extends HTMLAttributes<HTMLDivElement>, VariantProps<typeof spinnerVariants> {}

export const Spinner = forwardRef<HTMLDivElement, SpinnerProps>(({ className, size, ...props }, ref) => {
  return <div ref={ref} className={cn(spinnerVariants({ size }), className)} {...props} />;
});
