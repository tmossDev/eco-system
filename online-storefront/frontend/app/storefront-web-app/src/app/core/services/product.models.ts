export interface ProductPhoto {
  url: string;
  thumbnail_url: string;
  alt_text: string;
  is_primary: boolean;
}

export interface Product {
  id: string;
  sku: string;
  name: string;
  short_description: string;
  description: string;
  category: string;
  price_cents: number;
  currency: string;
  inventory_count: number;
  photos: ProductPhoto[];
  labels: string[];
}
