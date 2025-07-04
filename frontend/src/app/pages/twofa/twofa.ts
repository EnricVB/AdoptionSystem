import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-twofa',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './twofa.html',
  styleUrl: './twofa.css'
})
export class Twofa {
codeForm!: FormGroup;
  error: string | null = null;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {
    this.codeForm = this.fb.group({
      code: ['', [Validators.required, Validators.pattern(/^[A-Z0-9]{6}$/)]]
    });

    this.codeForm.get('code')?.valueChanges.subscribe(value => {
    if (value && value !== value.toUpperCase()) {
      this.codeForm.get('code')?.setValue(value.toUpperCase(), { emitEvent: false });
    }
  })
  }

  onSubmit() {
    this.error = null;

    if (this.codeForm.invalid) return;

    const payload = { code: this.codeForm.value.code };

    this.http.post('/api/auth/verify-2fa', payload, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    }).subscribe({
      next: () => this.router.navigate(['/dashboard']),
      error: (err) => {
        this.error = err.error?.message || 'Invalid 2FA code.';
      }
    });
  }
}
