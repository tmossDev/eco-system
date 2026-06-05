create table if not exists orders (
  id bigserial primary key,
  user_id bigint not null references users,
  cart_id bigint not null references carts,
  status varchar(20) not null default 'Created',
  item_count bigint not null default 0,
  subtotal_cents bigint not null default 0,
  currency varchar(3) not null default 'USD',
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  constraint uk_orders_cart unique (cart_id),
  constraint ck_orders_status check (status in ('Created', 'Paid', 'Cancelled', 'Fulfilled')),
  constraint ck_orders_item_count_non_negative check (item_count >= 0),
  constraint ck_orders_subtotal_non_negative check (subtotal_cents >= 0)
);

create index if not exists idx_orders_user_id on orders(user_id);
create index if not exists idx_orders_status on orders(status);

create table if not exists order_items (
  order_id bigint not null references orders on delete cascade,
  product_id bigint not null references products,
  sku varchar(80) not null,
  name varchar(140) not null,
  quantity bigint not null,
  price_cents bigint not null,
  currency varchar(3) not null,
  line_total_cents bigint not null,
  thumbnail_url text not null default '',
  created_at timestamp not null default now(),
  primary key (order_id, product_id),
  constraint ck_order_items_quantity_positive check (quantity > 0),
  constraint ck_order_items_price_non_negative check (price_cents >= 0),
  constraint ck_order_items_line_total_non_negative check (line_total_cents >= 0)
);

create index if not exists idx_order_items_product_id on order_items(product_id);
