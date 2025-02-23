import { defineConfig, loadEnv, ConfigEnv, UserConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vitejs.dev/config/
export default defineConfig(({ mode }: ConfigEnv): UserConfig => {
  const env = loadEnv(mode, process.cwd(), "");

  const proxyConfig =
    mode === "development"
      ? {
          // In development, proxy all requests to localhost
          '/api/auth': {
            target: env.VITE_AUTH_BASE_URL || 'http://localhost:8000',
            changeOrigin: true,
            secure: false,
            rewrite: (path: string) => path.replace(/^\/api\/auth/, ''),
          },
          '/api/ws': {
            target: env.VITE_WS_BASE_URL || 'http://localhost:8001',
            ws: true,
            changeOrigin: true,
            secure: false,
            rewrite: (path: string) => path.replace(/^\/api\/ws/, ''),
          },
          '/api/chat': {
            target: env.VITE_CHAT_BASE_URL || 'http://localhost:8002',
            changeOrigin: true,
            secure: false,
            rewrite: (path: string) => path.replace(/^\/api\/chat/, ''),
          },
        }
      : undefined; // In production, no proxy needed as we'll use the full URLs

  return {
    plugins: [react(), tailwindcss()],
    server: {
      proxy: proxyConfig,
    },
  };
});
