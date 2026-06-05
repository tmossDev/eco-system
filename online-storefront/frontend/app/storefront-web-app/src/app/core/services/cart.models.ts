export interface Cart {
  id: number;
  user_id: number;
  items: CartItem[];
  item_count: number;
  subtotal_cents: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface CartItem {
  product_id: number;
  sku: string;
  name: string;
  quantity: number;
  price_cents: number;
  currency: string;
  line_total_cents: number;
  thumbnail_url?: string;
}

export interface Order {
  id: number;
  user_id: number;
  cart_id: number;
  status: string;
  items: OrderItem[];
  item_count: number;
  subtotal_cents: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  product_id: number;
  sku: string;
  name: string;
  quantity: number;
  price_cents: number;
  currency: string;
  line_total_cents: number;
  thumbnail_url?: string;
}
