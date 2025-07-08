import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms'
import { Router, RouterModule } from '@angular/router';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-login',
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: './login.html',
  styleUrls: ['../../../styles/login.css'],
  host: {
    'style': 'view-transition-name: auth-form'
  }
})
export class Login {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  loginForm!: FormGroup;
  submitted = false;
  error: string | null = null;

  // Focus state tracking for floating labels
  isEmailFocused = false;
  isPasswordFocused = false;

  // Show password toggle
  showPassword = false;

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
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required]
    });
  }

  // ======================================
  // FORM SUBMISSION
  // ======================================

  onSubmit(): void {
    this.submitted = true;
    this.error = null;

    const payload = this.buildLoginPayload();
    console.log(this.loginForm.value.email);

    this.apiService.login(payload).subscribe({
      next: (response) => this.onLoginSuccess(response),
      error: (err) => this.onLoginError(err)
    });
  }

  sendMail(email: string): void {
    this.apiService.refresh2FAToken({ email }).subscribe({
      error: (err) => this.onSendMailError(err)
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
    const sessionID = response.content.session_id;
    
    this.sendMail(this.loginForm.value.email);
    this.router.navigate(['/twofa'], {state: {sessionID}});
  }

  /**
   * Error handler for the login API call.
   * This method is called when the login API returns an error response.
   * 
   * @param error Any error object returned from the login API.
   */
  private onLoginError(error: any): void {
    setTimeout(() => {
      this.submitted = false;
    }, 2000);

    this.error = error.error?.message || 'Login failed. Please check your credentials.';
  }

  /**
   * Error handler for the sendMail API call
   * This method is called when the sendMail API returns an error response.
   * 
   * @param error Any error object returned from the mailer API.
   */
  private onSendMailError(error: any): void {
    this.error = error.error?.message || 'Failed to send email. Please try again later.';

    // Wait for 3 seconds before allowing another submission
    setTimeout(() => {
      this.submitted = false;
    }, 2000);
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
}
