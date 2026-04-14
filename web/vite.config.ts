import { defineConfig } from "vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [
    tanstackStart({
      srcDirectory: "app",
      vite: {
        plugins: [viteReact()],
      },
    }),
  ],
  server: {
    port: 4000,
  },
  resolve: {
    alias: {
      "~": `${import.meta.dirname}/app`,
    },
  },
});
