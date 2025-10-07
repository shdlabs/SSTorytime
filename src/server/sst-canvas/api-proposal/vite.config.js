import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss()],
  server: {
    host: true, // Listen on all network interfaces
    port: 8000, // Specify the port number
  },
});
