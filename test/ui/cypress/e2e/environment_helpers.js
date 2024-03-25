export function verifyRequiredEnvVariables() {

  if (!Cypress.env('DEPLOYMENT_DRIVER_URL')) {
    throw new Error(`Missing required environment variable for DEPLOYMENT_DRIVER_URL`);
  }

  if (!Cypress.env('DEPLOYMENT_ENGINE_UI_PASSWORD')) {
    throw new Error('Missing required environment variable for DEPLOYMENT_ENGINE_UI_PASSWORD');
  }

  if (!Cypress.env('RH_SSO_URL')) {
    throw new Error('Missing required environment variable for RH_SSO_URL');
  }

  if (!Cypress.env('RH_ACCOUNT_USERNAME')) {
    throw new Error('Missing required environment variable for RH_ACCOUNT_USERNANE');
  }

  if (!Cypress.env('RH_ACCOUNT_PASSWORD')) {
    throw new Error('Missing required environment variable for RH_ACCOUNT_PASSWORD');

  }

	return 
}
