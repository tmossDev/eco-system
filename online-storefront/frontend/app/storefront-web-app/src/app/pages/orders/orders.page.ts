import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';
import { Order } from '../../core/services/cart.models';
import { CartService } from '../../core/services/cart.service';

@Component({
  selector: 'app-orders-page',
  imports: [RouterLink],
  template: `
    <section class="orders-page">
      <div class="heading">
        <div>
          <p class="eyebrow">Your purchases</p>
          <h1>Order history</h1>
        </div>
        <a routerLink="/">Continue shopping</a>
      </div>

      @if (!authService.isAuthenticated()) {
        <div class="empty">
          <h2>Sign in to view orders</h2>
          <p>Your checkout history is saved with your account.</p>
          <a class="primary" routerLink="/auth/login">Sign in</a>
        </div>
      } @else if (errorMessage()) {
        <p class="error">{{ errorMessage() }}</p>
      } @else if (isLoading()) {
        <p class="loading">Loading your orders...</p>
      } @else if (!orders().length) {
        <div class="empty">
          <h2>No orders yet</h2>
          <p>Completed checkouts will appear here.</p>
          <a class="primary" routerLink="/">Explore the collection</a>
        </div>
      } @else {
        <div class="orders">
          @for (order of orders(); track order.id) {
            <article>
              <header>
                <div>
                  <h2>Order #{{ order.id }}</h2>
                  <p>{{ formatDate(order.created_at) }}</p>
                </div>
                <span>{{ order.status }}</span>
              </header>

              <div class="summary">
                <div><span>Items</span><strong>{{ order.item_count }}</strong></div>
                <div><span>Total</span><strong>{{ formatMoney(order.subtotal_cents, order.currency) }}</strong></div>
              </div>

              <div class="items">
                @for (item of order.items; track item.product_id) {
                  <div class="item">
                    <span>{{ item.quantity }}x</span>
                    <strong>{{ item.name }}</strong>
                    <small>{{ formatMoney(item.line_total_cents, item.currency) }}</small>
                  </div>
                }
              </div>
            </article>
          }
        </div>
      }
    </section>
  `,
  styles: `
    .orders-page { max-width: 76rem; min-height: 55vh; margin: 0 auto; padding: clamp(2rem, 6vw, 5rem) clamp(1rem, 4vw, 3rem); }
    .heading { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-bottom: 2rem; }
    .eyebrow { margin: 0 0 .5rem; color: #728070; font-size: .74rem; font-weight: 700; letter-spacing: .16em; text-transform: uppercase; }
    h1, h2 { color: #28432f; font-family: Georgia, serif; }
    h1 { margin: 0; font-size: clamp(3rem, 6vw, 5rem); letter-spacing: -.07em; line-height: .95; }
    h2 { margin: 0; font-size: 1.35rem; }
    a { color: #31543c; font-weight: 700; text-decoration: none; }
    .orders { display: grid; gap: 1rem; }
    article { display: grid; gap: 1rem; border-radius: 1rem; background: #fffdf8; padding: 1rem; }
    article header, .summary, .item { display: flex; align-items: center; justify-content: space-between; gap: 1rem; }
    article header p { margin: .35rem 0 0; color: #718071; }
    article header span { border-radius: 999px; background: #e8ede2; padding: .45rem .7rem; color: #31543c; font-weight: 700; }
    .summary { border-top: 1px solid #e4dfd2; border-bottom: 1px solid #e4dfd2; padding: .85rem 0; color: #31543c; }
    .items { display: grid; gap: .6rem; }
    .item { justify-content: flex-start; color: #31543c; }
    .item span { color: #718071; font-weight: 700; }
    .item strong { flex: 1; }
    .item small, .loading, .empty p { color: #718071; line-height: 1.6; }
    .empty { border-radius: 1rem; background: #fffdf8; padding: 2rem; text-align: center; }
    .primary { display: inline-flex; margin-top: .5rem; border-radius: 999px; background: #31543c; padding: .8rem 1rem; color: white; }
    .error { color: #9b3f34; }
    @media (max-width: 640px) {
      .heading, article header, .summary { align-items: flex-start; flex-direction: column; }
      .item { align-items: flex-start; flex-direction: column; }
    }
  `,
})
export class OrdersPage implements OnInit {
  protected readonly authService = inject(AuthService);
  private readonly cartService = inject(CartService);
  protected readonly orders = signal<Order[]>([]);
  protected readonly isLoading = signal(false);
  protected readonly errorMessage = signal('');

  ngOnInit(): void {
    if (!this.authService.isAuthenticated()) {
      return;
    }
    this.isLoading.set(true);
    this.cartService.listOrders()
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (orders) => this.orders.set(orders),
        error: () => this.errorMessage.set('Unable to load your orders. Please try again.'),
      });
  }

  protected formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency || 'USD',
    }).format(cents / 100);
  }

  protected formatDate(value: string): string {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value;
    }
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    }).format(date);
  }
}
