import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { ApiService } from '../../../services/api.service';
import { CookieService } from '@app/services/cookie.service';

@Component({
  selector: 'app-change-pass',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './change-pass.html',
  styleUrls: ['../../../../styles/login.css', '../../../../styles/register.css'],
  host: {
    'style': 'view-transition-name: auth-form'
  }
})
export class ChangePass {
  
  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  changePassForm!: FormGroup;
  submitted = false;
  error: string | null = null;
  success: string | null = null;

  // Show password toggle
  showPassword = false;

  // Focus state tracking for floating labels
  isPasswordFocused = false;
  isConfirmPasswordFocused = false;

  // User information from query parameters
  email!: string;
  userID!: string;

  // ======================================
  // CONSTRUCTOR
  // ======================================
  constructor(
    private fb: FormBuilder,
    private apiService: ApiService,
    private router: Router,
    private route: ActivatedRoute,
    private cookieService: CookieService
  ) {
    this.initializeForm();
  }

  // Initializes the registration form with validation rules.
  // This method is called in the constructor to set up the form controls.
  private initializeForm(): void {
    this.changePassForm = this.fb.group({
      password: ['', [Validators.required, Validators.minLength(8)]],
      confirmPassword: ['', [Validators.required]]
    }, { validators: this.passwordMatchValidator });
  }

  // Custom validator to check if passwords match
  private passwordMatchValidator(form: FormGroup) {
    const password = form.get('password')?.value;
    const confirmPassword = form.get('confirmPassword')?.value;

    if (!password || !confirmPassword) return null;
    
    return password !== confirmPassword ? { passwordMismatch: true } : null;
  }

  /**
   * Lifecycle hook that is called after the component has been initialized.
   * It subscribes to the route's query parameters to retrieve the session ID, user ID, and email.
   * 
   * This allows the component to access these parameters for further processing, such as
   * displaying user-specific information or handling password change requests.
   * 
   * @returns {void}
   */
  ngOnInit(): void {
    this.email = history.state.email;
    this.userID = history.state.userID;

    if (!this.email || !this.userID) {
      this.error = 'Email and User ID are required to change password.';
      return;
    }
  }
 

  // ======================================
  // FORM SUBMISSION
  // ======================================

  /**
   * Handles the form submission for the change password.
   * This method is called when the user submits the change password form.
   * 
   * It builds the payload from the form data and calls the API service to change the password.
   * If the submission is successful, it redirects the user to the dashboard page.
   * If there is an error, it sets the error message to be displayed.
   * 
   * @returns {void}
   */
  onSubmit(): void {
    this.submitted = true;
    this.error = null;

    const payload = this.buildChangePassPayload();

    this.apiService.changePassword(payload).subscribe({
      next: (response) => this.onChangePassSuccess(response),
      error: (err) => this.onChangePassError(err)
    });
  }

  // ======================================
  // AUTH FLOW
  // ======================================

  /**
   * Builds the payload for the change pass request.
   * This method extracts the password from the change pass form
   * 
   * @returns An object containing the password from the change pass form.
   */
  private buildChangePassPayload(): { email: string; password: string } {
    const password = this.changePassForm.value.password?.trimEnd() || '';

    return {
      email: this.email,
      password: password
    };
  }

  /**
   * Success handler for the change pass API call.
   * This method is called when the change pass API returns a successful response.
   * 
   * @param response Any response from the change pass API.
   */
  private onChangePassSuccess(response: any): void {
    this.success = 'Password changed successfully!';

    setTimeout(() => {
      this.cookieService.setCookie('sessionID', response.content.session_id, 1);
      this.router.navigate(['/dashboard'], { state: { email: this.email }});
    });
  }

  /**
   * Error handler for the login API call.
   * This method is called when the login API returns an error response.
   * 
   * @param error Any error object returned from the login API.
   */
  private onChangePassError(error: any): void {
    this.submitted = false;
    this.error = error.error?.message || 'Change Password failed. Please check your password.';
  }

  // ======================================
  // VFX
  // ======================================

  
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

  /**
   * Handle confirm password input focus event
   * Used for floating label animation
   */
  onConfirmPasswordFocus(): void {
    this.isConfirmPasswordFocused = true;
  }

  /**
   * Handle confirm password input blur event
   * Used for floating label animation
   */
  onConfirmPasswordBlur(): void {
    this.isConfirmPasswordFocused = false;
  }
}