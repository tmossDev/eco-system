create table if not exists products (
  id bigserial primary key,
  sku varchar(80) not null constraint uk_products_sku unique,
  name varchar(140) not null,
  description varchar(2000) not null default '',
  category varchar(120) not null,
  price_cents bigint not null default 0,
  currency varchar(3) not null default 'USD',
  inventory_count bigint not null default 0,
  status varchar(20) not null default 'Draft',
  created_user bigint not null references users,
  created_at timestamp not null default now(),
  updated_user bigint references users,
  updated_at timestamp not null default now(),
  deleted_user bigint references users,
  deleted_at timestamp,
  constraint ck_products_price_non_negative check (price_cents >= 0),
  constraint ck_products_inventory_non_negative check (inventory_count >= 0),
  constraint ck_products_status check (status in ('Draft', 'Active', 'Archived'))
);

create index if not exists idx_products_status on products(status);
create index if not exists idx_products_category on products(category);
