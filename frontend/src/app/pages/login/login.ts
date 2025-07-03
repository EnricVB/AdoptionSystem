import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms'
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-login',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './login.html',
  styleUrl: './login.css'
})
export class Login {
  loginForm: FormGroup;
  submitted = false;
  error: string | null = null;

  constructor(private fb: FormBuilder, private http: HttpClient) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required]
    });
  }

  onSubmit() {
    this.submitted = true;
    this.error = null;

    // if (!this.loginForm.valid) return;

    const payload = {
      email: this.loginForm.value.email,
      password: this.loginForm.value.password
    };

    this.http.post(`/api/login`, payload, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      })
    }).subscribe({
      next: (response) => {
        console.log('Login successful', response);
      },
      error: (err) => {
        console.error('Login failed:', err);
        this.error = err.error?.message || 'Login failed. Please check your credentials.';
      }
    });
  }
}
