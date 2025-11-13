import React from 'react';

export const Loading = ({ size = 'md', fullScreen = false, text = '' }) => {
  const sizeClasses = {
    sm: 'h-8 w-8 border-2',
    md: 'h-12 w-12 border-4',
    lg: 'h-16 w-16 border-4',
    xl: 'h-20 w-20 border-4',
  };

  const spinner = (
    <div className="flex flex-col items-center justify-center gap-3">
      <div
        className={`${sizeClasses[size]} animate-spin rounded-full border-gray-200 border-t-primary-600 border-r-primary-400`}
      ></div>
      {text && (
        <p className="text-gray-600 font-medium animate-pulse">{text}</p>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-95 backdrop-blur-sm flex items-center justify-center z-50">
        {spinner}
      </div>
    );
  }

  return <div className="flex items-center justify-center w-full py-8">{spinner}</div>;
};
