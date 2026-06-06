export interface Order {
  id: number;
  user_id: number;
  cart_id: number;
  status: OrderStatus;
  items: OrderItem[];
  item_count: number;
  subtotal_cents: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export type OrderStatus =
  | 'Order Submitted'
  | 'Order Confirmed'
  | 'Order Fulfillment'
  | 'Order Out For Delivery'
  | 'Order Delivered'
  | 'Order Complete'
  | 'Order Returned'
  | 'Order Cancelled';

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

export interface UpdateOrderStatusRequest {
  status: OrderStatus;
}
