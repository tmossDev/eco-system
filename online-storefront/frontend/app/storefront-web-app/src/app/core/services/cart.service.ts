import { HttpClient } from '@angular/common/http';
import { Injectable, computed, effect, inject, signal } from '@angular/core';
import { Observable, tap } from 'rxjs';

import { AuthService } from './auth.service';
import { Cart, Order } from './cart.models';

@Injectable({ providedIn: 'root' })
export class CartService {
  private readonly authService = inject(AuthService);
  private readonly http = inject(HttpClient);
  private readonly cartSignal = signal<Cart | null>(null);

  readonly cart = computed(() => this.cartSignal());
  readonly itemCount = computed(() => this.cartSignal()?.item_count ?? 0);

  constructor() {
    effect(() => {
      if (this.authService.isAuthenticated()) {
        this.refresh().subscribe({ error: () => this.cartSignal.set(null) });
      } else {
        this.cartSignal.set(null);
      }
    });
  }

  refresh(): Observable<Cart> {
    return this.http
      .get<Cart>('/api/cart')
      .pipe(tap((cart) => this.cartSignal.set(cart)));
  }

  addItem(productId: string, quantity = 1): Observable<Cart> {
    return this.http
      .post<Cart>('/api/cart/items', {
        product_id: Number(productId),
        quantity,
      })
      .pipe(tap((cart) => this.cartSignal.set(cart)));
  }

  updateItem(productId: number, quantity: number): Observable<Cart> {
    return this.http
      .put<Cart>(`/api/cart/items/${productId}`, { quantity })
      .pipe(tap((cart) => this.cartSignal.set(cart)));
  }

  removeItem(productId: number): Observable<Cart> {
    return this.http
      .delete<Cart>(`/api/cart/items/${productId}`)
      .pipe(tap((cart) => this.cartSignal.set(cart)));
  }

  clear(): Observable<Cart> {
    return this.http
      .delete<Cart>('/api/cart')
      .pipe(tap((cart) => this.cartSignal.set(cart)));
  }

  checkout(): Observable<Order> {
    return this.http
      .post<Order>('/api/cart/checkout', {})
      .pipe(tap(() => this.cartSignal.set(null)));
  }

  listOrders(): Observable<Order[]> {
    return this.http.get<Order[]>('/api/orders');
  }
}
