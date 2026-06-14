import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { Order, OrderService } from '@eco/admin-features';

@Component({
  selector: 'app-dashboard-page',
  imports: [RouterLink],
  template: `
    <section class="dashboard-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">Operations</p>
          <h1>Dashboard</h1>
          <p class="description">Track open orders, fulfillment work, and cancelled orders.</p>
        </div>

        <a routerLink="/orders" class="primary-action">Manage orders</a>
      </div>

      @if (isLoading()) { <p class="state-message">Loading dashboard...</p> }
      @if (errorMessage()) { <p class="error-message">{{ errorMessage() }}</p> }

      <div class="stats-grid">
        <article><span>Total orders</span><strong>{{ orders().length }}</strong></article>
        <article><span>Submitted</span><strong>{{ statusCount('Order Submitted') }}</strong></article>
        <article><span>Fulfillment</span><strong>{{ statusCount('Order Fulfillment') }}</strong></article>
        <article><span>Out for delivery</span><strong>{{ statusCount('Order Out For Delivery') }}</strong></article>
      </div>

      <div class="panel">
        <h2>Recent orders</h2>
        @for (order of recentOrders(); track order.id) {
          <a class="order-row" [routerLink]="['/orders', order.id]">
            <span>#{{ order.id }}</span>
            <strong>{{ order.status }}</strong>
            <span>{{ formatMoney(order.subtotal_cents, order.currency) }}</span>
          </a>
        } @empty {
          <p class="description">No orders yet.</p>
        }
      </div>
    </section>
  `,
  styles: `
    .dashboard-page { padding: 2rem; color: #172033; }
    .page-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; margin-bottom: 2rem; }
    .eyebrow { margin: 0 0 .5rem; color: #56657f; font-size: .8rem; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; }
    h1, h2 { margin: 0; }
    h1 { font-size: clamp(2rem, 4vw, 3rem); letter-spacing: 0; }
    h2 { margin-bottom: 1rem; font-size: 1.15rem; }
    .description { max-width: 42rem; margin: .75rem 0 0; color: #56657f; line-height: 1.6; }
    .primary-action { display: inline-flex; align-items: center; border-radius: 999px; background: #2563eb; color: #fff; padding: .75rem 1rem; font-weight: 700; text-decoration: none; white-space: nowrap; }
    .state-message, .error-message { margin: 0 0 1rem; border-radius: 8px; padding: .75rem .9rem; font-weight: 700; }
    .state-message { background: #eff6ff; color: #1d4ed8; }
    .error-message { background: #fee2e2; color: #991b1b; }
    .stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 1rem; margin-bottom: 1rem; }
    article, .panel { border: 1px solid #dbe3ef; border-radius: 8px; background: #fff; box-shadow: 0 10px 30px rgb(15 23 42 / 6%); }
    article { padding: 1.25rem; }
    article strong { display: block; margin-top: .5rem; font-size: 2rem; }
    article span { color: #56657f; }
    .panel { padding: 1.25rem; }
    .order-row { display: grid; grid-template-columns: 1fr 1fr auto; gap: 1rem; align-items: center; border-top: 1px solid #eef2f7; padding: .85rem 0; color: #172033; text-decoration: none; }
    .order-row span { color: #56657f; }
    @media (max-width: 900px) { .page-header, .stats-grid { display: grid; grid-template-columns: 1fr; } }
  `,
})
export class DashboardPage implements OnInit {
  private readonly orderService = inject(OrderService);
  protected readonly orders = signal<Order[]>([]);
  protected readonly isLoading = signal(false);
  protected readonly errorMessage = signal('');
  protected readonly recentOrders = computed(() => this.orders().slice(0, 5));

  ngOnInit(): void {
    this.isLoading.set(true);
    this.orderService.listOrders()
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (orders) => this.orders.set(orders),
        error: () => this.errorMessage.set('Unable to load order dashboard.'),
      });
  }

  protected statusCount(status: Order['status']): number {
    return this.orders().filter((order) => order.status === status).length;
  }

  protected formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: currency || 'USD' }).format(cents / 100);
  }
}
