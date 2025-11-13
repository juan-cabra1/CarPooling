import React from 'react';
import { Navbar } from './Navbar';
import { Footer } from './Footer';

export const MainLayout = ({ children, maxWidth = 'container' }) => {
  const maxWidthClasses = {
    container: 'max-w-container',
    content: 'max-w-content',
    form: 'max-w-form',
    full: 'max-w-full',
  };

  return (
    <div className="flex flex-col min-h-screen bg-gray-50">
      <Navbar />
      <main className={`flex-grow container mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 lg:py-12 ${maxWidthClasses[maxWidth]}`}>
        {children}
      </main>
      <Footer />
    </div>
  );
};
