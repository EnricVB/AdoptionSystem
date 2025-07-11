# Google Authentication Implementation

This document describes the Google Sign-In authentication implementation for the AdoptionSystem.

## Overview

The system now supports Google OAuth 2.0 authentication alongside the existing email/password authentication. Users can sign in with their Google accounts, and the system will automatically create or update their user profiles.

## Key Features

- **Dual Authentication**: Users can choose between email/password or Google Sign-In
- **Skip 2FA for Google**: Google authenticated users bypass the 2FA requirement
- **Auto User Creation**: New users are automatically created when signing in with Google
- **Provider Tracking**: Users are tracked with their authentication provider (`local` or `google`)
- **Consistent UI**: Google Sign-In button matches the pet adoption theme

## Implementation Details

### Frontend Configuration

**Location**: `frontend/src/app/config/google-auth.config.ts`

```typescript
const googleAuthConfig: GoogleAuthConfig = {
  clientId: 'YOUR_GOOGLE_CLIENT_ID',
  enabled: true
};
```

### Backend Configuration

**Location**: `backend/internal/services/backend_calls/user_services.go`

```go
// Google Client ID - in production, set via environment variables
clientID := "YOUR_GOOGLE_CLIENT_ID"
```

## Setup Instructions

### 1. Google Cloud Console Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the "Google+ API" and "Google Sign-In API"
4. Go to "APIs & Services" > "Credentials"
5. Create a new OAuth 2.0 Client ID
6. Add your domain to authorized origins:
   - `http://localhost:4200` (for development)
   - `https://yourdomain.com` (for production)

### 2. Frontend Configuration

Update `frontend/src/app/config/google-auth.config.ts`:

```typescript
const googleAuthConfig: GoogleAuthConfig = {
  clientId: 'YOUR_ACTUAL_GOOGLE_CLIENT_ID.apps.googleusercontent.com',
  enabled: true
};
```

### 3. Backend Configuration

Update the `verifyGoogleToken` function in `backend/internal/services/backend_calls/user_services.go`:

```go
// Option 1: Direct configuration
clientID := "YOUR_ACTUAL_GOOGLE_CLIENT_ID.apps.googleusercontent.com"

// Option 2: Environment variable (recommended for production)
clientID := os.Getenv("GOOGLE_CLIENT_ID")
if clientID == "" {
    clientID = "YOUR_FALLBACK_CLIENT_ID"
}
```

## Database Schema

The existing user table supports Google authentication through these fields:

```sql
-- Users table already includes:
Provider VARCHAR(255) DEFAULT 'local'    -- 'local' or 'google'
Provider_ID VARCHAR(255)                 -- Google user ID
Password VARCHAR(255)                    -- Empty for Google users
```

## Authentication Flow

### Google Sign-In Flow

1. User clicks "Continuar con Google" button
2. Google Sign-In popup appears
3. User authenticates with Google
4. Google returns ID token
5. Frontend sends token to backend `/api/auth/login/google`
6. Backend verifies token with Google's public keys
7. Backend creates/updates user with Google provider data
8. Backend generates session ID and returns user data
9. Frontend redirects to dashboard (skipping 2FA)

### Traditional Flow (unchanged)

1. User enters email/password
2. Backend validates credentials
3. Backend generates 2FA code and sends via email
4. User enters 2FA code
5. Backend verifies 2FA code
6. Frontend redirects to dashboard

## Security Considerations

1. **ID Token Verification**: All Google ID tokens are verified using Google's public keys
2. **Provider Isolation**: Google users cannot use password authentication
3. **Session Management**: Google users receive the same session handling as local users
4. **No 2FA Required**: Google authentication is considered secure enough to skip 2FA
5. **Account Linking**: Existing email accounts can be linked to Google accounts

## Testing

### Development Testing

1. Use a test Google Client ID for development
2. Add `http://localhost:4200` to authorized origins
3. Test with personal Google accounts

### Production Testing

1. Set production Google Client ID
2. Add production domain to authorized origins
3. Test with various Google account types

## Troubleshooting

### Common Issues

1. **"Google Sign-In no está disponible"**: 
   - Check if Google SDK is loaded
   - Verify Client ID configuration
   - Check browser console for errors

2. **Token verification fails**:
   - Verify Client ID matches in frontend and backend
   - Check if Google APIs are enabled in Google Cloud Console
   - Ensure authorized origins are correctly configured

3. **User creation fails**:
   - Check database connectivity
   - Verify user model fields support Google data
   - Check backend logs for specific errors

### Debug Mode

Enable debug logging in the backend:

```go
// In user_services.go, add logging:
log.Printf("Google token verification for email: %s", email)
log.Printf("Google user data: %+v", claims)
```

## File Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── config/
│   │   │   └── google-auth.config.ts     # Google auth configuration
│   │   ├── pages/
│   │   │   └── login/
│   │   │       ├── login.html           # Updated with Google button
│   │   │       └── login.ts             # Google auth logic
│   │   └── services/
│   │       └── api.service.ts           # API calls (already had Google endpoint)
│   └── index.html                       # Google SDK script tag

backend/
├── internal/
│   ├── services/
│   │   └── backend_calls/
│   │       └── user_services.go         # Google auth implementation
│   ├── db/
│   │   └── dao/
│   │       └── user.go                  # Updated for Google users
│   └── api/
│       └── routes/
│           └── user.go                  # Google auth endpoint (already existed)
└── go.mod                              # Added Google OAuth dependencies
```

## Dependencies

### Frontend
- Google Sign-In SDK (loaded via CDN)
- Angular (existing)

### Backend
- `google.golang.org/api/idtoken` - Google ID token verification
- `gorm.io/gorm` - Database ORM (existing)
- `github.com/labstack/echo/v4` - HTTP framework (existing)

## Future Enhancements

1. **Multiple OAuth Providers**: Add Facebook, GitHub, etc.
2. **Account Linking**: Allow users to link multiple auth providers
3. **Profile Sync**: Sync Google profile picture and other data
4. **Enhanced Security**: Add additional security checks for Google accounts
5. **Admin Dashboard**: Manage users by authentication provider