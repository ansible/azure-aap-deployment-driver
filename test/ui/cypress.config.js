const { defineConfig } = require('cypress')

module.exports = defineConfig({
  // These settings apply everywhere unless overridden
  viewportWidth: 1920,
  viewportHeight: 1080,
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  }
})
