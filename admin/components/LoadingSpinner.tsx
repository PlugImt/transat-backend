import { memo } from 'react';
import { Loader2 } from 'lucide-react';
import clsx from 'clsx';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
  text?: string;
}

function LoadingSpinner({ 
  size = 'md', 
  className,
  text
}: LoadingSpinnerProps) {
  const sizeClasses = {
    sm: 'h-4 w-4',
    md: 'h-6 w-6',
    lg: 'h-8 w-8',
    xl: 'h-12 w-12',
  };

  return (
    <div className={clsx('flex flex-col items-center justify-center', className)}>
      <Loader2 className={clsx('animate-spin text-blue-600', sizeClasses[size])} />
      {text && (
        <p className="mt-2 text-sm text-gray-600">{text}</p>
      )}
    </div>
  );
}

export default memo(LoadingSpinner);

export const FullPageLoading = memo(function FullPageLoading({ text = "Chargement..." }: { text?: string }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <LoadingSpinner size="xl" text={text} />
    </div>
  );
});

export const PageLoading = memo(function PageLoading({ text = "Chargement..." }: { text?: string }) {
  return (
    <div className="p-8 flex items-center justify-center">
      <LoadingSpinner size="lg" text={text} />
    </div>
  );
});

export const ComponentLoading = memo(function ComponentLoading({ text = "Chargement..." }: { text?: string }) {
  return (
    <div className="flex items-center justify-center p-4">
      <LoadingSpinner size="md" text={text} />
    </div>
  );
});