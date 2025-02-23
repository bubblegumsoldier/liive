import axios, { InternalAxiosRequestConfig, AxiosError } from 'axios';

// Create separate API instances for different services
export const authApi = axios.create({
  baseURL: '/api/auth',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const chatApi = axios.create({
  baseURL: '/api/chat',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add JWT token to requests
const addAuthToken = (config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

// Handle 401 errors
const handle401Error = async (error: AxiosError) => {
  if (error.response?.status === 401) {
    localStorage.removeItem('token');
    window.location.href = '/login';
  }
  return Promise.reject(error);
};

// Add interceptors to both API instances
[authApi, chatApi].forEach(api => {
  api.interceptors.request.use(addAuthToken, error => Promise.reject(error));
  api.interceptors.response.use(response => response, handle401Error);
});

export default {
  auth: authApi,
  chat: chatApi,
}; 