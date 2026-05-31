import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { Product, ProductPhoto } from '../../core/services/product.models';
import { ProductService } from '../../core/services/product.service';

@Component({
  selector: 'app-catalog-page',
  imports: [RouterLink],
  template: `
    <section class="hero">
      <div>
        <p class="eyebrow">Freshly selected</p>
        <h1>Good things for the everyday.</h1>
        <p class="intro">
          A considered collection of useful pieces, made to settle naturally
          into your home and routines.
        </p>
        <a href="#catalog">Explore the collection</a>
      </div>
      <div class="hero-card">
        <span>New season</span>
        <strong>Simple pieces.<br>Long lives.</strong>
        <small>Curated for slower, better days.</small>
      </div>
    </section>

    <section class="catalog" id="catalog">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Shop the catalog</p>
          <h2>Made for living</h2>
        </div>
        <span>{{ products().length }} products</span>
      </div>

      @if (isLoading()) {
        <p class="message">Loading the collection...</p>
      }
      @if (errorMessage()) {
        <p class="message error">{{ errorMessage() }}</p>
      }

      <div class="grid">
        @for (product of products(); track product.id) {
          <a class="product" [routerLink]="['/products', product.id]">
            <div class="media">
              @if (primaryPhoto(product); as photo) {
                <img [src]="photo.thumbnail_url || photo.url" [alt]="photo.alt_text || product.name" loading="lazy">
              } @else {
                <span>{{ product.category }}</span>
              }
            </div>
            <div class="product-copy">
              <div>
                <p>{{ product.category }}</p>
                <h3>{{ product.name }}</h3>
              </div>
              <strong>{{ formatMoney(product) }}</strong>
            </div>
            <small>{{ product.short_description }}</small>
          </a>
        }
      </div>
    </section>
  `,
  styles: `
    .hero { display: grid; grid-template-columns: 1.2fr .8fr; gap: 2rem; align-items: center; padding: clamp(3rem, 8vw, 7rem) clamp(1rem, 7vw, 7rem); background: #e7ecdf; }
    .eyebrow { margin: 0 0 .8rem; color: #6d7b68; font-size: .72rem; font-weight: 700; letter-spacing: .16em; text-transform: uppercase; }
    h1, h2, h3, .hero-card strong { font-family: Georgia, serif; }
    h1 { max-width: 42rem; margin: 0; color: #28432f; font-size: clamp(3.2rem, 7vw, 6.3rem); letter-spacing: -.07em; line-height: .95; }
    .intro { max-width: 36rem; margin: 1.5rem 0 2rem; color: #59685d; font-size: 1.08rem; line-height: 1.75; }
    .hero a { display: inline-flex; border-radius: 999px; background: #31543c; padding: .85rem 1.2rem; color: white; font-weight: 700; text-decoration: none; }
    .hero-card { display: grid; align-content: end; min-height: 23rem; border-radius: 1.5rem; background: linear-gradient(145deg, #c2b08f, #7e9276); padding: 2rem; color: white; box-shadow: 0 24px 50px rgb(49 84 60 / 20%); }
    .hero-card span, .hero-card small { font-weight: 700; letter-spacing: .06em; text-transform: uppercase; }
    .hero-card span { margin-bottom: auto; font-size: .7rem; }
    .hero-card strong { font-size: clamp(2rem, 4vw, 3.4rem); line-height: 1; }
    .hero-card small { margin-top: 1rem; font-size: .66rem; }
    .catalog { padding: 4.5rem clamp(1rem, 7vw, 7rem); }
    .section-heading { display: flex; align-items: end; justify-content: space-between; margin-bottom: 2rem; }
    h2 { margin: 0; color: #28432f; font-size: clamp(2.3rem, 4vw, 3.5rem); letter-spacing: -.05em; }
    .section-heading span, .message { color: #6d7b68; }
    .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(15rem, 1fr)); gap: 1.3rem; }
    .product { display: block; color: inherit; text-decoration: none; }
    .media { display: grid; aspect-ratio: 1 / 1.12; place-items: center; overflow: hidden; border-radius: 1rem; background: linear-gradient(145deg, #e2ddcf, #cdd8c3); color: #61705f; font-size: .75rem; font-weight: 700; letter-spacing: .14em; text-transform: uppercase; }
    .media img { width: 100%; height: 100%; object-fit: cover; transition: transform 240ms ease; }
    .product:hover img { transform: scale(1.04); }
    .product-copy { display: flex; align-items: start; justify-content: space-between; gap: .8rem; padding-top: .9rem; }
    .product-copy p, h3, .product small { margin: 0; }
    .product-copy p { color: #7a867b; font-size: .7rem; font-weight: 700; letter-spacing: .12em; text-transform: uppercase; }
    h3 { margin-top: .25rem; color: #314538; font-size: 1.25rem; }
    .product-copy strong { color: #31543c; white-space: nowrap; }
    .product small { display: block; margin-top: .5rem; color: #758077; line-height: 1.5; }
    .error { color: #9b3f34; }
    @media (max-width: 760px) { .hero { grid-template-columns: 1fr; } .hero-card { min-height: 15rem; } }
  `,
})
export class CatalogPage implements OnInit {
  private readonly productService = inject(ProductService);
  protected readonly products = signal<Product[]>([]);
  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');

  ngOnInit(): void {
    this.productService
      .getProducts()
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: (products) => this.products.set(products),
        error: () => this.errorMessage.set('We could not load the collection. Please try again soon.'),
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
