import { defineConfig } from "vite";
import fs from "fs";
import path from "path";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  build: {
    cssCodeSplit: false,
    outDir: "../public",
    rollupOptions: {
      output: {
        entryFileNames: "index.js",
        chunkFileNames: "index.js",
        assetFileNames: "index.css",
      },
    },
  },
  server: {
    https: {
      key: fs.readFileSync(path.resolve(__dirname, "../localhost-key.pem")),
      cert: fs.readFileSync(path.resolve(__dirname, "../localhost.pem")),
    },
    proxy: {
      "/api": {
        target: "http://127.0.0.1:3000",
        changeOrigin: true,
        rewrite: (path) => {
          console.log(path);
          return path;
        },
      },
    },
  },
});
