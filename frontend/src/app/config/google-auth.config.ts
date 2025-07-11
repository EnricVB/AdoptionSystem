// Google Authentication Configuration
// This file contains the configuration for Google Sign-In integration

export interface GoogleAuthConfig {
  clientId: string;
  enabled: boolean;
}

// Google OAuth Client ID configuration
// In production, this should be set via environment variables
const googleAuthConfig: GoogleAuthConfig = {
  // Replace with your actual Google Client ID
  // You can get this from Google Cloud Console > APIs & Services > Credentials
  clientId: '1017473621019-9hbmho8kqgq7pjhvjl4nqsjq6kc6q5qv.apps.googleusercontent.com',
  enabled: true
};

export default googleAuthConfig;