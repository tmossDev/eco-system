export type ProductStatus = 'Draft' | 'Active' | 'Archived';

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
