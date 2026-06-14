import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { Order } from '../../core/services/order/order.models';
import { OrderService } from '../../core/services/order/order.service';

@Component({
  selector: 'app-order-list-page',
  imports: [RouterLink],
  template: `
    <section class="orders-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">Fulfillment queue</p>
          <h1>Orders</h1>
          <p class="description">Review customer orders and move them through fulfillment.</p>
        </div>
      </div>

      @if (isLoading()) { <p class="state-message">Loading orders...</p> }
      @if (errorMessage()) { <p class="error-message">{{ errorMessage() }}</p> }

      <div class="table-shell">
        <table>
          <thead>
            <tr>
              <th>Order</th>
              <th>Status</th>
              <th>Items</th>
              <th>Total</th>
              <th>Created</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            @for (order of orders(); track order.id) {
              <tr>
                <td>#{{ order.id }}</td>
                <td><span class="status">{{ order.status }}</span></td>
                <td>{{ order.item_count }}</td>
                <td>{{ formatMoney(order.subtotal_cents, order.currency) }}</td>
                <td>{{ formatDate(order.created_at) }}</td>
                <td><a [routerLink]="['/orders', order.id]">Open</a></td>
              </tr>
            } @empty {
              <tr><td colspan="6">No orders yet.</td></tr>
            }
          </tbody>
        </table>
      </div>
    </section>
  `,
  styles: `
    .orders-page { padding: 2rem; color: #172033; }
    .page-header { margin-bottom: 1.5rem; }
    .eyebrow { margin: 0 0 .5rem; color: #56657f; font-size: .8rem; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; }
    h1 { margin: 0; font-size: clamp(2rem, 4vw, 3rem); letter-spacing: 0; }
    .description { max-width: 42rem; margin: .75rem 0 0; color: #56657f; line-height: 1.6; }
    .state-message, .error-message { margin: 0 0 1rem; border-radius: 8px; padding: .75rem .9rem; font-weight: 700; }
    .state-message { background: #eff6ff; color: #1d4ed8; }
    .error-message { background: #fee2e2; color: #991b1b; }
    .table-shell { overflow: auto; border: 1px solid #dbe3ef; border-radius: 8px; background: #fff; box-shadow: 0 10px 30px rgb(15 23 42 / 6%); }
    table { width: 100%; border-collapse: collapse; min-width: 44rem; }
    th, td { padding: .9rem 1rem; border-bottom: 1px solid #eef2f7; text-align: left; }
    th { color: #56657f; font-size: .78rem; text-transform: uppercase; }
    tr:last-child td { border-bottom: 0; }
    a { color: #2563eb; font-weight: 700; text-decoration: none; }
    .status { display: inline-flex; border-radius: 999px; background: #e8eef8; padding: .35rem .65rem; color: #172033; font-weight: 700; }
  `,
})
export class OrderListPage implements OnInit {
  private readonly orderService = inject(OrderService);
  protected readonly orders = signal<Order[]>([]);
  protected readonly isLoading = signal(false);
  protected readonly errorMessage = signal('');

  ngOnInit(): void {
    this.isLoading.set(true);
    this.orderService.listOrders()
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (orders) => this.orders.set(orders),
        error: () => this.errorMessage.set('Unable to load orders.'),
      });
  }

  protected formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: currency || 'USD' }).format(cents / 100);
  }

  protected formatDate(value: string): string {
    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? value : new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', year: 'numeric' }).format(date);
  }
}
