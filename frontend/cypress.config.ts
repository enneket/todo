import { defineConfig } from "cypress";
import task from "@cypress/code-coverage/task";

export default defineConfig({
  e2e: {
    baseUrl: "http://localhost:5173", // Using the frontend dev server port
    setupNodeEvents(on, config) {
      task(on, config);
      
      on('task', {
        log(message) {
          console.log(message)
          return null
        },
        table(message) {
          console.table(message)
          return null
        }
      })

      // It's IMPORTANT to return the config object
      // with any changed environment variables
      return config;
    },
  },
});
