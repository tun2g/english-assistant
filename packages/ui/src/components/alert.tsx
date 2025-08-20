import { cva, type VariantProps } from 'class-variance-authority';
import { forwardRef, type HTMLAttributes } from 'react';
import { cn } from '../lib/utils';

const alertVariants = cva(
  'relative w-full rounded-lg border p-4 [&>svg~*]:pl-7 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground',
  {
    variants: {
      variant: {
        default: 'bg-background text-foreground',
        destructive: 'border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive',
        warning:
          'border-yellow-500/50 text-yellow-900 dark:text-yellow-100 bg-yellow-50 dark:bg-yellow-900/20 [&>svg]:text-yellow-600',
        success:
          'border-green-500/50 text-green-900 dark:text-green-100 bg-green-50 dark:bg-green-900/20 [&>svg]:text-green-600',
        info: 'border-blue-500/50 text-blue-900 dark:text-blue-100 bg-blue-50 dark:bg-blue-900/20 [&>svg]:text-blue-600',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
);

export const Alert = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement> & VariantProps<typeof alertVariants>>(
  ({ className, variant, ...props }, ref) => {
    return <div ref={ref} role="alert" className={cn(alertVariants({ variant }), className)} {...props} />;
  }
);

export const AlertTitle = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => {
    return <h5 ref={ref} className={cn('mb-1 font-medium leading-none tracking-tight', className)} {...props} />;
  }
);

export const AlertDescription = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLParagraphElement>>(
  ({ className, ...props }, ref) => {
    return <div ref={ref} className={cn('text-sm [&_p]:leading-relaxed', className)} {...props} />;
  }
);
