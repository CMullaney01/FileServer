import axios from 'axios';

const AuthService = {
  // Function to authenticate the user
  login: async (username, password) => {
    try {
      const response = await axios.post('http://localhost:8080/login', {
        username,
        password,
      });

      // Assuming the server returns a token upon successful login
      const token = response.data.token;

      // Store the token in local storage or a secure storage mechanism
      localStorage.setItem('token', token);

      return token;
    } catch (error) {
      console.error('Error logging in:', error);
      throw error;
    }
  },

  // Function to check if the user is authenticated
  isAuthenticated: () => {
    // Check if the token exists in local storage or a secure storage mechanism
    const token = localStorage.getItem('token');
    return !!token; // Returns true if the token exists, false otherwise
  },

  // Function to log out the user
  logout: () => {
    // Remove the token from local storage or a secure storage mechanism
    localStorage.removeItem('token');
  },
};

export default AuthService;