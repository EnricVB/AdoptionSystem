import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms'
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

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

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {
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

    this.http.post<any>(`/api/auth/login`, payload, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      })
    }).subscribe({
      next: (response) => {
        console.log('Login successful', response);
        const sessionID = response.content.session_id;
        this.sendmail(this.loginForm.value.email);
        this.router.navigate(['/twofa'], {state: {sessionID}});
      },
      error: (err) => {
        console.error('Login failed:', err);
        this.error = err.error?.message || 'Login failed. Please check your credentials.';
      }
    });
  }

 sendmail(email: string) {
  const payload = { email };

  this.http.post<any>('/api/auth/refresh-token', payload, {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  }).subscribe({
    next: (response) => {
      console.log('Email sent successfully', response);
    },
    error: (error) => {
      console.error('Error sending email:', error);
      this.error = error.error?.message || 'Failed to send email. Please try again later.';
    }
  });
}
}
