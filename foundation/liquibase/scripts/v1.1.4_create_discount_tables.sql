create table if not exists discounts (
  id bigserial primary key,
  name varchar(140) not null,
  description varchar(2000) not null default '',
  discount_type varchar(20) not null,
  scope varchar(20) not null,
  percentage_basis_points integer,
  amount_cents bigint,
  currency varchar(3),
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
  constraint ck_discounts_type check (discount_type in ('Percentage', 'Amount')),
  constraint ck_discounts_scope check (scope in ('Global', 'ProductSet')),
  constraint ck_discounts_status check (status in ('Draft', 'Active', 'Archived')),
  constraint ck_discounts_percentage_value check (
    (discount_type = 'Percentage' and percentage_basis_points between 1 and 10000 and amount_cents is null and currency is null)
    or
    (discount_type = 'Amount' and amount_cents > 0 and currency is not null and percentage_basis_points is null)
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
