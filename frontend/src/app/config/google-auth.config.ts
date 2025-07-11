// Google Authentication Configuration
// This file contains the configuration for Google Sign-In integration

export interface GoogleAuthConfig {
  clientId: string;
  enabled: boolean;
}

// Google OAuth Client ID configuration
// In production, this should be set via environment variables
const googleAuthConfig: GoogleAuthConfig = {
  clientId: '800054744191-9a91feuu075kn7f4rigapeqvgvp2nl00.apps.googleusercontent.com',
  enabled: true
};

export default googleAuthConfig;