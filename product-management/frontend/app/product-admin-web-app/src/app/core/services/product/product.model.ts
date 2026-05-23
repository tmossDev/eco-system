export type ProductStatus = 'Draft' | 'Active' | 'Archived';
export type DiscountStatus = 'Draft' | 'Active' | 'Archived';
export type DiscountType = 'Percentage' | 'Amount';
export type DiscountScope = 'Global' | 'ProductSet';

export interface ProductSummary {
  id: string;
  sku: string;
  name: string;
  short_description: string;
  description: string;
  category: string;
  price_cents: number;
  currency: string;
  inventory_count: number;
  status: ProductStatus;
  photos: ProductPhoto[];
  discounts: Discount[];
}

export interface ProductDetails extends ProductSummary {}

export interface CreateProductRequest {
  sku: string;
  name: string;
  short_description: string;
  description: string;
  category: string;
  price_cents: number;
  currency: string;
  inventory_count: number;
  status: ProductStatus;
  photos: ProductPhoto[];
}

export interface UpdateProductRequest extends CreateProductRequest {}

export interface ProductPhoto {
  url: string;
  thumbnail_url: string;
  alt_text: string;
  is_primary: boolean;
}

export interface Discount {
  id: string;
  name: string;
  description: string;
  discount_type: DiscountType;
  scope: DiscountScope;
  percentage_basis_points: number | null;
  amount_cents: number | null;
  currency: string;
  min_product_count: number;
  starts_at: string;
  ends_at: string;
  status: DiscountStatus;
  product_ids: ProductID[];
}

export type ProductID = string | number;

export interface CreateDiscountRequest {
  name: string;
  description: string;
  discount_type: DiscountType;
  scope: DiscountScope;
  percentage_basis_points: number | null;
  amount_cents: number | null;
  currency: string;
  min_product_count: number;
  starts_at: string;
  ends_at: string;
  status: DiscountStatus;
  product_ids: ProductID[];
}

export interface UpdateDiscountRequest extends CreateDiscountRequest {}
