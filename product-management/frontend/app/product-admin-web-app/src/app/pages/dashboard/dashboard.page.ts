import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { DashboardService } from '../../core/services/dashboard/dashboard.service';
import { DashboardStat } from '../../core/services/dashboard/dashboard.models';

@Component({
  selector: 'app-dashboard-page',
  imports: [RouterLink],
  template: `
    <section class="dashboard-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">Overview</p>
          <h1>Dashboard</h1>
          <p class="description">
            Monitor product coverage, catalog readiness, and inventory health.
          </p>
        </div>

        <a routerLink="/products" class="primary-action">Manage products</a>
      </div>

      @if (isLoading()) {
        <p class="state-message">Loading dashboard...</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <div class="stats-grid">
        @for (stat of stats(); track stat.label) {
          <article class="stat-card">
            <span class="stat-label">{{ stat.label }}</span>
            <strong>{{ stat.value }}</strong>
            <small>{{ stat.caption }}</small>
          </article>
        }
      </div>

      <div class="content-grid">
        <article class="panel">
          <h2>Recent activity</h2>

          <ul class="activity-list">
            @for (activity of recentActivity(); track activity) {
              <li>{{ activity }}</li>
            }
          </ul>
        </article>

        <article class="panel">
          <h2>Quick actions</h2>

          <div class="quick-actions">
            <a routerLink="/products" class="action-card">
              <strong>Products</strong>
              <span>View and manage catalog items</span>
            </a>

            <a routerLink="/settings" class="action-card">
              <strong>Settings</strong>
              <span>Configure admin preferences</span>
            </a>
          </div>
        </article>
      </div>
    </section>
  `,
  styles: `
    .dashboard-page {
      padding: 2rem;
      color: #172033;
    }

    .page-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
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
      margin-bottom: 1rem;
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

    .primary-action,
    .action-card {
      text-decoration: none;
    }

    .primary-action {
      display: inline-flex;
      align-items: center;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font-weight: 700;
      white-space: nowrap;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 1rem;
      margin-bottom: 1rem;
    }

    .stat-card,
    .panel,
    .action-card {
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .stat-card {
      padding: 1.25rem;
    }

    .stat-card strong {
      display: block;
      margin: 0.5rem 0;
      font-size: 2rem;
    }

    .stat-label,
    .stat-card small {
      color: #56657f;
    }

    .content-grid {
      display: grid;
      grid-template-columns: 1.2fr 0.8fr;
      gap: 1rem;
    }

    .panel {
      padding: 1.25rem;
    }

    .activity-list {
      margin: 0;
      padding-left: 1.2rem;
      color: #56657f;
      line-height: 1.8;
    }

    .quick-actions {
      display: grid;
      gap: 0.75rem;
    }

    .action-card {
      display: grid;
      gap: 0.25rem;
      padding: 1rem;
      color: #172033;
    }

    .action-card span {
      color: #56657f;
    }

    @media (max-width: 900px) {
      .page-header,
      .content-grid {
        grid-template-columns: 1fr;
      }

      .page-header {
        display: grid;
      }

      .stats-grid {
        grid-template-columns: repeat(2, minmax(0, 1fr));
      }
    }

    @media (max-width: 560px) {
      .dashboard-page {
        padding: 1rem;
      }

      .stats-grid {
        grid-template-columns: 1fr;
      }
    }
  `,
})
export class DashboardPage implements OnInit {
  private readonly dashboardService = inject(DashboardService);

  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');
  protected readonly stats = signal<DashboardStat[]>([]);
  protected readonly recentActivity = signal<string[]>([]);

  public ngOnInit(): void {
    this.dashboardService
      .getSummary()
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (summary) => {
          this.stats.set(summary.stats);
          this.recentActivity.set(summary.recentActivity);
        },
        error: () => {
          this.errorMessage.set('Unable to load dashboard summary.');
        },
      });
  }
}
