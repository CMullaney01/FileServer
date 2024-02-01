import React, { createContext, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthService from './AuthService';

const AuthContext = createContext();

const AuthProvider = ({ children }) => {
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuthentication = async () => {
      if (!AuthService.isAuthenticated()) {
        // Redirect to the login page if not authenticated
        navigate('/login');
      }
    };

    checkAuthentication();
  }, [navigate]);

  return <AuthContext.Provider value={{}}>{children}</AuthContext.Provider>;
};

const useAuth = () => {
  return useContext(AuthContext);
};

export { AuthProvider, useAuth };