import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    env: {
      baseUrl: 'Deployment Engine UI Url',
      DEPLOYMENT_ENGINE_UI_PASSWORD: 'Admin password to login Deployment Engine UI',
      RH_SSO_URL: 'https://sso.redhat.com',
      RH_ACCOUNT_USERNAME: 'User to login https://sso.redhat.com',
      RH_ACCOUNT_PASSWORD: 'Password to login https://sso.redhat.com',
    },
  }
});
