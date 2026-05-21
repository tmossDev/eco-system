export type ProductStatus = 'Draft' | 'Active' | 'Archived';

export interface ProductSummary {
  id: string;
  sku: string;
  name: string;
  description: string;
  category: string;
  price_cents: number;
  currency: string;
  inventory_count: number;
  status: ProductStatus;
}

export interface ProductDetails extends ProductSummary {}

export interface CreateProductRequest {
  sku: string;
  name: string;
  description: string;
  category: string;
  price_cents: number;
  currency: string;
  inventory_count: number;
  status: ProductStatus;
}

export interface UpdateProductRequest extends CreateProductRequest {}
