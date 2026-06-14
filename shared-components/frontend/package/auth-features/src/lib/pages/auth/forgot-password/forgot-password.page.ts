import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';

import { AuthService } from '../../../core/services/auth/auth.service';

@Component({
  selector: 'app-forgot-password-page',
  imports: [FormsModule, RouterLink],
  template: `
    <form class="forgot-password-form" (ngSubmit)="sendResetLink()">
      <div class="form-header">
        <p class="eyebrow">Password recovery</p>
        <h2>Forgot your password?</h2>
        <p>
          Enter your email address and we’ll send you instructions to reset your
          password.
        </p>
      </div>

      @if (message) {
        <p class="success-message">{{ message }}</p>
      }

      @if (errorMessage) {
        <p class="error-message">{{ errorMessage }}</p>
      }

      <label>
        <span>Email address</span>
        <input
          type="email"
          name="email"
          autocomplete="email"
          placeholder="you@example.com"
          [(ngModel)]="email"
        />
      </label>

      <button type="submit" [disabled]="isSubmitting">
        {{ isSubmitting ? 'Sending...' : 'Send reset link' }}
      </button>

      <a routerLink="/auth/login" class="back-link">Back to login</a>
    </form>
  `,
  styles: `
    .forgot-password-form {
      display: grid;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 1.25rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 20px 60px rgb(15 23 42 / 10%);
    }

    .form-header {
      margin-bottom: 0.5rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #2563eb;
      font-size: 0.8rem;
      font-weight: 800;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h2 {
      margin: 0;
      color: #172033;
      font-size: 2rem;
      letter-spacing: -0.04em;
    }

    .form-header p {
      margin: 0.75rem 0 0;
      color: #64748b;
      line-height: 1.6;
    }

    .success-message,
    .error-message {
      margin: 0;
      border-radius: 0.75rem;
      padding: 0.75rem 0.9rem;
      font-weight: 700;
    }

    .success-message {
      background: #dcfce7;
      color: #166534;
    }

    .error-message {
      background: #fee2e2;
      color: #991b1b;
    }

    label {
      display: grid;
      gap: 0.4rem;
      color: #172033;
      font-weight: 700;
    }

    input {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.8rem 0.9rem;
      color: #172033;
      font: inherit;
    }

    input:focus {
      border-color: #2563eb;
      outline: 3px solid #dbeafe;
    }

    button {
      border: 0;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.85rem 1rem;
      font: inherit;
      font-weight: 800;
      cursor: pointer;
    }

    button:disabled {
      cursor: default;
      opacity: 0.65;
    }

    .back-link {
      justify-self: center;
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
    }
  `,
})
export class ForgotPasswordPage {
  private readonly authService = inject(AuthService);

  protected email = '';
  protected isSubmitting = false;
  protected message = '';
  protected errorMessage = '';

  protected sendResetLink(): void {
    if (this.isSubmitting) {
      return;
    }

    this.isSubmitting = true;
    this.message = '';
    this.errorMessage = '';

    this.authService.forgotPassword({ email: this.email }).subscribe({
      next: () => {
        this.message = 'If an account exists for this email, reset instructions will be sent.';
        this.isSubmitting = false;
      },
      error: () => {
        this.errorMessage = 'Unable to send reset instructions. Please try again.';
        this.isSubmitting = false;
      },
    });
  }
}
