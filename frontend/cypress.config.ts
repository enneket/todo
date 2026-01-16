import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: "http://localhost:5173", // Using the frontend dev server port
    setupNodeEvents() {
      // implement node event listeners here
    },
  },
});
