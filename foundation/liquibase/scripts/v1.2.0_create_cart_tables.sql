create table if not exists carts (
  id bigserial primary key,
  user_id bigint not null references users,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  checked_out_at timestamp
);

create unique index if not exists uk_carts_active_user
  on carts(user_id)
  where checked_out_at is null;

create table if not exists cart_items (
  cart_id bigint not null references carts on delete cascade,
  product_id bigint not null references products,
  quantity bigint not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  primary key (cart_id, product_id),
  constraint ck_cart_items_quantity_positive check (quantity > 0)
);

create index if not exists idx_cart_items_product_id on cart_items(product_id);
