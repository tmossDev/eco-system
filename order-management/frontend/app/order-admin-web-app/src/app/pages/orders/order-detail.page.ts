import { Component, OnInit, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { Order, OrderStatus } from '../../core/services/order/order.models';
import { OrderService } from '../../core/services/order/order.service';

const STATUSES: OrderStatus[] = ['Created', 'Paid', 'Fulfilled', 'Cancelled'];

@Component({
  selector: 'app-order-detail-page',
  imports: [RouterLink],
  template: `
    <section class="order-page">
      <a routerLink="/orders" class="back-link">Back to orders</a>

      @if (isLoading()) { <p class="state-message">Loading order...</p> }
      @if (errorMessage()) { <p class="error-message">{{ errorMessage() }}</p> }

      @if (order()) {
        <div class="page-header">
          <div>
            <p class="eyebrow">Order #{{ order()!.id }}</p>
            <h1>{{ order()!.status }}</h1>
            <p class="description">{{ order()!.item_count }} items for {{ formatMoney(order()!.subtotal_cents, order()!.currency) }}</p>
          </div>

          <div class="status-actions">
            @for (status of statuses; track status) {
              <button type="button" [disabled]="isSaving() || order()!.status === status" (click)="setStatus(status)">
                {{ status }}
              </button>
            }
          </div>
        </div>

        <div class="items">
          @for (item of order()!.items; track item.product_id) {
            <article>
              <div>
                <strong>{{ item.name }}</strong>
                <span>{{ item.sku }}</span>
              </div>
              <span>{{ item.quantity }} x {{ formatMoney(item.price_cents, item.currency) }}</span>
              <strong>{{ formatMoney(item.line_total_cents, item.currency) }}</strong>
            </article>
          }
        </div>
      }
    </section>
  `,
  styles: `
    .order-page { padding: 2rem; color: #172033; }
    .back-link { display: inline-flex; margin-bottom: 1rem; color: #2563eb; font-weight: 700; text-decoration: none; }
    .page-header { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; margin-bottom: 1.5rem; }
    .eyebrow { margin: 0 0 .5rem; color: #56657f; font-size: .8rem; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; }
    h1 { margin: 0; font-size: clamp(2rem, 4vw, 3rem); letter-spacing: 0; }
    .description { margin: .75rem 0 0; color: #56657f; line-height: 1.6; }
    .state-message, .error-message { margin: 0 0 1rem; border-radius: 8px; padding: .75rem .9rem; font-weight: 700; }
    .state-message { background: #eff6ff; color: #1d4ed8; }
    .error-message { background: #fee2e2; color: #991b1b; }
    .status-actions { display: flex; flex-wrap: wrap; gap: .5rem; }
    button { border: 0; border-radius: 999px; background: #2563eb; padding: .65rem .9rem; color: white; cursor: pointer; font-weight: 700; }
    button:disabled { cursor: default; opacity: .55; }
    .items { display: grid; gap: .75rem; }
    article { display: grid; grid-template-columns: minmax(0, 1fr) auto auto; gap: 1rem; align-items: center; border: 1px solid #dbe3ef; border-radius: 8px; background: #fff; padding: 1rem; }
    article div { display: grid; gap: .25rem; }
    article span { color: #56657f; }
    @media (max-width: 760px) { .page-header, article { display: grid; grid-template-columns: 1fr; } }
  `,
})
export class OrderDetailPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly orderService = inject(OrderService);
  protected readonly statuses = STATUSES;
  protected readonly order = signal<Order | null>(null);
  protected readonly isLoading = signal(false);
  protected readonly isSaving = signal(false);
  protected readonly errorMessage = signal('');

  ngOnInit(): void {
    this.load();
  }

  protected setStatus(status: OrderStatus): void {
    const orderID = this.route.snapshot.paramMap.get('id');
    if (!orderID) return;
    this.isSaving.set(true);
    this.errorMessage.set('');
    this.orderService.updateStatus(orderID, { status })
      .pipe(finalize(() => this.isSaving.set(false)))
      .subscribe({
        next: (order) => this.order.set(order),
        error: () => this.errorMessage.set('Unable to update order status.'),
      });
  }

  protected formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: currency || 'USD' }).format(cents / 100);
  }

  private load(): void {
    const orderID = this.route.snapshot.paramMap.get('id');
    if (!orderID) return;
    this.isLoading.set(true);
    this.orderService.getOrder(orderID)
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (order) => this.order.set(order),
        error: () => this.errorMessage.set('Unable to load order.'),
      });
  }
}
