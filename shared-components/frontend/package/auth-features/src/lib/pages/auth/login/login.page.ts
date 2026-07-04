import { Component, inject, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { LoginForm, LoginFormSubmitEvent } from '@shared/login-form';
import { finalize } from 'rxjs';

import { AuthService } from '../../../core/services/auth/auth.service';

@Component({
  selector: 'app-login-page',
  imports: [LoginForm],
  template: `
    <shared-login-form
      [isSubmitting]="isSubmitting()"
      [errorMessage]="errorMessage()"
      [emailPlaceholder]="emailPlaceholder"
      forgotPasswordHref="/auth/forgot-password"
      (loginSubmit)="login($event)"
    />
  `,
})
export class LoginPage {
  private readonly authService = inject(AuthService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  protected readonly emailPlaceholder =
    this.route.snapshot.data['loginEmailPlaceholder'] ?? 'admin@example.com';
  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal('');

  protected login(event: LoginFormSubmitEvent): void {
    if (this.isSubmitting()) {
      return;
    }

    this.isSubmitting.set(true);
    this.errorMessage.set('');

    this.authService
      .login(event)
      .pipe(
        finalize(() => {
          this.isSubmitting.set(false);
        }),
      )
      .subscribe({
        next: () => {
          void this.router.navigate(['/']);
        },
        error: () => {
          this.errorMessage.set('Unable to sign in. Please check your details and try again.');
        },
      });
  }
}
