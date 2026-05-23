import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import {
  CreateDiscountRequest,
  Discount,
  ProductPhoto,
  ProductSummary,
} from '../../../core/services/product/product.model';
import {
  calculateDiscountedPrice,
  formatDiscountValue,
} from '../../../core/services/product/product-pricing';
import { ProductService } from '../../../core/services/product/product.service';

@Component({
  selector: 'app-product-list-page',
  imports: [FormsModule, RouterLink],
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

      @if (message()) {
        <p class="success-message">{{ message() }}</p>
      }

      <section class="discount-panel" aria-label="Bulk discount controls">
        <div class="discount-panel-header">
          <div>
            <h2>Discount tools</h2>
            <p>
              Apply a discount to every product, a category group, or the
              selected rows.
            </p>
          </div>
          <span>{{ selectedProductIds().length }} selected</span>
        </div>

        <div class="discount-controls">
          <label>
            <span>Apply to</span>
            <select
              name="bulk_scope"
              [ngModel]="bulkScope()"
              (ngModelChange)="setBulkScope($event)"
            >
              <option value="all">All products</option>
              <option value="category">Category group</option>
              <option value="selected">Selected products</option>
            </select>
          </label>

          @if (bulkScope() === 'category') {
            <label>
              <span>Category</span>
              <select
                name="bulk_category"
                [ngModel]="bulkCategory()"
                (ngModelChange)="bulkCategory.set($event)"
              >
                @for (category of categories(); track category) {
                  <option [value]="category">{{ category }}</option>
                }
              </select>
            </label>
          }

          <label>
            <span>Discount type</span>
            <select
              name="bulk_discount_type"
              [ngModel]="bulkDiscountForm().discount_type"
              (ngModelChange)="setBulkDiscountType($event)"
            >
              <option>Percentage</option>
              <option>Amount</option>
            </select>
          </label>

          <label>
            <span>{{ bulkDiscountForm().discount_type === 'Percentage' ? 'Percent' : 'Amount cents' }}</span>
            <input
              type="number"
              min="1"
              name="bulk_discount_value"
              [ngModel]="bulkDiscountValue()"
              (ngModelChange)="updateBulkDiscountValue($event)"
            />
          </label>

          <label>
            <span>Status</span>
            <select
              name="bulk_discount_status"
              [ngModel]="bulkDiscountForm().status"
              (ngModelChange)="updateBulkDiscountForm('status', $event)"
            >
              <option>Active</option>
              <option>Draft</option>
              <option>Archived</option>
            </select>
          </label>
        </div>

        <label>
          <span>Discount name</span>
          <input
            type="text"
            name="bulk_discount_name"
            [ngModel]="bulkDiscountForm().name"
            (ngModelChange)="updateBulkDiscountForm('name', $event)"
          />
        </label>

        <button
          type="button"
          class="secondary-action discount-action"
          [disabled]="isSavingDiscount() || targetProducts().length === 0"
          (click)="createBulkDiscount()"
        >
          {{ isSavingDiscount() ? 'Applying...' : 'Apply discount' }}
        </button>
      </section>

      <div class="table-card">
        <div class="table-header">
          <h2>All products</h2>
          <span>{{ products().length }} products</span>
        </div>

        <div class="table-wrap">
          <table>
            <thead>
            <tr>
              <th class="select-column">
                <input
                  type="checkbox"
                  aria-label="Select all products"
                  [checked]="allProductsSelected()"
                  (change)="toggleAllProducts($event)"
                />
              </th>
              <th>Name</th>
              <th>SKU</th>
              <th>Category</th>
              <th>Price</th>
              <th>Discounts</th>
              <th>Inventory</th>
              <th>Status</th>
              <th class="actions-column">Actions</th>
            </tr>
            </thead>

            <tbody>
              @for (product of products(); track product.id) {
                <tr>
                  <td>
                    <input
                      type="checkbox"
                      [attr.aria-label]="'Select ' + product.name"
                      [checked]="isProductSelected(product.id)"
                      (change)="toggleProduct(product.id, $event)"
                    />
                  </td>
                  <td>
                    <div class="product-cell">
                      @if (primaryPhoto(product)) {
                        <img
                          [src]="photoPreviewUrl(primaryPhoto(product)!)"
                          [alt]="primaryPhoto(product)!.alt_text || product.name"
                          loading="lazy"
                          decoding="async"
                          width="44"
                          height="44"
                        />
                      } @else {
                        <span class="photo-fallback">
                          {{ product.sku.slice(0, 2) }}
                        </span>
                      }

                      <span>
                        <strong>{{ product.name }}</strong>
                        <small>{{ product.short_description }}</small>
                      </span>
                    </div>
                  </td>
                  <td>{{ product.sku }}</td>
                  <td>{{ product.category }}</td>
                  <td>
                    <span class="price-stack">
                      @if (discountedPrice(product).hasDiscount) {
                        <strong>
                          {{ formatMoney(discountedPrice(product).finalCents, product.currency) }}
                        </strong>
                        <span>{{ formatMoney(product.price_cents, product.currency) }}</span>
                      } @else {
                        <strong>{{ formatMoney(product.price_cents, product.currency) }}</strong>
                      }
                    </span>
                  </td>
                  <td>{{ discountSummary(product) }}</td>
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

    .primary-action,
    .secondary-action {
      display: inline-flex;
      border: 0;
      border-radius: 999px;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
      cursor: pointer;
      text-decoration: none;
      white-space: nowrap;
    }

    .primary-action {
      background: #2563eb;
      color: #ffffff;
    }

    .secondary-action {
      border: 1px solid #cbd5e1;
      background: #ffffff;
      color: #2563eb;
    }

    .secondary-action:disabled {
      cursor: default;
      opacity: 0.65;
    }

    .discount-panel {
      display: grid;
      gap: 1rem;
      margin-bottom: 1.25rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.25rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .discount-panel-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
    }

    .discount-panel-header p {
      margin: 0.35rem 0 0;
      color: #56657f;
      line-height: 1.5;
    }

    .discount-panel-header span {
      color: #56657f;
      font-weight: 700;
      white-space: nowrap;
    }

    .discount-controls {
      display: grid;
      grid-template-columns: repeat(5, minmax(8rem, 1fr));
      gap: 0.85rem;
    }

    .discount-panel label {
      display: grid;
      gap: 0.4rem;
      font-weight: 700;
    }

    .discount-panel input,
    .discount-panel select {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.75rem 0.85rem;
      color: #172033;
      font: inherit;
    }

    .discount-action {
      width: fit-content;
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
      min-width: 1040px;
    }

    th,
    td {
      padding: 1rem 1.25rem;
      border-bottom: 1px solid #eef2f7;
      text-align: left;
    }

    .select-column {
      width: 3rem;
    }

    td input[type='checkbox'],
    th input[type='checkbox'] {
      width: 1rem;
      height: 1rem;
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

    .product-cell {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      min-width: 18rem;
    }

    .product-cell img,
    .photo-fallback {
      flex: 0 0 auto;
      width: 2.75rem;
      height: 2.75rem;
      border-radius: 0.5rem;
    }

    .product-cell img {
      object-fit: cover;
      background: #eef2f7;
    }

    .photo-fallback {
      display: grid;
      place-items: center;
      background: #dbeafe;
      color: #1d4ed8;
      font-size: 0.8rem;
      font-weight: 800;
    }

    .product-cell span:last-child {
      display: grid;
      gap: 0.2rem;
      min-width: 0;
    }

    .product-cell small {
      max-width: 24rem;
      overflow: hidden;
      color: #56657f;
      font-size: 0.85rem;
      font-weight: 500;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .price-stack {
      display: inline-grid;
      gap: 0.15rem;
      white-space: nowrap;
    }

    .price-stack span {
      color: #64748b;
      font-size: 0.85rem;
      text-decoration: line-through;
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

      .discount-panel-header,
      .discount-controls {
        display: grid;
        grid-template-columns: 1fr;
      }
    }
  `,
})
export class ProductListPage implements OnInit {
  private readonly productService = inject(ProductService);

  protected readonly isLoading = signal(true);
  protected readonly isSavingDiscount = signal(false);
  protected readonly message = signal('');
  protected readonly errorMessage = signal('');
  protected readonly products = signal<ProductSummary[]>([]);
  protected readonly selectedProductIds = signal<string[]>([]);
  protected readonly bulkScope = signal<'all' | 'category' | 'selected'>('selected');
  protected readonly bulkCategory = signal('');
  protected readonly bulkDiscountForm = signal<CreateDiscountRequest>({
    name: 'Catalog discount',
    description: '',
    discount_type: 'Percentage',
    scope: 'ProductSet',
    percentage_basis_points: 1000,
    amount_cents: null,
    currency: '',
    min_product_count: 1,
    starts_at: '',
    ends_at: '',
    status: 'Active',
    product_ids: [],
  });

  protected readonly categories = computed(() => [
    ...new Set(this.products().map((product) => product.category).filter(Boolean)),
  ]);

  protected readonly targetProducts = computed(() => {
    if (this.bulkScope() === 'all') {
      return this.products();
    }

    if (this.bulkScope() === 'category') {
      return this.products().filter(
        (product) => product.category === this.bulkCategory(),
      );
    }

    const selected = new Set(this.selectedProductIds());

    return this.products().filter((product) => selected.has(String(product.id)));
  });

  public ngOnInit(): void {
    this.loadProducts();
  }

  protected formatMoney(priceCents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(priceCents / 100);
  }

  protected discountedPrice(product: ProductSummary) {
    return calculateDiscountedPrice(product.price_cents, product.discounts);
  }

  protected discountSummary(product: ProductSummary): string {
    const discounts = product.discounts ?? [];
    if (discounts.length === 0) {
      return 'None';
    }

    return discounts
      .map((discount) => {
        return formatDiscountValue(discount, this.formatMoney.bind(this));
      })
      .join(', ');
  }

  protected allProductsSelected(): boolean {
    return (
      this.products().length > 0 &&
      this.products().every((product) => this.isProductSelected(product.id))
    );
  }

  protected isProductSelected(id: string): boolean {
    return this.selectedProductIds().includes(String(id));
  }

  protected toggleProduct(id: string, event: Event): void {
    const checked = (event.target as HTMLInputElement).checked;
    const productId = String(id);

    this.selectedProductIds.update((ids) =>
      checked ? [...new Set([...ids, productId])] : ids.filter((id) => id !== productId),
    );
  }

  protected toggleAllProducts(event: Event): void {
    const checked = (event.target as HTMLInputElement).checked;

    this.selectedProductIds.set(
      checked ? this.products().map((product) => String(product.id)) : [],
    );
  }

  protected setBulkScope(scope: 'all' | 'category' | 'selected'): void {
    this.bulkScope.set(scope);

    if (scope === 'category' && !this.bulkCategory()) {
      this.bulkCategory.set(this.categories()[0] ?? '');
    }
  }

  protected updateBulkDiscountForm<Key extends keyof CreateDiscountRequest>(
    key: Key,
    value: CreateDiscountRequest[Key],
  ): void {
    this.bulkDiscountForm.update((form) => ({
      ...form,
      [key]: value,
    }));
  }

  protected bulkDiscountValue(): number {
    const form = this.bulkDiscountForm();

    return form.discount_type === 'Percentage'
      ? (form.percentage_basis_points ?? 0) / 100
      : form.amount_cents ?? 0;
  }

  protected updateBulkDiscountValue(value: number | string): void {
    const numericValue = Number(value) || 0;

    this.bulkDiscountForm.update((form) => ({
      ...form,
      percentage_basis_points:
        form.discount_type === 'Percentage' ? Math.round(numericValue * 100) : null,
      amount_cents: form.discount_type === 'Amount' ? Math.round(numericValue) : null,
      currency: form.discount_type === 'Amount' ? this.defaultCurrency() : '',
    }));
  }

  protected setBulkDiscountType(
    discountType: CreateDiscountRequest['discount_type'],
  ): void {
    this.bulkDiscountForm.update((form) => ({
      ...form,
      discount_type: discountType,
      percentage_basis_points:
        discountType === 'Percentage' ? form.percentage_basis_points ?? 1000 : null,
      amount_cents: discountType === 'Amount' ? form.amount_cents ?? 500 : null,
      currency: discountType === 'Amount' ? this.defaultCurrency() : '',
    }));
  }

  protected createBulkDiscount(): void {
    if (this.isSavingDiscount() || this.targetProducts().length === 0) {
      return;
    }

    this.isSavingDiscount.set(true);
    this.message.set('');
    this.errorMessage.set('');

    const isGlobal = this.bulkScope() === 'all';
    const request: CreateDiscountRequest = {
      ...this.bulkDiscountForm(),
      scope: isGlobal ? 'Global' : 'ProductSet',
      currency:
        this.bulkDiscountForm().discount_type === 'Amount'
          ? this.defaultCurrency()
          : '',
      product_ids: isGlobal
        ? []
        : this.targetProducts().map((product) => this.normalizedProductId(product.id)),
    };

    this.productService
      .createDiscount(request)
      .pipe(
        finalize(() => {
          this.isSavingDiscount.set(false);
        }),
      )
      .subscribe({
        next: () => {
          this.message.set('Discount applied successfully.');
          this.loadProducts(false);
        },
        error: () => {
          this.errorMessage.set('Unable to apply discount.');
        },
      });
  }

  protected primaryPhoto(product: ProductSummary): ProductPhoto | null {
    const photos = product.photos ?? [];

    return (
      photos.find((photo) => photo.is_primary) ?? photos[0] ?? null
    );
  }

  protected photoPreviewUrl(photo: ProductPhoto): string {
    return photo.thumbnail_url || photo.url;
  }

  private loadProducts(showLoading = true): void {
    if (showLoading) {
      this.isLoading.set(true);
    }

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
          this.selectedProductIds.update((ids) =>
            ids.filter((id) =>
              products.some((product) => String(product.id) === String(id)),
            ),
          );

          if (!this.bulkCategory()) {
            this.bulkCategory.set(this.categories()[0] ?? '');
          }
        },
        error: () => {
          this.errorMessage.set('Unable to load products.');
        },
      });
  }

  private defaultCurrency(): string {
    return this.targetProducts()[0]?.currency ?? this.products()[0]?.currency ?? 'USD';
  }

  private normalizedProductId(id: string): number | string {
    const numericId = Number(id);

    return Number.isFinite(numericId) ? numericId : id;
  }
}
