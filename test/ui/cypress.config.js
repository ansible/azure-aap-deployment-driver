const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    supportFile: false,
  },
  env: {
    "DEPLOYMENT_DRIVER_URL": "http://localhost:9999"
  }
});
