import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { finalize, forkJoin } from 'rxjs';

import {
  CreateDiscountRequest,
  Discount,
  DiscountType,
  formatDiscountValue,
  PromotionSettings,
  ProductService,
  ProductSummary,
} from '@eco/admin-features';

type PromotionTarget = 'all' | 'category' | 'label' | 'selected';

@Component({
  selector: 'app-promotions-page',
  imports: [FormsModule],
  template: `
    <section class="promotions-page">
      <div class="page-header">
        <div>
          <p class="eyebrow">Promotions</p>
          <h1>Sales and discounts</h1>
          <p class="description">
            Build scheduled promotions, target product labels or groups, and
            keep active and inactive offers visible in one workspace.
          </p>
        </div>
      </div>

      @if (isLoading()) {
        <p class="state-message">Loading promotions...</p>
      }

      @if (message()) {
        <p class="success-message">{{ message() }}</p>
      }

      @if (errorMessage()) {
        <p class="error-message">{{ errorMessage() }}</p>
      }

      <section class="builder">
        <div class="section-header">
          <div>
            <h2>Promotion builder</h2>
            <p>{{ targetProducts().length }} products targeted</p>
          </div>
          <button
            type="button"
            class="primary-action"
            [disabled]="isSaving() || targetProducts().length === 0"
            (click)="createPromotion()"
          >
            {{ isSaving() ? 'Saving...' : 'Create promotion' }}
          </button>
        </div>

        <div class="global-toggle">
          <div>
            <strong>Global promotions</strong>
            <span>{{ promotionSettings().promotions_enabled ? 'Promotions are applied' : 'All promotions are paused' }}</span>
          </div>
          <label>
            <input
              type="checkbox"
              [checked]="promotionSettings().promotions_enabled"
              [disabled]="isSaving()"
              (change)="toggleGlobalPromotions($event)"
            />
            <span>Enabled</span>
          </label>
        </div>

        <div class="form-grid">
          <label class="wide">
            <span>Name</span>
            <input
              type="text"
              name="promotion_name"
              [ngModel]="form().name"
              (ngModelChange)="updateForm('name', $event)"
            />
          </label>

          <label class="wide">
            <span>Description</span>
            <input
              type="text"
              name="promotion_description"
              [ngModel]="form().description"
              (ngModelChange)="updateForm('description', $event)"
            />
          </label>

          <label>
            <span>Target</span>
            <select
              name="target"
              [ngModel]="target()"
              (ngModelChange)="setTarget($event)"
            >
              <option value="all">All products</option>
              <option value="category">Category</option>
              <option value="label">Label</option>
              <option value="selected">Selected products</option>
            </select>
          </label>

          @if (target() === 'category') {
            <label>
              <span>Category</span>
              <select
                name="category"
                [ngModel]="targetCategory()"
                (ngModelChange)="targetCategory.set($event)"
              >
                @for (category of categories(); track category) {
                  <option [value]="category">{{ category }}</option>
                }
              </select>
            </label>
          }

          @if (target() === 'label') {
            <label>
              <span>Label</span>
              <select
                name="label"
                [ngModel]="targetLabel()"
                (ngModelChange)="targetLabel.set($event)"
              >
                @for (label of labels(); track label) {
                  <option [value]="label">{{ label }}</option>
                }
              </select>
            </label>
          }

          <label>
            <span>Discount type</span>
            <select
              name="discount_type"
              [ngModel]="form().discount_type"
              (ngModelChange)="setDiscountType($event)"
            >
              <option>Percentage</option>
              <option>Amount</option>
              <option>QuantityBonus</option>
            </select>
          </label>

          @if (form().discount_type === 'QuantityBonus') {
            <label>
              <span>Buy quantity</span>
              <input
                type="number"
                min="1"
                name="buy_quantity"
                [ngModel]="form().buy_quantity"
                (ngModelChange)="updateForm('buy_quantity', $event)"
              />
            </label>

            <label>
              <span>Free quantity</span>
              <input
                type="number"
                min="1"
                name="free_quantity"
                [ngModel]="form().free_quantity"
                (ngModelChange)="updateForm('free_quantity', $event)"
              />
            </label>
          } @else {
            <label>
              <span>{{ form().discount_type === 'Percentage' ? 'Percent' : 'Amount cents' }}</span>
              <input
                type="number"
                min="1"
                name="discount_value"
                [ngModel]="discountValue()"
                (ngModelChange)="updateDiscountValue($event)"
              />
            </label>
          }

          <label>
            <span>Status</span>
            <select
              name="status"
              [ngModel]="form().status"
              (ngModelChange)="updateForm('status', $event)"
            >
              <option>Active</option>
              <option>Draft</option>
              <option>Archived</option>
            </select>
          </label>

          <label>
            <span>Starts</span>
            <input
              type="datetime-local"
              name="starts_at"
              [ngModel]="form().starts_at"
              (ngModelChange)="updateForm('starts_at', $event)"
            />
          </label>

          <label>
            <span>Ends</span>
            <input
              type="datetime-local"
              name="ends_at"
              [ngModel]="form().ends_at"
              (ngModelChange)="updateForm('ends_at', $event)"
            />
          </label>
        </div>

        @if (target() === 'selected') {
          <div class="selection-list">
            @for (product of products(); track product.id) {
              <label>
                <input
                  type="checkbox"
                  [checked]="selectedProductIds().includes(product.id)"
                  (change)="toggleProduct(product.id, $event)"
                />
                <span>{{ product.name }}</span>
              </label>
            }
          </div>
        }
      </section>

      <section class="promotion-list">
        <div class="section-header">
          <div>
            <h2>Existing promotions</h2>
            <p>{{ promotionSummary() }}</p>
          </div>
        </div>

        <div class="promotion-grid">
          @for (discount of discounts(); track discount.id) {
            <article class="promotion-card">
              <div>
                <strong>{{ discount.name }}</strong>
                <span>{{ promotionWindow(discount) }}</span>
              </div>
              <p>{{ discount.description || 'No description' }}</p>
              <dl>
                <div>
                  <dt>Value</dt>
                  <dd>{{ formatDiscount(discount) }}</dd>
                </div>
                <div>
                  <dt>Status</dt>
                  <dd>{{ discount.status }}</dd>
                </div>
                <div>
                  <dt>Type</dt>
                  <dd>{{ promotionTypeLabel(discount) }}</dd>
                </div>
                <div>
                  <dt>Products</dt>
                  <dd>{{ discount.scope === 'Global' ? 'All' : discount.product_ids.length }}</dd>
                </div>
              </dl>
              <button type="button" class="secondary-action" (click)="togglePromotion(discount)">
                {{ discount.status === 'Active' ? 'Archive' : 'Activate' }}
              </button>
            </article>
          }
        </div>
      </section>
    </section>
  `,
  styles: `
    .promotions-page {
      padding: 2rem;
      color: #172033;
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
    h2,
    p {
      margin: 0;
    }

    h1 {
      font-size: clamp(2rem, 4vw, 3rem);
      letter-spacing: -0.04em;
    }

    h2 {
      font-size: 1.15rem;
    }

    .description,
    .section-header p,
    .promotion-card p,
    .promotion-card span {
      color: #56657f;
      line-height: 1.5;
    }

    .description {
      max-width: 45rem;
      margin-top: 0.75rem;
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

    .builder,
    .promotion-list {
      display: grid;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 1rem;
      background: #ffffff;
      padding: 1.25rem;
      box-shadow: 0 10px 30px rgb(15 23 42 / 6%);
    }

    .builder {
      margin-bottom: 1.25rem;
    }

    .section-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 1rem;
    }

    .form-grid {
      display: grid;
      grid-template-columns: repeat(4, minmax(9rem, 1fr));
      gap: 0.85rem;
    }

    .wide {
      grid-column: span 2;
    }

    label {
      display: grid;
      gap: 0.4rem;
      font-weight: 700;
    }

    input,
    select {
      width: 100%;
      box-sizing: border-box;
      border: 1px solid #cbd5e1;
      border-radius: 0.75rem;
      padding: 0.75rem 0.85rem;
      color: #172033;
      font: inherit;
    }

    .selection-list {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr));
      gap: 0.65rem;
    }

    .global-toggle {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 1rem;
      border: 1px solid #dbe3ef;
      border-radius: 0.75rem;
      padding: 0.9rem;
    }

    .global-toggle div,
    .global-toggle label {
      display: grid;
      gap: 0.25rem;
    }

    .global-toggle label {
      grid-auto-flow: column;
      align-items: center;
    }

    .global-toggle input {
      width: auto;
    }

    .selection-list label {
      display: flex;
      align-items: center;
    }

    .selection-list input {
      width: auto;
    }

    .promotion-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(18rem, 1fr));
      gap: 0.85rem;
    }

    .promotion-card {
      display: grid;
      gap: 0.85rem;
      border: 1px solid #dbe3ef;
      border-radius: 0.75rem;
      padding: 1rem;
    }

    .promotion-card div:first-child {
      display: grid;
      gap: 0.25rem;
    }

    dl {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 0.65rem;
      margin: 0;
    }

    dt {
      color: #56657f;
      font-size: 0.75rem;
      font-weight: 700;
      text-transform: uppercase;
    }

    dd {
      margin: 0.2rem 0 0;
      font-weight: 800;
    }

    .primary-action,
    .secondary-action {
      width: fit-content;
      border-radius: 999px;
      padding: 0.75rem 1rem;
      font: inherit;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
    }

    .primary-action {
      border: 0;
      background: #2563eb;
      color: #ffffff;
    }

    .secondary-action {
      border: 1px solid #cbd5e1;
      background: #ffffff;
      color: #2563eb;
    }

    .primary-action:disabled {
      cursor: default;
      opacity: 0.65;
    }

    @media (max-width: 760px) {
      .promotions-page {
        padding: 1rem;
      }

      .section-header,
      .form-grid {
        display: grid;
        grid-template-columns: 1fr;
      }

      .wide {
        grid-column: auto;
      }
    }
  `,
})
export class PromotionsPage implements OnInit {
  private readonly productService = inject(ProductService);

  protected readonly isLoading = signal(true);
  protected readonly isSaving = signal(false);
  protected readonly message = signal('');
  protected readonly errorMessage = signal('');
  protected readonly products = signal<ProductSummary[]>([]);
  protected readonly discounts = signal<Discount[]>([]);
  protected readonly promotionSettings = signal<PromotionSettings>({
    promotions_enabled: true,
    updated_user: '',
    updated_at: '',
  });
  protected readonly target = signal<PromotionTarget>('label');
  protected readonly targetCategory = signal('');
  protected readonly targetLabel = signal('');
  protected readonly selectedProductIds = signal<string[]>([]);
  protected readonly form = signal<CreateDiscountRequest>({
    name: 'Seasonal promotion',
    description: '',
    discount_type: 'Percentage',
    scope: 'ProductSet',
    percentage_basis_points: 1000,
    amount_cents: null,
    currency: '',
    buy_quantity: 0,
    free_quantity: 0,
    min_product_count: 1,
    starts_at: '',
    ends_at: '',
    status: 'Active',
    product_ids: [],
    target_scope: 'Label',
  });

  protected readonly labels = computed(() => [
    ...new Set(this.products().flatMap((product) => product.labels ?? [])),
  ]);

  protected readonly categories = computed(() => [
    ...new Set(this.products().map((product) => product.category).filter(Boolean)),
  ]);

  protected readonly targetProducts = computed(() => {
    if (this.target() === 'all') {
      return this.products();
    }

    if (this.target() === 'category') {
      return this.products().filter((product) => product.category === this.targetCategory());
    }

    if (this.target() === 'label') {
      return this.products().filter((product) =>
        (product.labels ?? []).includes(this.targetLabel()),
      );
    }

    const selected = new Set(this.selectedProductIds());

    return this.products().filter((product) => selected.has(product.id));
  });

  public ngOnInit(): void {
    this.load();
  }

  protected updateForm<Key extends keyof CreateDiscountRequest>(
    key: Key,
    value: CreateDiscountRequest[Key],
  ): void {
    this.form.update((form) => ({
      ...form,
      [key]:
        key === 'buy_quantity' || key === 'free_quantity' || key === 'min_product_count'
          ? Number(value) || 0
          : value,
    }));
  }

  protected setTarget(target: PromotionTarget): void {
    this.target.set(target);
    this.targetCategory.set(this.targetCategory() || this.categories()[0] || '');
    this.targetLabel.set(this.targetLabel() || this.labels()[0] || '');
  }

  protected discountValue(): number {
    const form = this.form();

    return form.discount_type === 'Percentage'
      ? (form.percentage_basis_points ?? 0) / 100
      : form.amount_cents ?? 0;
  }

  protected updateDiscountValue(value: number | string): void {
    const numericValue = Number(value) || 0;

    this.form.update((form) => ({
      ...form,
      percentage_basis_points:
        form.discount_type === 'Percentage' ? Math.round(numericValue * 100) : null,
      amount_cents: form.discount_type === 'Amount' ? Math.round(numericValue) : null,
      currency: form.discount_type === 'Amount' ? this.defaultCurrency() : '',
    }));
  }

  protected setDiscountType(discountType: DiscountType): void {
    this.form.update((form) => ({
      ...form,
      discount_type: discountType,
      percentage_basis_points:
        discountType === 'Percentage' ? form.percentage_basis_points ?? 1000 : null,
      amount_cents: discountType === 'Amount' ? form.amount_cents ?? 500 : null,
      currency: discountType === 'Amount' ? this.defaultCurrency() : '',
      buy_quantity: discountType === 'QuantityBonus' ? form.buy_quantity || 1 : 0,
      free_quantity: discountType === 'QuantityBonus' ? form.free_quantity || 1 : 0,
      min_product_count:
        discountType === 'QuantityBonus'
          ? Math.max(1, (form.buy_quantity || 1) + (form.free_quantity || 1))
          : form.min_product_count,
    }));
  }

  protected toggleProduct(id: string, event: Event): void {
    const checked = (event.target as HTMLInputElement).checked;

    this.selectedProductIds.update((ids) =>
      checked ? [...new Set([...ids, id])] : ids.filter((candidate) => candidate !== id),
    );
  }

  protected createPromotion(): void {
    if (this.isSaving() || this.targetProducts().length === 0) {
      return;
    }

    this.isSaving.set(true);
    this.message.set('');
    this.errorMessage.set('');

    const targetScope = this.targetScopeLabel();
    const buyQuantity = Number(this.form().buy_quantity) || 0;
    const freeQuantity = Number(this.form().free_quantity) || 0;
    const request: CreateDiscountRequest = {
      ...this.form(),
      scope: this.target() === 'all' ? 'Global' : 'ProductSet',
      currency: this.form().discount_type === 'Amount' ? this.defaultCurrency() : '',
      buy_quantity: this.form().discount_type === 'QuantityBonus' ? buyQuantity : 0,
      free_quantity: this.form().discount_type === 'QuantityBonus' ? freeQuantity : 0,
      min_product_count:
        this.form().discount_type === 'QuantityBonus'
          ? Math.max(1, buyQuantity + freeQuantity)
          : this.form().min_product_count,
      product_ids:
        this.target() === 'all'
          ? []
          : this.targetProducts().map((product) => this.normalizedProductId(product.id)),
      target_scope: targetScope,
      target_label: this.target() === 'label' ? this.targetLabel() : '',
      target_category: this.target() === 'category' ? this.targetCategory() : '',
    };

    this.productService
      .createDiscount(request)
      .pipe(finalize(() => this.isSaving.set(false)))
      .subscribe({
        next: () => {
          this.message.set('Promotion created successfully.');
          this.load(false);
        },
        error: () => {
          this.errorMessage.set('Unable to create promotion.');
        },
      });
  }

  protected togglePromotion(discount: Discount): void {
    this.productService
      .updateDiscount(discount.id, {
        ...discount,
        status: discount.status === 'Active' ? 'Archived' : 'Active',
      })
      .subscribe({
        next: () => this.load(false),
        error: () => {
          this.errorMessage.set('Unable to update promotion.');
        },
      });
  }

  protected toggleGlobalPromotions(event: Event): void {
    const promotionsEnabled = (event.target as HTMLInputElement).checked;
    this.isSaving.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.productService
      .updatePromotionSettings({ promotions_enabled: promotionsEnabled })
      .pipe(finalize(() => this.isSaving.set(false)))
      .subscribe({
        next: (settings) => {
          this.promotionSettings.set(settings);
          this.message.set(
            settings.promotions_enabled
              ? 'Promotions enabled.'
              : 'Promotions paused globally.',
          );
        },
        error: () => {
          this.errorMessage.set('Unable to update promotion settings.');
        },
      });
  }

  protected formatDiscount(discount: Discount): string {
    return formatDiscountValue(discount, this.formatMoney.bind(this));
  }

  protected promotionWindow(discount: Discount): string {
    if (!discount.starts_at && !discount.ends_at) {
      return 'No schedule';
    }

    return `${discount.starts_at || 'Any time'} to ${discount.ends_at || 'No end'}`;
  }

  protected promotionTypeLabel(discount: Discount): string {
    const labels: Record<DiscountType, string> = {
      Percentage: 'Percentage',
      Amount: 'Amount off',
      QuantityBonus: 'Quantity bonus',
    };

    return labels[discount.discount_type];
  }

  protected promotionSummary(): string {
    const typeCounts = this.discounts().reduce<Record<DiscountType, number>>(
      (counts, discount) => ({
        ...counts,
        [discount.discount_type]: counts[discount.discount_type] + 1,
      }),
      {
        Percentage: 0,
        Amount: 0,
        QuantityBonus: 0,
      },
    );

    return `${this.discounts().length} configured: ${typeCounts.Percentage} percentage, ${typeCounts.Amount} amount, ${typeCounts.QuantityBonus} quantity bonus`;
  }

  private load(showLoading = true): void {
    if (showLoading) {
      this.isLoading.set(true);
    }

    forkJoin({
      products: this.productService.getProducts(),
      discounts: this.productService.getDiscounts(),
      settings: this.productService.getPromotionSettings(),
    })
      .pipe(finalize(() => this.isLoading.set(false)))
      .subscribe({
        next: ({ products, discounts, settings }) => {
          this.products.set(products);
          this.discounts.set(discounts);
          this.promotionSettings.set(settings);
          this.targetCategory.set(this.targetCategory() || this.categories()[0] || '');
          this.targetLabel.set(this.targetLabel() || this.labels()[0] || '');
        },
        error: () => {
          this.errorMessage.set('Unable to load promotions.');
        },
      });
  }

  private defaultCurrency(): string {
    return this.targetProducts()[0]?.currency ?? this.products()[0]?.currency ?? 'USD';
  }

  private formatMoney(priceCents: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(priceCents / 100);
  }

  private normalizedProductId(id: string): number | string {
    const numericId = Number(id);

    return Number.isFinite(numericId) ? numericId : id;
  }

  private targetScopeLabel(): CreateDiscountRequest['target_scope'] {
    const scopeByTarget: Record<PromotionTarget, CreateDiscountRequest['target_scope']> = {
      all: 'All',
      category: 'Category',
      label: 'Label',
      selected: 'Selected',
    };

    return scopeByTarget[this.target()];
  }
}
