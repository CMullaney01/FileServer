import React, { createContext, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const AuthContext = createContext();

const AuthProvider = ({ children }) => {
  const navigate = useNavigate();

  useEffect(() => {
    const isAuthenticated = () => {
      // Implement your logic to check authentication status by inspecting cookies
      const sessionToken = document.cookie.replace(
        /(?:(?:^|.*;\s*)session_token\s*=\s*([^;]*).*$)|^.*$/,
        '$1'
      );
      return sessionToken !== '';
    };

    const checkAuthentication = async () => {
      if (!isAuthenticated()) {
        // Redirect to login if not authenticated
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
