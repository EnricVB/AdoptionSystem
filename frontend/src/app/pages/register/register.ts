import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';

@Component({
  selector: 'app-register',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './register.html',
  styleUrls: ['../../../styles/login.css', '../../../styles/register.css'],
  host: {
    'style': 'view-transition-name: auth-form'
  }
})
export class Register {

  // ======================================
  // COMPONENT PROPERTIES
  // ======================================
  registerForm!: FormGroup;
  submitted = false;
  error: string | null = null;
  success: string | null = null;

  // Focus state tracking for floating labels
  isFirstNameFocused = false;
  isLastNameFocused = false;
  isEmailFocused = false;
  isPasswordFocused = false;
  isConfirmPasswordFocused = false;

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

  // Initializes the registration form with validation rules.
  // This method is called in the constructor to set up the form controls.
  private initializeForm(): void {
    this.registerForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
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

  // ======================================
  // FORM SUBMISSION
  // ======================================

  onSubmit(): void {
    this.submitted = true;
    this.error = null;
    this.success = null;
    console.log('Registering user');

    const payload = this.buildRegisterPayload();

    this.apiService.register(payload).subscribe({
      next: (response) => this.onRegisterSuccess(response),
      error: (err) => this.onRegisterError(err)
    });
  }

  // ======================================
  // AUTH FLOW
  // ======================================

  /**
   * Builds the payload for the registration request.
   * This method extracts the form data and prepares it for the API call.
   * 
   * @returns An object containing the registration data.
   */
  private buildRegisterPayload(): any {
    const formValue = this.registerForm.value;
    return {
      Name: formValue.firstName,
      Surname: formValue.lastName,
      Email: formValue.email,
      Password: formValue.password
    };
  }

  /**
   * Success handler for the registration API call.
   * This method is called when the registration API returns a successful response.
   * 
   * @param response Any response from the registration API.
   */
  private onRegisterSuccess(response: any): void {
    this.success = '¡Registro exitoso! Redirigiendo al login...';
    
    setTimeout(() => {
      this.router.navigate(['/login']);
    }, 2000);
  }

  /**
   * Error handler for the registration API call.
   * This method is called when the registration API returns an error response.
   * 
   * @param error Any error object returned from the registration API.
   */
  private onRegisterError(error: any): void {
    this.submitted = false;
    this.error = error.error?.message || 'Error en el registro. Por favor intenta de nuevo.';
  }

  // ======================================
  // FLOATING LABEL METHODS
  // ======================================

  /**
   * Handle first name input focus event
   * Used for floating label animation
   */
  onFirstNameFocus(): void {
    this.isFirstNameFocused = true;
  }

  /**
   * Handle first name input blur event
   * Used for floating label animation
   */
  onFirstNameBlur(): void {
    this.isFirstNameFocused = false;
  }

  /**
   * Handle last name input focus event
   * Used for floating label animation
   */
  onLastNameFocus(): void {
    this.isLastNameFocused = true;
  }

  /**
   * Handle last name input blur event
   * Used for floating label animation
   */
  onLastNameBlur(): void {
    this.isLastNameFocused = false;
  }

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

  // ======================================
  // NAVIGATION METHODS
  // ======================================

  /**
   * Navigate to login page
   */
  goToLogin(): void {
    this.router.navigate(['/login']);
  }
}
