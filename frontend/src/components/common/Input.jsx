import { forwardRef } from 'react';

export const Input = forwardRef(
  ({ label, error, helperText, fullWidth = false, className = '', ...props }, ref) => {
    const baseClasses = 'px-4 py-2.5 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 transition-all duration-200 text-gray-900 placeholder-gray-400 min-h-[42px]';
    const errorClasses = error
      ? 'border-red-500 focus:ring-red-500 focus:border-red-500 animate-shake'
      : 'border-gray-300 hover:border-gray-400';
    const widthClass = fullWidth ? 'w-full' : '';

    return (
      <div className={`${fullWidth ? 'w-full' : ''}`}>
        {label && (
          <label className="block text-sm font-semibold text-gray-700 mb-2">
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={`${baseClasses} ${errorClasses} ${widthClass} ${className}`}
          {...props}
        />
        {error && (
          <p className="mt-2 text-sm text-red-600 flex items-center gap-1">
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            {error}
          </p>
        )}
        {helperText && !error && (
          <p className="mt-2 text-sm text-gray-500">{helperText}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
