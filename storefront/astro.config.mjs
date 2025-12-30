// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';
import node from '@astrojs/node';

// https://astro.build/config
export default defineConfig({
  output: 'server',
  vite: {
    plugins: [tailwindcss()],
    define: {
      // Expose environment variables to client-side code
      'import.meta.env.PUBLIC_API_URL': JSON.stringify(process.env.PUBLIC_API_URL || 'http://localhost:8080/api/v1'),
      'import.meta.env.PUBLIC_MIDTRANS_CLIENT_KEY': JSON.stringify(process.env.PUBLIC_MIDTRANS_CLIENT_KEY || ''),
      'import.meta.env.PUBLIC_MIDTRANS_ENVIRONMENT': JSON.stringify(process.env.PUBLIC_MIDTRANS_ENVIRONMENT || 'sandbox'),
    }
  },

  adapter: node({
    mode: 'standalone'
  })
});