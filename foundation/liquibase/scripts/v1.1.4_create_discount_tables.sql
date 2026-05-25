create table if not exists discounts (
  id bigserial primary key,
  name varchar(140) not null,
  description varchar(2000) not null default '',
  discount_type varchar(20) not null,
  scope varchar(20) not null,
  percentage_basis_points integer,
  amount_cents bigint,
  currency varchar(3),
  buy_quantity bigint,
  free_quantity bigint,
  min_product_count bigint not null default 1,
  starts_at timestamp,
  ends_at timestamp,
  status varchar(20) not null default 'Draft',
  created_user bigint not null references users,
  created_at timestamp not null default now(),
  updated_user bigint references users,
  updated_at timestamp not null default now(),
  deleted_user bigint references users,
  deleted_at timestamp,
  constraint ck_discounts_type check (discount_type in ('Percentage', 'Amount', 'QuantityBonus')),
  constraint ck_discounts_scope check (scope in ('Global', 'ProductSet')),
  constraint ck_discounts_status check (status in ('Draft', 'Active', 'Archived')),
  constraint ck_discounts_value check (
    (discount_type = 'Percentage' and percentage_basis_points between 1 and 10000 and amount_cents is null and currency is null and buy_quantity is null and free_quantity is null)
    or
    (discount_type = 'Amount' and amount_cents > 0 and currency is not null and percentage_basis_points is null and buy_quantity is null and free_quantity is null)
    or
    (discount_type = 'QuantityBonus' and buy_quantity >= 1 and free_quantity >= 1 and min_product_count >= buy_quantity + free_quantity and percentage_basis_points is null and amount_cents is null and currency is null)
  ),
  constraint ck_discounts_min_product_count check (min_product_count >= 1),
  constraint ck_discounts_dates check (ends_at is null or starts_at is null or ends_at > starts_at)
);

create table if not exists discount_products (
  discount_id bigint not null references discounts(id) on delete cascade,
  product_id bigint not null references products(id) on delete cascade,
  created_at timestamp not null default now(),
  primary key (discount_id, product_id)
);

create index if not exists idx_discounts_scope on discounts(scope);
create index if not exists idx_discounts_status on discounts(status);
create index if not exists idx_discount_products_product on discount_products(product_id);

create table if not exists promotion_settings (
  id smallint primary key default 1,
  promotions_enabled boolean not null default true,
  updated_user bigint references users,
  updated_at timestamp not null default now(),
  constraint ck_promotion_settings_singleton check (id = 1)
);

insert into promotion_settings (id, promotions_enabled, updated_at)
values (1, true, now())
on conflict (id) do nothing;
