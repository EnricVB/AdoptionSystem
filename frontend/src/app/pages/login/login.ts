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

    this.router.navigate(['/twofa'], {state: {email: this.loginForm.value.email, sessionID: sessionID }});
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
