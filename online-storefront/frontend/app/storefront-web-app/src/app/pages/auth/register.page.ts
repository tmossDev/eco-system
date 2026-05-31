import { Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-register-page',
  imports: [ReactiveFormsModule, RouterLink],
  template: `
    <section class="auth">
      <form [formGroup]="form" (ngSubmit)="submit()">
        <p class="eyebrow">Join Northstar</p>
        <h1>Create account</h1>
        <p class="intro">Set up your account to keep your cart close and make future checkout smoother.</p>
        <div class="row">
          <label>First name <input type="text" formControlName="first_name" autocomplete="given-name"></label>
          <label>Last name <input type="text" formControlName="last_name" autocomplete="family-name"></label>
        </div>
        <label>Email <input type="email" formControlName="email" autocomplete="email"></label>
        <label>Password <input type="password" formControlName="password" autocomplete="new-password"></label>
        <label>Confirm password <input type="password" formControlName="confirm_password" autocomplete="new-password"></label>
        @if (errorMessage()) { <p class="error">{{ errorMessage() }}</p> }
        <button type="submit" [disabled]="form.invalid || isSubmitting()">
          {{ isSubmitting() ? 'Creating account...' : 'Create account' }}
        </button>
        <p class="switch">Already registered? <a routerLink="/auth/login">Sign in</a></p>
      </form>
    </section>
  `,
  styleUrl: './account-form.scss',
})
export class RegisterPage {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  protected readonly isSubmitting = signal(false);
  protected readonly errorMessage = signal('');
  protected readonly form = new FormGroup({
    first_name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    last_name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    email: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.email] }),
    password: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    confirm_password: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
  });

  protected submit(): void {
    if (this.form.invalid || this.isSubmitting()) return;
    const request = this.form.getRawValue();
    if (request.password !== request.confirm_password) {
      this.errorMessage.set('Passwords must match.');
      return;
    }
    this.isSubmitting.set(true);
    this.errorMessage.set('');
    this.authService.register(request)
      .pipe(finalize(() => this.isSubmitting.set(false)))
      .subscribe({
        next: () => void this.router.navigate(['/']),
        error: () => this.errorMessage.set('Unable to create your account. Please check your details and try again.'),
      });
  }
}
