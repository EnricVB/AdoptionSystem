import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { CookieService } from '@app/services/cookie.service';

@Component({
  selector: 'app-twofa',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './twofa.html',
  styleUrls: ['../../../styles/twofa.css', '../../../styles/transitions/2fa-transition.css'],
  host: {
    'style': 'view-transition-name: twofa-page'
  }
})
export class Twofa {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  codeForm!: FormGroup;
  submitted = false;
  error: string | null = null;
  sessionID: string = '';
  email: string = '';

  // Focus state tracking for floating labels
  isCodeFocused = false;

  // Resend 2FA code cooldown
  resendCooldown = 0;
  resendInterval: any;

  // ======================================
  // CONSTRUCTOR
  // ======================================

  constructor(
    private fb: FormBuilder, 
    private apiService: ApiService, 
    private router: Router,
    private cookieService: CookieService
  ) {
    this.initializeSessionAndForm();
  }

  // Initializes the session from router state and sets up the 2FA form with validation rules.
  // This method is called in the constructor to set up the session ID and form controls.
  private initializeSessionAndForm(): void {
    const navigation = this.router.getCurrentNavigation();
    const state = navigation?.extras.state as { sessionID?: string, email?: string };
    this.sessionID = state?.sessionID || '';
    this.email = state?.email || '';

    this.codeForm = this.fb.group({
      code: ['', [Validators.required, Validators.pattern(/^[A-Z0-9]{6}$/)]]
    });

    // Auto-convert input to uppercase
    this.codeForm.get('code')?.valueChanges.subscribe(value => {
      if (value && value !== value.toUpperCase()) {
        this.codeForm.get('code')?.setValue(value.toUpperCase(), { emitEvent: false });
      }
    });
  }

  // ======================================
  // FORM SUBMISSION
  // ======================================
  onResend2FA(): void {
    if (this.resendCooldown > 0) return;
    const payload = this.build2FAEmailPayload();

    this.startResendCooldown();
    this.apiService.refresh2FAToken(payload).subscribe({
      next: () => {},
      error: (err) => {}
    });
  }

  startResendCooldown(): void {
    this.resendCooldown = 30;
    this.resendInterval = setInterval(() => {
      this.resendCooldown--;
      if (this.resendCooldown <= 0) {
        clearInterval(this.resendInterval);
      }
    }, 1000);
  }

  onSubmit(): void {
    this.submitted = true;
    this.error = null;

    if (this.codeForm.invalid) {
      this.submitted = false;
      return;
    }

    const payload = this.build2FAPayload();

    this.apiService.verify2FA(payload).subscribe({
      next: (response) => this.on2FASuccess(response),
      error: (err) => this.on2FAError(err)
    });
  }

  // ======================================
  // AUTH FLOW
  // ======================================

  /**
   * Builds the payload for the 2FA email request.
   * This method extracts the email from the session and prepares it for the API call.
   * 
   * @returns An object containing the email for the 2FA request.
   */
  private build2FAEmailPayload(): { email: string; } {
    return {
      email: this.email
    };
  }

  /**
   * Builds the payload for the 2FA verification request.
   * This method extracts the code from the form and combines it with the session ID.
   * 
   * @returns An object containing the verification code and session ID.
   */
  private build2FAPayload(): { code: string; session_id: string } {
    return {
      session_id: this.sessionID,
      code: this.codeForm.value.code
    };
  }

  /**
   * Success handler for the 2FA verification API call.
   * This method is called when the 2FA API returns a successful response.
   * 
   * @param response Any response from the 2FA verification API.
   */
  private on2FASuccess(response: any): void {
    this.cookieService.setCookie('sessionID', this.sessionID, 7);
    
    this.router.navigate(['/dashboard']);
  }

  /**
   * Error handler for the 2FA verification API call.
   * This method is called when the 2FA API returns an error response.
   * 
   * @param error Any error object returned from the 2FA verification API.
   */
  private on2FAError(error: any): void {
    this.submitted = false;
    this.error = error.error?.message || 'Invalid 2FA code. Please try again.';
  }

  // ======================================
  // FLOATING LABEL METHODS
  // ======================================

  /**
   * Handle code input focus event
   * Used for floating label animation
   */
  onCodeFocus(): void {
    this.isCodeFocused = true;
  }

  /**
   * Handle code input blur event
   * Used for floating label animation
   */
  onCodeBlur(): void {
    this.isCodeFocused = false;
  }
}
