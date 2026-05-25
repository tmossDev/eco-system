alter table discounts
  add column if not exists buy_quantity bigint,
  add column if not exists free_quantity bigint;

alter table discounts
  drop constraint if exists ck_discounts_type,
  add constraint ck_discounts_type check (discount_type in ('Percentage', 'Amount', 'QuantityBonus'));

alter table discounts
  drop constraint if exists ck_discounts_percentage_value,
  drop constraint if exists ck_discounts_value,
  add constraint ck_discounts_value check (
    (discount_type = 'Percentage'
      and percentage_basis_points between 1 and 10000
      and amount_cents is null
      and currency is null
      and buy_quantity is null
      and free_quantity is null)
    or
    (discount_type = 'Amount'
      and amount_cents > 0
      and currency is not null
      and percentage_basis_points is null
      and buy_quantity is null
      and free_quantity is null)
    or
    (discount_type = 'QuantityBonus'
      and buy_quantity >= 1
      and free_quantity >= 1
      and min_product_count >= buy_quantity + free_quantity
      and percentage_basis_points is null
      and amount_cents is null
      and currency is null)
  );

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

GRANT SELECT, INSERT, UPDATE, DELETE ON promotion_settings TO app_user;
