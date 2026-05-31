import { Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-login-page',
  imports: [ReactiveFormsModule, RouterLink],
  template: `
    <section class="auth">
      <form [formGroup]="form" (ngSubmit)="submit()">
        <p class="eyebrow">Welcome back</p>
        <h1>Sign in</h1>
        <p class="intro">Pick up where you left off and keep your account close at hand.</p>
        <label>Email <input type="email" formControlName="email" autocomplete="email"></label>
        <label>Password <input type="password" formControlName="password" autocomplete="current-password"></label>
        @if (errorMessage()) { <p class="error">{{ errorMessage() }}</p> }
        <button type="submit" [disabled]="form.invalid || isSubmitting()">
          {{ isSubmitting() ? 'Signing in...' : 'Sign in' }}
        </button>
        <p class="switch">New here? <a routerLink="/auth/register">Create an account</a></p>
      </form>
    </section>
  `,
  styleUrl: './account-form.scss',
})
export class LoginPage {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal('');
  protected readonly form = new FormGroup({
    email: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.email] }),
    password: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
  });

  protected submit(): void {
    if (this.form.invalid || this.isSubmitting()) return;
    this.isSubmitting.set(true);
    this.errorMessage.set('');
    this.authService.login(this.form.getRawValue())
      .pipe(finalize(() => this.isSubmitting.set(false)))
      .subscribe({
        next: () => void this.router.navigate(['/']),
        error: () => this.errorMessage.set('Unable to sign in. Please check your details and try again.'),
      });
  }
}
