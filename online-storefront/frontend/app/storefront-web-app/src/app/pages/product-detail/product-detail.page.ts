import { Component, OnInit, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { Product, ProductPhoto } from '../../core/services/product.models';
import { ProductService } from '../../core/services/product.service';

@Component({
  selector: 'app-product-detail-page',
  imports: [RouterLink],
  template: `
    <section class="detail">
      <a class="back" routerLink="/">Back to collection</a>

      @if (isLoading()) {
        <p>Loading product...</p>
      }
      @if (errorMessage()) {
        <p class="error">{{ errorMessage() }}</p>
      }

      @if (product(); as item) {
        <div class="product">
          <div class="media">
            @if (primaryPhoto(item); as photo) {
              <img [src]="photo.url" [alt]="photo.alt_text || item.name">
            } @else {
              <span>{{ item.category }}</span>
            }
          </div>

          <div class="copy">
            <p class="eyebrow">{{ item.category }}</p>
            <h1>{{ item.name }}</h1>
            <strong class="price">{{ formatMoney(item) }}</strong>
            <p class="lead">{{ item.short_description }}</p>
            <p class="description">{{ item.description }}</p>
            <div class="notice">
              <strong>Available now</strong>
              <span>{{ item.inventory_count }} currently in stock</span>
            </div>
            <p class="next">Cart and checkout are arriving in the next storefront update.</p>
          </div>
        </div>
      }
    </section>
  `,
  styles: `
    .detail { max-width: 76rem; margin: 0 auto; padding: clamp(2rem, 6vw, 5rem) clamp(1rem, 4vw, 3rem); }
    .back { display: inline-flex; margin-bottom: 1.5rem; color: #637064; font-weight: 700; text-decoration: none; }
    .product { display: grid; grid-template-columns: minmax(18rem, 1fr) minmax(18rem, .9fr); gap: clamp(2rem, 6vw, 5rem); align-items: center; }
    .media { display: grid; aspect-ratio: 1 / 1.08; place-items: center; overflow: hidden; border-radius: 1.4rem; background: linear-gradient(145deg, #e2ddcf, #cdd8c3); color: #61705f; font-size: .8rem; font-weight: 700; letter-spacing: .14em; text-transform: uppercase; }
    img { width: 100%; height: 100%; object-fit: cover; }
    .eyebrow { color: #728070; font-size: .74rem; font-weight: 700; letter-spacing: .16em; text-transform: uppercase; }
    h1 { margin: .5rem 0 1rem; color: #28432f; font-family: Georgia, serif; font-size: clamp(3rem, 6vw, 5rem); letter-spacing: -.07em; line-height: .95; }
    .price { color: #31543c; font-size: 1.35rem; }
    .lead { margin: 2rem 0 .5rem; color: #455448; font-size: 1.1rem; font-weight: 700; line-height: 1.6; }
    .description { color: #6d786e; line-height: 1.8; }
    .notice { display: grid; gap: .2rem; margin-top: 1.8rem; border-radius: .9rem; background: #e8ede2; padding: 1rem; color: #31543c; }
    .notice span, .next { color: #718071; font-size: .86rem; }
    .next { margin-top: 1rem; }
    .error { color: #9b3f34; }
    @media (max-width: 760px) { .product { grid-template-columns: 1fr; } }
  `,
})
export class ProductDetailPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly productService = inject(ProductService);
  protected readonly product = signal<Product | null>(null);
  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      this.isLoading.set(false);
      this.errorMessage.set('This product could not be found.');
      return;
    }

    this.productService
      .getProduct(id)
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (product) => this.product.set(product),
        error: () => this.errorMessage.set('This product could not be found.'),
      });
  }

  protected primaryPhoto(product: Product): ProductPhoto | undefined {
    return product.photos?.find((photo) => photo.is_primary) ?? product.photos?.[0];
  }

  protected formatMoney(product: Product): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: product.currency,
    }).format(product.price_cents / 100);
  }
}
