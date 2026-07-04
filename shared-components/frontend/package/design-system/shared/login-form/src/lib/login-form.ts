import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ButtonComponent } from '@primitive/button';
import { DsCheckboxDirective, DsFormField, DsInputDirective } from '@primitive/form';

export interface LoginFormSubmitEvent {
  email: string;
  password: string;
  rememberMe: boolean;
}

@Component({
  selector: 'shared-login-form',
  imports: [FormsModule, ButtonComponent, DsCheckboxDirective, DsFormField, DsInputDirective],
  template: `
    <form class="shared-login-form" (ngSubmit)="submitLogin()">
      <div class="shared-login-form__header">
        <p class="shared-login-form__eyebrow">{{ eyebrow }}</p>
        <h2>{{ heading }}</h2>
        <p>{{ description }}</p>
      </div>

      @if (errorMessage) {
        <p class="shared-login-form__error-message">{{ errorMessage }}</p>
      }

      <primitive-form-field label="Email address">
        <input
          primitiveInput
          type="email"
          name="email"
          autocomplete="email"
          [placeholder]="emailPlaceholder"
          [(ngModel)]="form.email"
        />
      </primitive-form-field>

      <primitive-form-field label="Password">
        <input
          primitiveInput
          type="password"
          name="password"
          autocomplete="current-password"
          [placeholder]="passwordPlaceholder"
          [(ngModel)]="form.password"
        />
      </primitive-form-field>

      <div class="shared-login-form__row">
        <label class="shared-login-form__remember-me">
          <input
            primitiveCheckbox
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

      <primitive-button
        type="submit"
        variant="primary"
        size="large"
        [disabled]="isSubmitting"
        [label]="isSubmitting ? submittingText : submitText"
      />
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
