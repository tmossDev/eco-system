alter table products
  add column if not exists labels jsonb not null default '[]'::jsonb;

create index if not exists idx_products_labels on products using gin (labels);
