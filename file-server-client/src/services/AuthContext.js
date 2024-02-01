import React, { createContext, useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const AuthContext = createContext();

const AuthProvider = ({ children }) => {
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuthentication = async () => {
      try {
        // Make a request to the server's authentication status endpoint
        const response = await axios.get('http://localhost:8080/authstatus', {
          withCredentials: true, // Include credentials (cookies) in the request
        });

        if (response.status === 200) {
          // User is authenticated
          console.log('User is authenticated');
        } else {
          // User is not authenticated
          console.log('User is not authenticated');
          // Redirect to login page
          navigate('/login');
        }
      } catch (error) {
        // An error occurred (e.g., network error, server error)
        console.error('Error checking authentication status:', error);
        // Redirect to login page in case of an error
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
