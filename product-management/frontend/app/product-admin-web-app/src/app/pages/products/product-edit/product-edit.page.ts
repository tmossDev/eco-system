import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { finalize } from 'rxjs';

import {
  CreateDiscountRequest,
  Discount,
  ProductPhoto,
  UpdateProductRequest,
} from '../../../core/services/product/product.model';
import {
  calculateDiscountedPrice,
  formatDiscountValue,
} from '../../../core/services/product/product-pricing';
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
          <span>Short description</span>
          <textarea
            name="short_description"
            rows="2"
            maxlength="280"
            [ngModel]="form().short_description"
            (ngModelChange)="updateForm('short_description', $event)"
          ></textarea>
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

        <section class="photo-editor" aria-label="Product photos">
          <div class="section-header">
            <div>
              <h2>Photos</h2>
            </div>
            <label class="upload-button">
              <input
                type="file"
                accept="image/png,image/jpeg"
                [disabled]="isUploading()"
                (change)="uploadPhoto($event)"
              />
              <span>{{ isUploading() ? 'Uploading...' : 'Upload photo' }}</span>
            </label>
          </div>

          @if (form().photos.length === 0) {
            <p class="empty-state">No photos added.</p>
          }

          <div class="photo-list">
            @for (photo of form().photos; track $index) {
              <div class="photo-row">
                <div class="photo-preview">
                  @if (photo.thumbnail_url || photo.url) {
                    <img
                      [src]="photo.thumbnail_url || photo.url"
                      [alt]="photo.alt_text || form().name || 'Product photo'"
                      loading="lazy"
                      decoding="async"
                      width="72"
                      height="72"
                    />
                  } @else {
                    <span>{{ $index + 1 }}</span>
                  }
                </div>

                <div class="photo-fields">
                  <label>
                    <span>Alt text</span>
                    <input
                      type="text"
                      [name]="'photo_alt_text_' + $index"
                      [ngModel]="photo.alt_text"
                      (ngModelChange)="updatePhoto($index, 'alt_text', $event)"
                    />
                  </label>
                </div>

                <div class="photo-actions">
                  <label class="checkbox-label">
                    <input
                      type="checkbox"
                      [name]="'photo_primary_' + $index"
                      [ngModel]="photo.is_primary"
                      (ngModelChange)="setPrimaryPhoto($index, $event)"
                    />
                    <span>Primary</span>
                  </label>

                  <button
                    type="button"
                    class="text-button"
                    (click)="removePhoto($index)"
                  >
                    Remove
                  </button>
                </div>
              </div>
            }
          </div>
        </section>

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

        <section class="discount-editor" aria-label="Product discounts">
          <div class="section-header">
            <div>
              <h2>Discounts</h2>
              <p class="section-note">
                Final price:
                <strong>{{ formatMoney(discountedPrice().finalCents, form().currency) }}</strong>
                @if (discountedPrice().hasDiscount) {
                  <span>{{ formatMoney(form().price_cents, form().currency) }}</span>
                }
              </p>
            </div>
          </div>

          @if (discounts().length === 0) {
            <p class="empty-state">No discounts are attached to this product.</p>
          }

          <div class="discount-list">
            @for (discount of discounts(); track discount.id) {
              <div class="discount-row">
                <div>
                  <strong>{{ discount.name }}</strong>
                  <span>{{ formatDiscount(discount) }} · {{ discount.status }}</span>
                </div>
                <button
                  type="button"
                  class="text-button"
                  [disabled]="isSavingDiscount()"
                  (click)="toggleDiscountStatus(discount)"
                >
                  {{ discount.status === 'Active' ? 'Archive' : 'Activate' }}
                </button>
              </div>
            }
          </div>

          <div class="discount-form">
            <label>
              <span>Name</span>
              <input
                type="text"
                name="discount_name"
                [ngModel]="discountForm().name"
                (ngModelChange)="updateDiscountForm('name', $event)"
              />
            </label>

            <div class="field-grid discount-grid">
              <label>
                <span>Type</span>
                <select
                  name="discount_type"
                  [ngModel]="discountForm().discount_type"
                  (ngModelChange)="setDiscountType($event)"
                >
                  <option>Percentage</option>
                  <option>Amount</option>
                </select>
              </label>

              <label>
                <span>{{ discountForm().discount_type === 'Percentage' ? 'Percent' : 'Amount cents' }}</span>
                <input
                  type="number"
                  min="1"
                  name="discount_value"
                  [ngModel]="discountValue()"
                  (ngModelChange)="updateDiscountValue($event)"
                />
              </label>

              <label>
                <span>Status</span>
                <select
                  name="discount_status"
                  [ngModel]="discountForm().status"
                  (ngModelChange)="updateDiscountForm('status', $event)"
                >
                  <option>Active</option>
                  <option>Draft</option>
                  <option>Archived</option>
                </select>
              </label>
            </div>

            <button
              type="button"
              class="secondary-button"
              [disabled]="isSavingDiscount()"
              (click)="createProductDiscount()"
            >
              {{ isSavingDiscount() ? 'Saving discount...' : 'Add product discount' }}
            </button>
          </div>
        </section>

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

    h2 {
      margin: 0;
      font-size: 1rem;
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

    .photo-editor {
      display: grid;
      gap: 0.9rem;
      border-top: 1px solid #eef2f7;
      border-bottom: 1px solid #eef2f7;
      padding: 1rem 0;
    }

    .discount-editor {
      display: grid;
      gap: 1rem;
      border-top: 1px solid #eef2f7;
      border-bottom: 1px solid #eef2f7;
      padding: 1rem 0;
    }

    .section-note {
      margin: 0.35rem 0 0;
      color: #56657f;
    }

    .section-note span {
      margin-left: 0.45rem;
      color: #64748b;
      text-decoration: line-through;
    }

    .discount-list,
    .discount-form {
      display: grid;
      gap: 0.7rem;
    }

    .discount-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 0.75rem;
      padding: 0.75rem;
    }

    .discount-row div {
      display: grid;
      gap: 0.2rem;
    }

    .discount-row span {
      color: #56657f;
    }

    .section-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
    }

    .empty-state {
      margin: 0.35rem 0 0;
      color: #56657f;
      line-height: 1.5;
    }

    .photo-list {
      display: grid;
      gap: 0.75rem;
    }

    .photo-row {
      display: grid;
      grid-template-columns: 4.5rem 1fr auto;
      gap: 0.9rem;
      border: 1px solid #dbe3ef;
      border-radius: 0.75rem;
      padding: 0.75rem;
    }

    .photo-preview,
    .photo-preview img {
      width: 4.5rem;
      height: 4.5rem;
      border-radius: 0.5rem;
    }

    .photo-preview {
      display: grid;
      place-items: center;
      overflow: hidden;
      background: #eef2f7;
      color: #56657f;
      font-weight: 800;
    }

    .photo-preview img {
      object-fit: cover;
    }

    .photo-fields {
      display: grid;
      gap: 0.65rem;
    }

    .photo-actions {
      display: flex;
      align-items: flex-start;
      flex-direction: column;
      gap: 0.75rem;
      min-width: 6rem;
    }

    .checkbox-label {
      display: inline-flex;
      align-items: center;
      gap: 0.45rem;
      font-weight: 700;
      white-space: nowrap;
    }

    .checkbox-label input {
      width: auto;
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
    .secondary-action,
    .secondary-button,
    .upload-button,
    .text-button {
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

    .secondary-button {
      width: fit-content;
      border: 1px solid #cbd5e1;
      background: #ffffff;
      color: #2563eb;
      cursor: pointer;
    }

    .upload-button,
    .text-button {
      border: 1px solid #cbd5e1;
      background: #ffffff;
      color: #2563eb;
      cursor: pointer;
    }

    .text-button {
      border: 0;
      padding: 0;
      background: transparent;
    }

    .upload-button {
      position: relative;
      display: inline-flex;
      overflow: hidden;
      white-space: nowrap;
    }

    .upload-button input {
      position: absolute;
      inset: 0;
      opacity: 0;
      cursor: pointer;
    }

    .upload-button:has(input:disabled) {
      cursor: default;
      opacity: 0.65;
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

      .discount-row {
        align-items: flex-start;
        flex-direction: column;
      }

      .section-header,
      .photo-row {
        grid-template-columns: 1fr;
      }

      .section-header {
        display: grid;
      }

      .photo-actions {
        flex-direction: row;
        justify-content: space-between;
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
  protected readonly isUploading = signal(false);
  protected readonly isSavingDiscount = signal(false);
  protected readonly message = signal('');
  protected readonly errorMessage = signal('');
  protected readonly discounts = signal<Discount[]>([]);

  protected readonly productId = computed(
    () => this.route.snapshot.paramMap.get('id') ?? '1',
  );

  protected readonly form = signal<UpdateProductRequest>({
    sku: '',
    name: '',
    description: '',
    short_description: '',
    category: 'General',
    price_cents: 0,
    currency: 'USD',
    inventory_count: 0,
    status: 'Draft',
    photos: [],
  });

  protected readonly discountForm = signal<CreateDiscountRequest>({
    name: 'Product discount',
    description: '',
    discount_type: 'Percentage',
    scope: 'ProductSet',
    percentage_basis_points: 1000,
    amount_cents: null,
    currency: 'USD',
    min_product_count: 1,
    starts_at: '',
    ends_at: '',
    status: 'Active',
    product_ids: [],
  });

  protected readonly discountedPrice = computed(() =>
    calculateDiscountedPrice(this.form().price_cents, this.discounts()),
  );

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
            short_description: product.short_description ?? '',
            description: product.description,
            category: product.category,
            price_cents: product.price_cents,
            currency: product.currency,
            inventory_count: product.inventory_count,
            status: product.status,
            photos: product.photos ?? [],
          });
          this.discounts.set(product.discounts ?? []);
          this.discountForm.update((form) => ({
            ...form,
            currency: product.currency,
            product_ids: [this.normalizedProductId()],
          }));
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

  protected updateDiscountForm<Key extends keyof CreateDiscountRequest>(
    key: Key,
    value: CreateDiscountRequest[Key],
  ): void {
    this.discountForm.update((form) => ({
      ...form,
      [key]: value,
    }));
  }

  protected discountValue(): number {
    const form = this.discountForm();

    return form.discount_type === 'Percentage'
      ? (form.percentage_basis_points ?? 0) / 100
      : form.amount_cents ?? 0;
  }

  protected updateDiscountValue(value: number | string): void {
    const numericValue = Number(value) || 0;

    this.discountForm.update((form) => ({
      ...form,
      percentage_basis_points:
        form.discount_type === 'Percentage' ? Math.round(numericValue * 100) : null,
      amount_cents: form.discount_type === 'Amount' ? Math.round(numericValue) : null,
      currency: form.discount_type === 'Amount' ? this.form().currency : '',
    }));
  }

  protected setDiscountType(discountType: CreateDiscountRequest['discount_type']): void {
    this.discountForm.update((form) => ({
      ...form,
      discount_type: discountType,
      percentage_basis_points:
        discountType === 'Percentage' ? form.percentage_basis_points ?? 1000 : null,
      amount_cents: discountType === 'Amount' ? form.amount_cents ?? 500 : null,
      currency: discountType === 'Amount' ? this.form().currency : '',
    }));
  }

  protected formatMoney(priceCents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(priceCents / 100);
  }

  protected formatDiscount(discount: Discount): string {
    return formatDiscountValue(discount, this.formatMoney.bind(this));
  }

  protected updatePhoto<Key extends keyof ProductPhoto>(
    index: number,
    key: Key,
    value: ProductPhoto[Key],
  ): void {
    const photos = this.form().photos.map((photo, photoIndex) =>
      photoIndex === index ? { ...photo, [key]: value } : photo,
    );

    this.updateForm('photos', photos);
  }

  protected setPrimaryPhoto(index: number, isPrimary: boolean): void {
    const photos = this.form().photos.map((photo, photoIndex) => ({
      ...photo,
      is_primary: isPrimary ? photoIndex === index : false,
    }));

    this.updateForm('photos', photos);
  }

  protected removePhoto(index: number): void {
    const photos = this.form().photos.filter((_, photoIndex) => photoIndex !== index);

    if (photos.length > 0 && !photos.some((photo) => photo.is_primary)) {
      photos[0] = {
        ...photos[0],
        is_primary: true,
      };
    }

    this.updateForm('photos', photos);
  }

  protected uploadPhoto(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    input.value = '';

    if (!file || this.isUploading()) {
      return;
    }

    this.isUploading.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.productService
      .uploadProductPhoto(this.productId(), file)
      .pipe(
        finalize(() => {
          this.isUploading.set(false);
        }),
      )
      .subscribe({
        next: (product) => {
          this.form.update((form) => ({
            ...form,
            photos: product.photos ?? [],
          }));
          this.message.set('Photo uploaded successfully.');
        },
        error: () => {
          this.errorMessage.set('Unable to upload photo.');
        },
      });
  }

  protected createProductDiscount(): void {
    if (this.isSavingDiscount()) {
      return;
    }

    this.isSavingDiscount.set(true);
    this.message.set('');
    this.errorMessage.set('');

    const request: CreateDiscountRequest = {
      ...this.discountForm(),
      currency:
        this.discountForm().discount_type === 'Amount' ? this.form().currency : '',
      product_ids: [this.normalizedProductId()],
    };

    this.productService
      .createDiscount(request)
      .pipe(
        finalize(() => {
          this.isSavingDiscount.set(false);
        }),
      )
      .subscribe({
        next: (discount) => {
          this.discounts.update((discounts) => [discount, ...discounts]);
          this.message.set('Discount added successfully.');
        },
        error: () => {
          this.errorMessage.set('Unable to add discount.');
        },
      });
  }

  protected toggleDiscountStatus(discount: Discount): void {
    if (this.isSavingDiscount()) {
      return;
    }

    this.isSavingDiscount.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.productService
      .updateDiscount(discount.id, {
        ...discount,
        status: discount.status === 'Active' ? 'Archived' : 'Active',
        product_ids: discount.product_ids,
      })
      .pipe(
        finalize(() => {
          this.isSavingDiscount.set(false);
        }),
      )
      .subscribe({
        next: (updatedDiscount) => {
          this.discounts.update((discounts) =>
            discounts.map((candidate) =>
              candidate.id === updatedDiscount.id ? updatedDiscount : candidate,
            ),
          );
          this.message.set('Discount updated successfully.');
        },
        error: () => {
          this.errorMessage.set('Unable to update discount.');
        },
      });
  }

  protected saveProduct(): void {
    if (this.isSaving()) {
      return;
    }

    this.isSaving.set(true);
    this.message.set('');
    this.errorMessage.set('');

    const request: UpdateProductRequest = {
      ...this.form(),
      photos: this.form()
        .photos.filter((photo) => photo.url.trim() !== '')
        .map((photo) => ({
          ...photo,
          thumbnail_url: photo.thumbnail_url || photo.url,
        })),
    };

    this.productService
      .updateProduct(this.productId(), request)
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

  private normalizedProductId(): number | string {
    const numericId = Number(this.productId());

    return Number.isFinite(numericId) ? numericId : this.productId();
  }
}
