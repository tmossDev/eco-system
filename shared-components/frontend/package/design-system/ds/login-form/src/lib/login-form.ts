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
  styleUrl: './login-form.scss',
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
