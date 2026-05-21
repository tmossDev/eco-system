import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { ProductSummary } from '../../../core/services/product/product.model';
import { ProductService } from '../../../core/services/product/product.service';

@Component({
  selector: 'app-product-list-page',
  imports: [RouterLink],
  template: `
    <section class="products-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">Product management</p>
          <h1>Products</h1>
          <p class="description">
            Manage the catalog for physical, digital, bundled, and service
            products from one product workspace.
          </p>
        </div>

        <a [routerLink]="['/products', '1', 'edit']" class="primary-action">
          Edit sample
        </a>
      </div>

      @if (isLoading()) {
        <p class="state-message">Loading products...</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <div class="table-card">
        <div class="table-header">
          <h2>All products</h2>
          <span>{{ products().length }} products</span>
        </div>

        <div class="table-wrap">
          <table>
            <thead>
            <tr>
              <th>Name</th>
              <th>SKU</th>
              <th>Category</th>
              <th>Price</th>
              <th>Inventory</th>
              <th>Status</th>
              <th class="actions-column">Actions</th>
            </tr>
            </thead>

            <tbody>
              @for (product of products(); track product.id) {
                <tr>
                  <td>
                    <strong>{{ product.name }}</strong>
                  </td>
                  <td>{{ product.sku }}</td>
                  <td>{{ product.category }}</td>
                  <td>{{ formatMoney(product.price_cents, product.currency) }}</td>
                  <td>{{ product.inventory_count }}</td>
                  <td>
                    <span class="status" [class]="product.status.toLowerCase()">
                      {{ product.status }}
                    </span>
                  </td>
                  <td class="actions">
                    <a [routerLink]="['/products', product.id]">View</a>
                    <a [routerLink]="['/products', product.id, 'edit']">Edit</a>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      </div>
    </section>
  `,
  styles: `
    .products-page {
      padding: 2rem;
      color: #172033;
    }

    .page-header,
    .table-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
    }

    .page-header {
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

    .primary-action {
      display: inline-flex;
      border: 0;
      border-radius: 999px;
      background: #2563eb;
      color: #ffffff;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
      cursor: pointer;
      text-decoration: none;
      white-space: nowrap;
    }

    .table-card {
      overflow: hidden;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .table-header {
      padding: 1.25rem;
      border-bottom: 1px solid #dbe3ef;
    }

    .table-header span {
      color: #56657f;
    }

    .table-wrap {
      overflow-x: auto;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 860px;
    }

    th,
    td {
      padding: 1rem 1.25rem;
      border-bottom: 1px solid #eef2f7;
      text-align: left;
    }

    th {
      color: #56657f;
      font-size: 0.8rem;
      letter-spacing: 0.04em;
      text-transform: uppercase;
    }

    tbody tr:last-child td {
      border-bottom: 0;
    }

    .status {
      display: inline-flex;
      border-radius: 999px;
      padding: 0.3rem 0.65rem;
      font-size: 0.85rem;
      font-weight: 700;
    }

    .active {
      background: #dcfce7;
      color: #166534;
    }

    .draft {
      background: #fef3c7;
      color: #92400e;
    }

    .archived {
      background: #fee2e2;
      color: #991b1b;
    }

    .actions-column {
      width: 10rem;
    }

    .actions {
      display: flex;
      gap: 0.75rem;
    }

    .actions a {
      color: #2563eb;
      font-weight: 700;
      text-decoration: none;
    }

    @media (max-width: 700px) {
      .products-page {
        padding: 1rem;
      }

      .page-header {
        display: grid;
      }
    }
  `,
})
export class ProductListPage implements OnInit {
  private readonly productService = inject(ProductService);

  protected readonly isLoading = signal(true);
  protected readonly errorMessage = signal('');
  protected readonly products = signal<ProductSummary[]>([]);

  public ngOnInit(): void {
    this.productService
      .getProducts()
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (products) => {
          this.products.set(products);
        },
        error: () => {
          this.errorMessage.set('Unable to load products.');
        },
      });
  }

  protected formatMoney(priceCents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(priceCents / 100);
  }
}
