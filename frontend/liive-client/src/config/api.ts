const isDevelopment = import.meta.env.MODE === 'development'

// In development, we use relative URLs because of the proxy
// In production, we use the full URLs from environment variables
export const API_CONFIG = {
  auth: {
    baseUrl: isDevelopment
      ? import.meta.env.VITE_AUTH_API
      : `${import.meta.env.VITE_AUTH_BASE_URL}${import.meta.env.VITE_AUTH_API}`,
  },
  ws: {
    baseUrl: isDevelopment
      ? import.meta.env.VITE_WS_API
      : `${import.meta.env.VITE_WS_BASE_URL}${import.meta.env.VITE_WS_API}`,
  },
  chat: {
    baseUrl: isDevelopment
      ? import.meta.env.VITE_CHAT_API
      : `${import.meta.env.VITE_CHAT_BASE_URL}${import.meta.env.VITE_CHAT_API}`,
  },
} 