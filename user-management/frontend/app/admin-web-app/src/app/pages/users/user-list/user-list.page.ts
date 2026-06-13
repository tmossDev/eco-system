import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { CrossAppNavigationService } from '../../../core/services/navigation/cross-app-navigation.service';
import { UserSummary } from '../../../core/services/user/user.model';
import { UserService } from '../../../core/services/user/user.service';

@Component({
  selector: 'app-user-list-page',
  imports: [RouterLink],
  template: `
    <section class="users-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">User management</p>
          <h1>Users</h1>
          <p class="description">
            View user accounts, check account status, and manage profile
            details.
          </p>
        </div>

        <button type="button" class="primary-action">Add user</button>
      </div>

      @if (isLoading()) {
        <p class="state-message">Loading users...</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <div class="table-card">
        <div class="table-header">
          <h2>All users</h2>
          <span>{{ users().length }} users</span>
        </div>

        <div class="table-wrap">
          <table>
            <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
              <th>Status</th>
              <th class="actions-column">Actions</th>
            </tr>
            </thead>

            <tbody>
              @for (user of users(); track user.id) {
                <tr>
                  <td>
                    <strong>{{ user.name }}</strong>
                  </td>
                  <td>{{ user.email }}</td>
                  <td>{{ user.role }}</td>
                  <td>
                    <span class="status" [class]="user.status.toLowerCase()">
                      {{ user.status }}
                    </span>
                  </td>
                  <td class="actions">
                    <a [routerLink]="['/users', user.id]">View</a>
                    <a [routerLink]="['/users', user.id, 'edit']">Edit</a>
                    <a [href]="storefrontUrl(user)">Storefront</a>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      </div>
    </section>
  `,
  styles: `
    .users-page {
      padding: 2rem;
      color: #172033;
    }

    .page-header,
    .table-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
    }

    .page-header {
      margin-bottom: 2rem;
    }

    .eyebrow {
      margin: 0 0 0.5rem;
      color: #56657f;
      font-size: 0.8rem;
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h1,
    h2 {
      margin: 0;
    }

    h1 {
      font-size: clamp(2rem, 4vw, 3rem);
      letter-spacing: -0.04em;
    }

    h2 {
      font-size: 1.15rem;
    }

    .description {
      max-width: 42rem;
      margin: 0.75rem 0 0;
      color: #56657f;
      line-height: 1.6;
    }

    .state-message,
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

    .error-message {
      background: #fee2e2;
      color: #991b1b;
    }

    .primary-action {
      border: 0;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
    }

    .table-card {
      overflow: hidden;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .table-header {
      padding: 1.25rem;
      border-bottom: 1px solid #dbe3ef;
    }

    .table-header span {
      color: #56657f;
    }

    .table-wrap {
      overflow-x: auto;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 720px;
    }

    th,
    td {
      padding: 1rem 1.25rem;
      border-bottom: 1px solid #eef2f7;
      text-align: left;
    }

    th {
      color: #56657f;
      font-size: 0.8rem;
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }

    tbody tr:last-child td {
      border-bottom: 0;
    }

    .status {
      display: inline-flex;
      border-radius: 999px;
      padding: 0.3rem 0.65rem;
      font-size: 0.85rem;
      font-weight: 700;
    }

    .active {
      background: #dcfce7;
      color: #166534;
    }

    .pending {
      background: #fef3c7;
      color: #92400e;
    }

    .suspended {
      background: #fee2e2;
      color: #991b1b;
    }

    .actions-column {
      width: 16rem;
    }

    .actions {
      display: flex;
      gap: 0.75rem;
    }

    .actions a {
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
    }

    @media (max-width: 700px) {
      .users-page {
        padding: 1rem;
      }

      .page-header {
        display: grid;
      }
    }
  `,
})
export class UserListPage implements OnInit {
  private readonly userService = inject(UserService);
  private readonly crossAppNavigation = inject(CrossAppNavigationService);

  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');
  protected readonly users = signal<UserSummary[]>([]);

  public ngOnInit(): void {
    this.userService
      .getUsers()
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (users) => {
          this.users.set(users);
        },
        error: () => {
          this.errorMessage.set('Unable to load users.');
        },
      });
  }

  protected storefrontUrl(user: UserSummary): string {
    return this.crossAppNavigation.buildUrl('storefront-web-app', '/', user);
  }
}
