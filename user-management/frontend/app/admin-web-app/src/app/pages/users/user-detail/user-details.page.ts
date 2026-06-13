import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { CrossAppNavigationService } from '../../../core/services/navigation/cross-app-navigation.service';
import { UserDetails } from '../../../core/services/user/user.model';
import { UserService } from '../../../core/services/user/user.service';

@Component({
  selector: 'app-user-detail-page',
  imports: [RouterLink],
  template: `
    <section class="user-detail-page">
      <a routerLink="/users" class="back-link">← Back to users</a>

      @if (isLoading()) {
        <p class="state-message">Loading user...</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      @if (user()) {
        <div class="profile-card">
          <div class="avatar">{{ initials() }}</div>

          <div class="profile-content">
            <p class="eyebrow">User profile</p>
            <h1>{{ user()!.name }}</h1>
            <p class="email">{{ user()!.email }}</p>

            <dl>
              <div>
                <dt>User ID</dt>
                <dd>{{ user()!.id }}</dd>
              </div>

              <div>
                <dt>Role</dt>
                <dd>{{ user()!.role }}</dd>
              </div>

              <div>
                <dt>Status</dt>
                <dd>{{ user()!.status }}</dd>
              </div>
            </dl>

            <div class="actions">
              <a [routerLink]="['/users', user()!.id, 'edit']" class="primary-action">
                Edit user
              </a>
              <a [href]="storefrontUrl()" class="secondary-action">
                View storefront as user
              </a>
            </div>
          </div>
        </div>
      }
    </section>
  `,
  styles: `
    .user-detail-page {
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

    .profile-card {
      display: grid;
      grid-template-columns: auto 1fr;
      gap: 1.5rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .avatar {
      display: grid;
      place-items: center;
      width: 5rem;
      height: 5rem;
      border-radius: 50%;
      background: #dbeafe;
      color: #1d4ed8;
      font-size: 1.5rem;
      font-weight: 800;
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

    .email {
      margin: 0.5rem 0 1.5rem;
      color: #56657f;
    }

    dl {
      display: grid;
      gap: 1rem;
      margin: 0 0 1.5rem;
    }

    dl div {
      display: grid;
      gap: 0.25rem;
    }

    dt {
      color: #56657f;
      font-size: 0.8rem;
      font-weight: 700;
      text-transform: uppercase;
    }

    dd {
      margin: 0;
      font-weight: 700;
    }

    .actions {
      display: flex;
      flex-wrap: wrap;
      gap: 0.75rem;
    }

    .primary-action,
    .secondary-action {
      display: inline-flex;
      width: fit-content;
      border-radius: 999px;
      padding: 0.75rem 1rem;
      font-weight: 700;
      text-decoration: none;
    }

    .primary-action {
      background: #2563eb;
      color: #ffffff;
    }

    .secondary-action {
      background: #e8eef8;
      color: #172033;
    }

    @media (max-width: 640px) {
      .user-detail-page {
        padding: 1rem;
      }

      .profile-card {
        grid-template-columns: 1fr;
      }
    }
  `,
})
export class UserDetailPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly userService = inject(UserService);
  private readonly crossAppNavigation = inject(CrossAppNavigationService);

  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');
  protected readonly user = signal<UserDetails | null>(null);

  protected readonly initials = computed(() => {
    const user = this.user();

    if (!user) {
      return '';
    }

    return user.name
      .split(' ')
      .map((namePart) => namePart[0])
      .join('')
      .toUpperCase();
  });

  public ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id') ?? '1';

    this.userService
      .getUserById(id)
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (user) => {
          this.user.set(user);
        },
        error: () => {
          this.errorMessage.set('Unable to load user.');
        },
      });
  }

  protected storefrontUrl(): string {
    const user = this.user();

    return this.crossAppNavigation.buildUrl(
      'storefront-web-app',
      '/',
      user ?? undefined,
    );
  }
}
