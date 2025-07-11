import { CommonModule } from '@angular/common';
import { Component, OnInit, NgZone } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms'
import { Router, RouterModule } from '@angular/router';
import { ApiService } from '../../services/api.service';
import googleAuthConfig from '../../config/google-auth.config';

declare const google: any;

@Component({
  selector: 'app-login',
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: './login.html',
  styleUrls: ['../../../styles/login.css'],
  host: {
    'style': 'view-transition-name: auth-form'
  }
})
export class Login implements OnInit {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  loginForm!: FormGroup;
  submitted = false;
  googleSubmitted = false;
  error: string | null = null;

  // Focus state tracking for floating labels
  isEmailFocused = false;
  isPasswordFocused = false;

  // Show password toggle
  showPassword = false;

  // Google Auth configuration
  googleAuthEnabled = googleAuthConfig.enabled;

  // ======================================
  // CONSTRUCTOR
  // ======================================

  constructor(
    private fb: FormBuilder, 
    private apiService: ApiService,
    private router: Router,
    private ngZone: NgZone
  ) {
    this.initializeForm();
  }

  ngOnInit(): void {
    this.initializeGoogleSignIn();
  }

  // Initializes the login form with validation rules.
  // This method is called in the constructor to set up the form controls.
  private initializeForm(): void {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required]
    });
  }

  // ======================================
  // FORM SUBMISSION
  // ======================================

  /**
   * Handles the form submission for the login.
   * This method is called when the user submits the login form.
   * 
   * It sets the submitted flag to true, clears any previous error messages,
   * builds the login payload from the form values, and calls the API service to perform the
   * login operation.
   * 
   * @returns {void}
   */
  onSubmit(): void {
    this.submitted = true;
    this.error = null;

    const payload = this.buildLoginPayload();

    this.apiService.login(payload).subscribe({
      next: (response) => this.onLoginSuccess(response),
      error: (err) => this.onLoginError(err)
    });
  }

  // ======================================
  // AUTH FLOW
  // ======================================

  /**
   * Builds the payload for the login request.
   * This method extracts the email and password from the login form
   * 
   * @returns An object containing the email and password from the login form.
   */
  private buildLoginPayload(): { email: string; password: string } {
    this.loginForm.setValue({
      email: this.loginForm.value.email.trimEnd(),
      password: this.loginForm.value.password.trimEnd()
    });

    return {
      email: this.loginForm.value.email,
      password: this.loginForm.value.password
    };
  }

  /**
   * Success handler for the login API call.
   * This method is called when the login API returns a successful response.
   * 
   * @param response Any response from the login API.
   */
  private onLoginSuccess(response: any): void {
    const userID = response.content.id;
    const changePass = response.content.change_pass;
    const sessionID = response.content.session_id;

    if (changePass) {
      this.router.navigate(['/change-pass'], { state: { email: this.loginForm.value.email, userID: userID }});
      return;
    }

    this.send2FAEmail(this.loginForm.value.email);
    this.router.navigate(['/twofa'], {state: { email: this.loginForm.value.email, sessionID: sessionID }});
  }

  private send2FAEmail(email: string): void {
    const EmailPayload = { email: email };

    this.apiService.refresh2FAToken(EmailPayload).subscribe({
      next: (response) => {},
      error: (err) => {}
    });
  }

  /**
   * Error handler for the login API call.
   * This method is called when the login API returns an error response.
   * 
   * @param error Any error object returned from the login API.
   */
  private onLoginError(error: any): void {
    this.submitted = false;
    this.error = error.error?.message || 'Login failed. Please check your credentials.';
  }

  // ======================================
  // FLOATING LABEL METHODS
  // ======================================

  /**
   * Handle email input focus event
   * Used for floating label animation
   */
  onEmailFocus(): void {
    this.isEmailFocused = true;
  }

  /**
   * Handle email input blur event
   * Used for floating label animation
   */
  onEmailBlur(): void {
    this.isEmailFocused = false;
  }

  /**
   * Handle password input focus event
   * Used for floating label animation
   */
  onPasswordFocus(): void {
    this.isPasswordFocused = true;
  }

  /**
   * Handle password input blur event
   * Used for floating label animation
   */
  onPasswordBlur(): void {
    this.isPasswordFocused = false;
  }

  // ======================================
  // GOOGLE SIGN-IN METHODS
  // ======================================

  /**
   * Initialize Google Sign-In
   */
  private initializeGoogleSignIn(): void {
    if (!googleAuthConfig.enabled) {
      console.warn('Google Sign-In is disabled in configuration');
      return;
    }

    if (typeof google !== 'undefined') {
      google.accounts.id.initialize({
        client_id: googleAuthConfig.clientId,
        callback: (response: any) => this.handleGoogleSignInResponse(response),
        auto_select: false,
        cancel_on_tap_outside: true
      });
    } else {
      console.error('Google Sign-In SDK not loaded');
    }
  }

  /**
   * Handle Google Sign-In button click
   */
  signInWithGoogle(): void {
    if (!googleAuthConfig.enabled) {
      this.error = 'Google Sign-In está deshabilitado en la configuración';
      return;
    }

    this.googleSubmitted = true;
    this.error = null;
    
    if (typeof google !== 'undefined') {
      google.accounts.id.prompt((notification: any) => {
        if (notification.isNotDisplayed() || notification.isSkippedMoment()) {
          google.accounts.id.prompt();
        }
      });
    } else {
      this.error = 'Google Sign-In no está disponible. Por favor, recarga la página.';
      this.googleSubmitted = false;
    }
  }

  /**
   * Handle Google Sign-In response
   */
  private handleGoogleSignInResponse(response: any): void {
    if (response.credential) {
      // Decode the JWT token to get user info
      const payload = this.decodeJwt(response.credential);
      
      const googleLoginData = {
        email: payload.email,
        id_token: response.credential
      };

      this.apiService.loginWithGoogle(googleLoginData).subscribe({
        next: (apiResponse) => this.onGoogleLoginSuccess(apiResponse),
        error: (err) => this.onGoogleLoginError(err)
      });
    } else {
      this.error = 'Error en la autenticación con Google';
      this.submitted = false;
    }
  }

  /**
   * Success handler for Google login
   */
  private onGoogleLoginSuccess(response: any): void {
    this.ngZone.run(() => {
      // Google authentication successful, redirect to dashboard
      // Skip 2FA for Google users as specified in requirements
      this.router.navigate(['/dashboard']);
    });
  }

  /**
   * Error handler for Google login
   */
  private onGoogleLoginError(error: any): void {
    this.ngZone.run(() => {
      this.submitted = false;
      this.error = error.error?.message || 'Error en la autenticación con Google. Por favor, inténtalo de nuevo.';
    });
  }

  /**
   * Decode JWT token to extract user information
   */
  private decodeJwt(token: string): any {
    try {
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));
      return JSON.parse(jsonPayload);
    } catch (error) {
      console.error('Error decoding JWT:', error);
      return {};
    }
  }
}
