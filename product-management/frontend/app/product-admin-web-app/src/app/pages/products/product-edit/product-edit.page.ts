import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import { UpdateProductRequest } from '../../../core/services/product/product.model';
import { ProductService } from '../../../core/services/product/product.service';

@Component({
  selector: 'app-product-edit-page',
  imports: [FormsModule, RouterLink],
  template: `
    <section class="product-edit-page">
      <a [routerLink]="['/products', productId()]" class="back-link">
        ← Back to profile
      </a>

      @if (isLoading()) {
        <p class="state-message">Loading product...</p>
      }

      @if (message()) {
        <p class="success-message">{{ message() }}</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <form class="form-card" (ngSubmit)="saveProduct()">
        <div class="form-header">
          <p class="eyebrow">Edit product</p>
          <h1>Product details</h1>
          <p>
            Update catalog information, pricing, inventory, and publication
            status.
          </p>
        </div>

        <label>
          <span>SKU</span>
          <input
            type="text"
            name="sku"
            [ngModel]="form().sku"
            (ngModelChange)="updateForm('sku', $event)"
          />
        </label>

        <label>
          <span>Name</span>
          <input
            type="text"
            name="name"
            [ngModel]="form().name"
            (ngModelChange)="updateForm('name', $event)"
          />
        </label>

        <label>
          <span>Description</span>
          <textarea
            name="description"
            rows="4"
            [ngModel]="form().description"
            (ngModelChange)="updateForm('description', $event)"
          ></textarea>
        </label>

        <label>
          <span>Category</span>
          <input
            type="text"
            name="category"
            [ngModel]="form().category"
            (ngModelChange)="updateForm('category', $event)"
          />
        </label>

        <div class="field-grid">
          <label>
            <span>Price cents</span>
            <input
              type="number"
              min="0"
              name="price_cents"
              [ngModel]="form().price_cents"
              (ngModelChange)="updateForm('price_cents', $event)"
            />
          </label>

          <label>
            <span>Currency</span>
            <input
              type="text"
              maxlength="3"
              name="currency"
              [ngModel]="form().currency"
              (ngModelChange)="updateForm('currency', $event)"
            />
          </label>

          <label>
            <span>Inventory</span>
            <input
              type="number"
              min="0"
              name="inventory_count"
              [ngModel]="form().inventory_count"
              (ngModelChange)="updateForm('inventory_count', $event)"
            />
          </label>
        </div>

        <label>
          <span>Status</span>
          <select
            name="status"
            [ngModel]="form().status"
            (ngModelChange)="updateForm('status', $event)"
          >
            <option>Draft</option>
            <option>Active</option>
            <option>Archived</option>
          </select>
        </label>

        <div class="form-actions">
          <button type="submit" class="primary-action" [disabled]="isSaving()">
            {{ isSaving() ? 'Saving...' : 'Save changes' }}
          </button>
          <a [routerLink]="['/products', productId()]" class="secondary-action">
            Cancel
          </a>
        </div>
      </form>
    </section>
  `,
  styles: `
    .product-edit-page {
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
    .success-message,
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

    .success-message {
      background: #dcfce7;
      color: #166534;
    }

    .error-message {
      background: #fee2e2;
      color: #991b1b;
    }

    .form-card {
      display: grid;
      gap: 1rem;
      max-width: 42rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.5rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .form-header {
      margin-bottom: 0.5rem;
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

    .form-header p {
      margin: 0.75rem 0 0;
      color: #56657f;
      line-height: 1.6;
    }

    label {
      display: grid;
      gap: 0.4rem;
      font-weight: 700;
    }

    input,
    select,
    textarea {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.8rem 0.9rem;
      color: #172033;
      font: inherit;
    }

    textarea {
      resize: vertical;
    }

    input:focus,
    select:focus,
    textarea:focus {
      border-color: #2563eb;
      outline: 3px solid #dbeafe;
    }

    .field-grid {
      display: grid;
      grid-template-columns: 1fr 10rem 1fr;
      gap: 1rem;
    }

    .form-actions {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-top: 0.5rem;
    }

    .primary-action,
    .secondary-action {
      border-radius: 999px;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
    }

    .primary-action {
      border: 0;
      background: #2563eb;
      color: #ffffff;
      cursor: pointer;
    }

    .primary-action:disabled {
      cursor: default;
      opacity: 0.65;
    }

    .secondary-action {
      color: #2563eb;
      text-decoration: none;
    }

    @media (max-width: 640px) {
      .product-edit-page {
        padding: 1rem;
      }

      .form-actions {
        align-items: stretch;
        flex-direction: column;
      }

      .field-grid {
        grid-template-columns: 1fr;
      }

      .primary-action,
      .secondary-action {
        text-align: center;
      }
    }
  `,
})
export class ProductEditPage implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly productService = inject(ProductService);

  protected readonly isLoading = signal(true);
  protected readonly isSaving = signal(false);
  protected readonly message = signal('');
  protected readonly errorMessage = signal('');

  protected readonly productId = computed(
    () => this.route.snapshot.paramMap.get('id') ?? '1',
  );

  protected readonly form = signal<UpdateProductRequest>({
    sku: '',
    name: '',
    description: '',
    category: 'General',
    price_cents: 0,
    currency: 'USD',
    inventory_count: 0,
    status: 'Draft',
  });

  public ngOnInit(): void {
    this.productService
      .getProductById(this.productId())
      .pipe(
        finalize(() => {
          this.isLoading.set(false);
        }),
      )
      .subscribe({
        next: (product) => {
          this.form.set({
            sku: product.sku,
            name: product.name,
            description: product.description,
            category: product.category,
            price_cents: product.price_cents,
            currency: product.currency,
            inventory_count: product.inventory_count,
            status: product.status,
          });
        },
        error: () => {
          this.errorMessage.set('Unable to load product.');
        },
      });
  }

  protected updateForm<Key extends keyof UpdateProductRequest>(
    key: Key,
    value: UpdateProductRequest[Key],
  ): void {
    this.form.update((form) => ({
      ...form,
      [key]: value,
    }));
  }

  protected saveProduct(): void {
    if (this.isSaving()) {
      return;
    }

    this.isSaving.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.productService
      .updateProduct(this.productId(), this.form())
      .pipe(
        finalize(() => {
          this.isSaving.set(false);
        }),
      )
      .subscribe({
        next: (product) => {
          this.message.set('Product saved successfully.');
          void this.router.navigate(['/products', product.id]);
        },
        error: () => {
          this.errorMessage.set('Unable to save product.');
        },
      });
  }
}
