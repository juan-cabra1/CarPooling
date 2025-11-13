import React from 'react';

export const Card = ({
  children,
  className = '',
  padding = 'md',
  hoverable = false,
  maxWidth = true,
}) => {
  const baseClasses = 'bg-white rounded-xl shadow-md border border-gray-100';
  const hoverClasses = hoverable
    ? 'hover:shadow-xl hover:scale-[1.02] transition-all duration-300 cursor-pointer'
    : 'transition-shadow duration-200';
  const maxWidthClass = maxWidth ? 'max-w-full' : '';

  const paddingClasses = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8',
    xl: 'p-10',
  };

  return (
    <div className={`${baseClasses} ${hoverClasses} ${paddingClasses[padding]} ${maxWidthClass} ${className}`}>
      {children}
    </div>
  );
};
