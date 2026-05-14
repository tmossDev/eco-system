import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { UpdateUserRequest } from '../../../core/services/user/user.model';
import { UserService } from '../../../core/services/user/user.service';

@Component({
  selector: 'app-user-edit-page',
  imports: [FormsModule, RouterLink],
  template: `
    <section class="user-edit-page">
      <a [routerLink]="['/users', userId()]" class="back-link">
        ← Back to profile
      </a>

      @if (isLoading()) {
        <p class="state-message">Loading user...</p>
      }

      @if (message()) {
        <p class="success-message">{{ message() }}</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <form class="form-card" (ngSubmit)="saveUser()">
        <div class="form-header">
          <p class="eyebrow">Edit user</p>
          <h1>User details</h1>
          <p>
            Update the user's profile information and account access level.
          </p>
        </div>

        <label>
          <span>Name</span>
          <input
            type="text"
            name="name"
            [ngModel]="form().name"
            (ngModelChange)="updateForm('name', $event)"
          />
        </label>

        <label>
          <span>Email</span>
          <input
            type="email"
            name="email"
            [ngModel]="form().email"
            (ngModelChange)="updateForm('email', $event)"
          />
        </label>

        <label>
          <span>Role</span>
          <select
            name="role"
            [ngModel]="form().role"
            (ngModelChange)="updateForm('role', $event)"
          >
            <option>Admin</option>
            <option>Manager</option>
            <option>User</option>
          </select>
        </label>

        <label>
          <span>Status</span>
          <select
            name="status"
            [ngModel]="form().status"
            (ngModelChange)="updateForm('status', $event)"
          >
            <option>Active</option>
            <option>Pending</option>
            <option>Suspended</option>
          </select>
        </label>

        <div class="form-actions">
          <button type="submit" class="primary-action" [disabled]="isSaving()">
            {{ isSaving() ? 'Saving...' : 'Save changes' }}
          </button>
          <a [routerLink]="['/users', userId()]" class="secondary-action">
            Cancel
          </a>
        </div>
      </form>
    </section>
  `,
  styles: `
    .user-edit-page {
      padding: 2rem;
      color: #172033;
    }

    .back-link {
      display: inline-flex;
      margin-bottom: 1rem;
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
    }

    .state-message,
    .success-message,
    .error-message {
      margin: 0 0 1rem;
      border-radius: 0.75rem;
      padding: 0.75rem 0.9rem;
      font-weight: 700;
    }

    .state-message {
      background: #eff6ff;
      color: #1d4ed8;
    }

    .success-message {
      background: #dcfce7;
      color: #166534;
    }

    .error-message {
      background: #fee2e2;
      color: #991b1b;
    }

    .form-card {
      display: grid;
      gap: 1rem;
      max-width: 42rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .form-header {
      margin-bottom: 0.5rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #56657f;
      font-size: 0.8rem;
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h1 {
      margin: 0;
      font-size: clamp(2rem, 4vw, 3rem);
      letter-spacing: -0.04em;
    }

    .form-header p {
      margin: 0.75rem 0 0;
      color: #56657f;
      line-height: 1.6;
    }

    label {
      display: grid;
      gap: 0.4rem;
      font-weight: 700;
    }

    input,
    select {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.8rem 0.9rem;
      color: #172033;
      font: inherit;
    }

    input:focus,
    select:focus {
      border-color: #2563eb;
      outline: 3px solid #dbeafe;
    }

    .form-actions {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-top: 0.5rem;
    }

    .primary-action,
    .secondary-action {
      border-radius: 999px;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
    }

    .primary-action {
      border: 0;
      background: #2563eb;
      color: #ffffff;
      cursor: pointer;
    }

    .primary-action:disabled {
      cursor: default;
      opacity: 0.65;
    }

    .secondary-action {
      color: #2563eb;
      text-decoration: none;
    }

    @media (max-width: 640px) {
      .user-edit-page {
        padding: 1rem;
      }

      .form-actions {
        align-items: stretch;
        flex-direction: column;
      }

      .primary-action,
      .secondary-action {
        text-align: center;
      }
    }
  `,
})
export class UserEditPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly userService = inject(UserService);

  protected readonly isLoading = signal(true);
  protected readonly isSaving = signal(false);
  protected readonly message = signal('');
  protected readonly errorMessage = signal('');

  protected readonly userId = computed(
    () => this.route.snapshot.paramMap.get('id') ?? '1',
  );

  protected readonly form = signal<UpdateUserRequest>({
    name: '',
    email: '',
    role: 'User',
    status: 'Active',
  });

  public ngOnInit(): void {
    this.userService
      .getUserById(this.userId())
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (user) => {
          this.form.set({
            name: user.name,
            email: user.email,
            role: user.role,
            status: user.status,
          });
        },
        error: () => {
          this.errorMessage.set('Unable to load user.');
        },
      });
  }

  protected updateForm<Key extends keyof UpdateUserRequest>(
    key: Key,
    value: UpdateUserRequest[Key],
  ): void {
    this.form.update((form) => ({
      ...form,
      [key]: value,
    }));
  }

  protected saveUser(): void {
    if (this.isSaving()) {
      return;
    }

    this.isSaving.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.userService
      .updateUser(this.userId(), this.form())
      .pipe(
        finalize(() => {
          this.isSaving.set(false);
        }),
      )
      .subscribe({
        next: (user) => {
          this.message.set('User saved successfully.');
          void this.router.navigate(['/users', user.id]);
        },
        error: () => {
          this.errorMessage.set('Unable to save user.');
        },
      });
  }
}
