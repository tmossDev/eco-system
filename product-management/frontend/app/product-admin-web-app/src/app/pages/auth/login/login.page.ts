import { Component, CUSTOM_ELEMENTS_SCHEMA, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { finalize } from 'rxjs';

import { AuthService } from '../../../core/services/auth/auth.service';

@Component({
  selector: 'app-login-page',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  template: `
    <ds-login-form
      [isSubmitting]="isSubmitting()"
      [errorMessage]="errorMessage()"
      emailPlaceholder="admin@test.com"
      forgotPasswordHref="/auth/forgot-password"
      (loginSubmit)="login($event)"
    />
  `,
})
export class LoginPage {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal('');

  protected login(event: any): void {
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
