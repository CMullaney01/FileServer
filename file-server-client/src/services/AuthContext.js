import React, { createContext, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthService from './AuthService';

const AuthContext = createContext();

const AuthProvider = ({ children }) => {
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuthentication = async () => {
      if (!AuthService.isAuthenticated()) {
        try {
          // Make a request to the server to get the redirection URL
          const response = await fetch('http://localhost:8080/getRedirectionURL');
          const data = await response.json();

          // Redirect to the URL from the server response
          if (data.redirectionURL) {
            navigate(data.redirectionURL);
          } else {
            // Fallback to a default login page if the response doesn't contain a URL
            navigate('/signin');
          }
        } catch (error) {
          console.error('Error fetching redirection URL:', error);
          // Fallback to a default login page in case of an error
          navigate('/signin');
        }
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