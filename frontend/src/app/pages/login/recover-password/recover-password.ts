import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms'
import { Router, RouterModule } from '@angular/router';
import { ApiService } from '../../../services/api.service';

@Component({
  selector: 'app-recover-password',
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
  ],
  templateUrl: './recover-password.html',
})
export class RecoverPassword {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  recoverPasswordForm!: FormGroup;
  submitted = false;
  error: string | null = null;
  success: string | null = null;

  // Focus state tracking for floating labels
  isEmailFocused = false;

  // ======================================
  // CONSTRUCTOR
  // ======================================

  constructor(
    private fb: FormBuilder, 
    private apiService: ApiService,
    private router: Router
  ) {
    this.initializeForm();
  }

  // Initializes the login form with validation rules.
  // This method is called in the constructor to set up the form controls.
  private initializeForm(): void {
    this.recoverPasswordForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]]
    });
  }

  // ======================================
  // FORM SUBMISSION
  // ======================================

  onSubmit(): void {
    this.submitted = true;
    this.error = null;

    const payload = this.buildPayload();

    this.apiService.sendRecoverPasswordMail(payload).subscribe({
      next: (response) => this.onRecoverPasswordSuccess(response),
      error: (err) => this.onSendMailError(err)
    });
  }

  /**
   * Builds the payload for the recover password API call.
   * @returns An object containing the email from the form.
   */
  private buildPayload(): { email: string } {
    return {
      email: this.recoverPasswordForm.value.email
    };
  }
  
  // ======================================
  // RECOVER PASSWORD FLOW
  // ======================================

  /**
   * Success handler for the recover password API call.
   * This method is called when the recover password API returns a successful response.
   * 
   * @param response Any response from the recover password API.
   */
  private onRecoverPasswordSuccess(response: any): void {
    const sessionID = response.content.session_id;
    this.success = 'An email has been sent to your address with instructions to reset your password. Redirecting to login...';
    
    setTimeout(() => {
      this.router.navigate(['/login'], {state: {sessionID}});
    }, 2000);
  }

  /**
   * Error handler for the sendMail API call
   * This method is called when the sendMail API returns an error response.
   * 
   * @param error Any error object returned from the mailer API.
   */
  private onSendMailError(error: any): void {
    this.submitted = false;
    this.error = error.error?.message || 'Failed to send email. Please try again later.';
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
}