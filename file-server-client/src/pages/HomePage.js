import React from 'react';
import { AuthProvider } from '../services/AuthContext';
import FileList from '../components/FileList';

function HomePage() {
  return (
    <AuthProvider>
        <div>
            Home Page
            <FileList />
        </div>
    </AuthProvider>
  );
}

export default HomePage