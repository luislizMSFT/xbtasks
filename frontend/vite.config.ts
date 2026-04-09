import path from "node:path";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";
import tailwindcss from "@tailwindcss/vite";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), wails("./bindings"), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    // Force dependency re-optimization on every dev server start
    force: true,
  },
  optimizeDeps: {
    // Don't cache pre-bundled deps between restarts
    force: true,
  },
});
