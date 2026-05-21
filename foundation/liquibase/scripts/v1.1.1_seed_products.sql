insert into products (
  sku,
  name,
  description,
  category,
  price_cents,
  currency,
  inventory_count,
  status,
  created_user
)
values
  (
    'GEN-MUG-001',
    'Everyday Ceramic Mug',
    'A durable 350ml mug for daily coffee, tea, or desk rituals.',
    'Home',
    1299,
    'USD',
    48,
    'Active',
    1
  ),
  (
    'APP-TEE-002',
    'Organic Cotton Tee',
    'Soft unisex cotton tee available in core store colors.',
    'Apparel',
    2499,
    'USD',
    82,
    'Active',
    1
  ),
  (
    'DIG-GUIDE-003',
    'Digital Buying Guide',
    'Downloadable product guide for new store customers.',
    'Digital',
    499,
    'USD',
    999,
    'Draft',
    1
  ),
  (
    'KIT-STARTER-004',
    'Starter Gift Kit',
    'A bundled kit that can represent arbitrary grouped products.',
    'Bundles',
    5499,
    'USD',
    16,
    'Archived',
    1
  )
on conflict (sku) do nothing;
