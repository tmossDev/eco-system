import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';

export interface LoginFormSubmitEvent {
  email: string;
  password: string;
  rememberMe: boolean;
}

@Component({
  selector: 'ds-login-form',
  imports: [FormsModule],
  template: `
    <form class="ds-login-form" (ngSubmit)="submitLogin()">
      <div class="ds-login-form__header">
        <p class="ds-login-form__eyebrow">{{ eyebrow }}</p>
        <h2>{{ heading }}</h2>
        <p>{{ description }}</p>
      </div>

      @if (errorMessage) {
        <p class="ds-login-form__error-message">{{ errorMessage }}</p>
      }

      <label>
        <span>Email address</span>
        <input
          type="email"
          name="email"
          autocomplete="email"
          [placeholder]="emailPlaceholder"
          [(ngModel)]="form.email"
        />
      </label>

      <label>
        <span>Password</span>
        <input
          type="password"
          name="password"
          autocomplete="current-password"
          [placeholder]="passwordPlaceholder"
          [(ngModel)]="form.password"
        />
      </label>

      <div class="ds-login-form__row">
        <label class="ds-login-form__remember-me">
          <input
            type="checkbox"
            name="rememberMe"
            [(ngModel)]="form.rememberMe"
          />
          <span>Remember me</span>
        </label>

        @if (showForgotPasswordLink) {
          <a [href]="forgotPasswordHref">Forgot password?</a>
        }
      </div>

      <button type="submit" [disabled]="isSubmitting">
        {{ isSubmitting ? submittingText : submitText }}
      </button>
    </form>
  `,
  styles: `
    .ds-login-form {
      display: grid;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 1.25rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 20px 60px rgb(15 23 42 / 10%);
    }

    .ds-login-form__header {
      margin-bottom: 0.5rem;
    }

    .ds-login-form__eyebrow {
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

    .ds-login-form__header p:last-child {
      margin: 0.75rem 0 0;
      color: #64748b;
      line-height: 1.6;
    }

    .ds-login-form__error-message {
      margin: 0;
      border-radius: 0.75rem;
      background: #fee2e2;
      color: #991b1b;
      padding: 0.75rem 0.9rem;
      font-weight: 700;
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

    .ds-login-form__row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
    }

    .ds-login-form__remember-me {
      display: flex;
      grid-template-columns: auto 1fr;
      align-items: center;
      gap: 0.45rem;
      color: #64748b;
      font-weight: 600;
    }

    .ds-login-form__remember-me input {
      width: 1rem;
      height: 1rem;
    }

    a {
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
      white-space: nowrap;
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

    @media (max-width: 420px) {
      .ds-login-form__row {
        align-items: flex-start;
        flex-direction: column;
      }
    }
  `,
})
export class LoginForm {
  @Input()
  public eyebrow = 'Welcome back';

  @Input()
  public heading = 'Sign in to your account';

  @Input()
  public description = 'Enter your details below to access the admin dashboard.';

  @Input()
  public emailPlaceholder = 'you@example.com';

  @Input()
  public passwordPlaceholder = 'Enter your password';

  @Input()
  public submitText = 'Sign in';

  @Input()
  public submittingText = 'Signing in...';

  @Input()
  public isSubmitting = false;

  @Input()
  public errorMessage = '';

  @Input()
  public showForgotPasswordLink = true;

  @Input()
  public forgotPasswordHref = '/auth/forgot-password';

  @Output()
  public loginSubmit = new EventEmitter<LoginFormSubmitEvent>();

  protected readonly form: LoginFormSubmitEvent = {
    email: '',
    password: '',
    rememberMe: true,
  };

  protected submitLogin(): void {
    if (this.isSubmitting) {
      return;
    }

    this.loginSubmit.emit({
      ...this.form,
    });
  }
}
