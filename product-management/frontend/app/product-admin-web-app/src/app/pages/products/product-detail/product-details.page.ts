import { Component, OnInit, inject, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import {
  ProductDetails,
  ProductPhoto,
} from '../../../core/services/product/product.model';
import { ProductService } from '../../../core/services/product/product.service';

@Component({
  selector: 'app-product-detail-page',
  imports: [RouterLink],
  template: `
    <section class="product-detail-page">
      <a routerLink="/products" class="back-link">← Back to products</a>

      @if (isLoading()) {
        <p class="state-message">Loading product...</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      @if (product()) {
        <div class="profile-card">
          <div class="media-panel">
            @if (primaryPhoto(product()!)) {
              <img
                class="hero-photo"
                [src]="primaryPhoto(product()!)!.url"
                [alt]="primaryPhoto(product()!)!.alt_text || product()!.name"
                fetchpriority="high"
              />
            } @else {
              <div class="avatar">{{ product()!.sku.slice(0, 2) }}</div>
            }

            @if (productPhotos(product()!).length > 1) {
              <div class="photo-strip">
                @for (photo of productPhotos(product()!); track photo.url) {
                  <img
                    [src]="photo.thumbnail_url || photo.url"
                    [alt]="photo.alt_text || product()!.name"
                    loading="lazy"
                    decoding="async"
                    width="56"
                    height="56"
                  />
                }
              </div>
            }
          </div>

          <div class="profile-content">
            <p class="eyebrow">Product details</p>
            <h1>{{ product()!.name }}</h1>
            <p class="short-description">{{ product()!.short_description }}</p>
            <p class="description">{{ product()!.description }}</p>

            <dl>
              <div>
                <dt>Product ID</dt>
                <dd>{{ product()!.id }}</dd>
              </div>

              <div>
                <dt>SKU</dt>
                <dd>{{ product()!.sku }}</dd>
              </div>

              <div>
                <dt>Category</dt>
                <dd>{{ product()!.category }}</dd>
              </div>

              <div>
                <dt>Price</dt>
                <dd>{{ formatMoney(product()!.price_cents, product()!.currency) }}</dd>
              </div>

              <div>
                <dt>Discounts</dt>
                <dd>{{ discountSummary(product()!) }}</dd>
              </div>

              <div>
                <dt>Inventory</dt>
                <dd>{{ product()!.inventory_count }}</dd>
              </div>

              <div>
                <dt>Status</dt>
                <dd>{{ product()!.status }}</dd>
              </div>
            </dl>

            <a [routerLink]="['/products', product()!.id, 'edit']" class="primary-action">
              Edit product
            </a>
          </div>
        </div>
      }
    </section>
  `,
  styles: `
    .product-detail-page {
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
      grid-template-columns: minmax(15rem, 22rem) 1fr;
      gap: 1.5rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .media-panel {
      display: grid;
      align-content: start;
      gap: 0.75rem;
    }

    .hero-photo,
    .avatar {
      width: 100%;
      aspect-ratio: 4 / 3;
      border-radius: 0.75rem;
    }

    .hero-photo {
      object-fit: cover;
      background: #eef2f7;
    }

    .avatar {
      display: grid;
      place-items: center;
      background: #dbeafe;
      color: #1d4ed8;
      font-size: 2rem;
      font-weight: 800;
    }

    .photo-strip {
      display: flex;
      gap: 0.5rem;
      overflow-x: auto;
      padding-bottom: 0.2rem;
    }

    .photo-strip img {
      width: 3.5rem;
      height: 3.5rem;
      flex: 0 0 auto;
      border-radius: 0.5rem;
      object-fit: cover;
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

    .short-description,
    .description {
      max-width: 48rem;
      color: #56657f;
      line-height: 1.6;
    }

    .short-description {
      margin: 0.5rem 0 0;
      font-weight: 700;
    }

    .description {
      margin: 0.5rem 0 1.5rem;
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

    .primary-action {
      display: inline-flex;
      width: fit-content;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font-weight: 700;
      text-decoration: none;
    }

    @media (max-width: 640px) {
      .product-detail-page {
        padding: 1rem;
      }

      .profile-card {
        grid-template-columns: 1fr;
      }
    }
  `,
})
export class ProductDetailPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly productService = inject(ProductService);

  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');
  protected readonly product = signal<ProductDetails | null>(null);

  public ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id') ?? '1';

    this.productService
      .getProductById(id)
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (product) => {
          this.product.set(product);
        },
        error: () => {
          this.errorMessage.set('Unable to load product.');
        },
      });
  }

  protected formatMoney(priceCents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(priceCents / 100);
  }

  protected discountSummary(product: ProductDetails): string {
    const discounts = product.discounts ?? [];
    if (discounts.length === 0) {
      return 'None';
    }

    return discounts
      .map((discount) => {
        const value =
          discount.discount_type === 'Percentage'
            ? `${(discount.percentage_basis_points ?? 0) / 100}%`
            : this.formatMoney(discount.amount_cents ?? 0, discount.currency);

        return `${discount.name} (${value})`;
      })
      .join(', ');
  }

  protected primaryPhoto(product: ProductDetails): ProductPhoto | null {
    const photos = this.productPhotos(product);

    return photos.find((photo) => photo.is_primary) ?? photos[0] ?? null;
  }

  protected productPhotos(product: ProductDetails): ProductPhoto[] {
    return product.photos ?? [];
  }
}
