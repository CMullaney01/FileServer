import React from 'react';
import { AuthProvider } from '../services/AuthContext';
import FileList from '../components/FileList';

function HomePage() {
  return (
    <div>
      Home Page
      <FileList />
    </div>
  );
}

export default () => (
  <AuthProvider>
    <HomePage />
  </AuthProvider>
);