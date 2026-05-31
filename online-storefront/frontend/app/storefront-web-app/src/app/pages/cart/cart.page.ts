import { Component, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { AuthService } from '../../core/services/auth.service';
import { CartItem } from '../../core/services/cart.models';
import { CartService } from '../../core/services/cart.service';

@Component({
  selector: 'app-cart-page',
  imports: [RouterLink],
  template: `
    <section class="cart-page">
      <div class="heading">
        <div>
          <p class="eyebrow">Your selection</p>
          <h1>Shopping cart</h1>
        </div>
        <a routerLink="/">Continue shopping</a>
      </div>

      @if (!authService.isAuthenticated()) {
        <div class="empty">
          <h2>Sign in to use your cart</h2>
          <p>Your cart is stored with your account, ready when you return.</p>
          <a class="primary" routerLink="/auth/login">Sign in</a>
        </div>
      } @else if (errorMessage()) {
        <p class="error">{{ errorMessage() }}</p>
      } @else if (!cartService.cart()) {
        <p class="loading">Loading your cart...</p>
      } @else if (!cartService.cart()!.items.length) {
        <div class="empty">
          <h2>Your cart is empty</h2>
          <p>There is plenty to browse when you are ready.</p>
          <a class="primary" routerLink="/">Explore the collection</a>
        </div>
      } @else {
        <div class="cart-grid">
          <div class="items">
            @for (item of cartService.cart()!.items; track item.product_id) {
              <article>
                <div class="media">
                  @if (item.thumbnail_url) {
                    <img [src]="item.thumbnail_url" [alt]="item.name">
                  } @else {
                    <span>{{ item.sku }}</span>
                  }
                </div>
                <div class="item-copy">
                  <h2>{{ item.name }}</h2>
                  <p>{{ formatMoney(item.price_cents, item.currency) }} each</p>
                  <div class="actions">
                    <button type="button" aria-label="Decrease quantity" [disabled]="isSaving()" (click)="changeQuantity(item, -1)">-</button>
                    <strong>{{ item.quantity }}</strong>
                    <button type="button" aria-label="Increase quantity" [disabled]="isSaving()" (click)="changeQuantity(item, 1)">+</button>
                    <button class="remove" type="button" [disabled]="isSaving()" (click)="remove(item)">Remove</button>
                  </div>
                </div>
                <strong>{{ formatMoney(item.line_total_cents, item.currency) }}</strong>
              </article>
            }
          </div>

          <aside>
            <p>Order summary</p>
            <div><span>Items</span><strong>{{ cartService.itemCount() }}</strong></div>
            <div><span>Subtotal</span><strong>{{ formatMoney(cartService.cart()!.subtotal_cents, cartService.cart()!.currency) }}</strong></div>
            <small>Checkout will be added in a future update.</small>
            <button class="clear" type="button" [disabled]="isSaving()" (click)="clear()">Clear cart</button>
          </aside>
        </div>
      }
    </section>
  `,
  styles: `
    .cart-page { max-width: 76rem; min-height: 55vh; margin: 0 auto; padding: clamp(2rem, 6vw, 5rem) clamp(1rem, 4vw, 3rem); }
    .heading, article, aside div, .actions { display: flex; align-items: center; }
    .heading { justify-content: space-between; gap: 1rem; margin-bottom: 2rem; }
    .eyebrow { margin: 0 0 .5rem; color: #728070; font-size: .74rem; font-weight: 700; letter-spacing: .16em; text-transform: uppercase; }
    h1, h2 { color: #28432f; font-family: Georgia, serif; }
    h1 { margin: 0; font-size: clamp(3rem, 6vw, 5rem); letter-spacing: -.07em; line-height: .95; }
    h2 { margin: 0; font-size: 1.35rem; }
    a { color: #31543c; font-weight: 700; text-decoration: none; }
    .cart-grid { display: grid; grid-template-columns: minmax(0, 1fr) 19rem; gap: 2rem; align-items: start; }
    .items { display: grid; gap: 1rem; }
    article { gap: 1rem; border-radius: 1rem; background: #fffdf8; padding: 1rem; }
    .media { display: grid; width: 5.5rem; height: 5.5rem; flex: 0 0 auto; place-items: center; overflow: hidden; border-radius: .8rem; background: #e2ddcf; color: #61705f; font-size: .65rem; font-weight: 700; text-align: center; }
    img { width: 100%; height: 100%; object-fit: cover; }
    .item-copy { flex: 1; }
    .item-copy p { margin: .35rem 0 .65rem; color: #718071; font-size: .86rem; }
    .actions { gap: .6rem; }
    button { border: 0; border-radius: 999px; background: #e8ede2; padding: .45rem .7rem; color: #31543c; cursor: pointer; font-weight: 700; }
    button:disabled { cursor: wait; opacity: .6; }
    .remove, .clear { background: transparent; color: #9b3f34; }
    aside { display: grid; gap: .9rem; border-radius: 1rem; background: #e8ede2; padding: 1.2rem; color: #31543c; }
    aside p { margin: 0; font-family: Georgia, serif; font-size: 1.35rem; font-weight: 700; }
    aside div { justify-content: space-between; }
    aside small, .loading, .empty p { color: #718071; line-height: 1.6; }
    .empty { border-radius: 1rem; background: #fffdf8; padding: 2rem; text-align: center; }
    .primary { display: inline-flex; margin-top: .5rem; border-radius: 999px; background: #31543c; padding: .8rem 1rem; color: white; }
    .error { color: #9b3f34; }
    @media (max-width: 760px) { .cart-grid { grid-template-columns: 1fr; } article { align-items: start; flex-wrap: wrap; } article > strong { margin-left: 6.5rem; } }
  `,
})
export class CartPage {
  protected readonly authService = inject(AuthService);
  protected readonly cartService = inject(CartService);
  protected readonly isSaving = signal(false);
  protected readonly errorMessage = signal('');

  protected changeQuantity(item: CartItem, difference: number): void {
    const quantity = item.quantity + difference;
    if (quantity < 1) {
      this.remove(item);
      return;
    }
    this.save(this.cartService.updateItem(item.product_id, quantity));
  }

  protected remove(item: CartItem): void {
    this.save(this.cartService.removeItem(item.product_id));
  }

  protected clear(): void {
    this.save(this.cartService.clear());
  }

  protected formatMoney(cents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency || 'USD',
    }).format(cents / 100);
  }

  private save(request: ReturnType<CartService['clear']>): void {
    this.isSaving.set(true);
    this.errorMessage.set('');
    request
      .pipe(finalize(() => this.isSaving.set(false)))
      .subscribe({
        error: () => this.errorMessage.set('Unable to update your cart. Please try again.'),
      });
  }
}
