alter table products
  add column if not exists short_description varchar(280) not null default '',
  add column if not exists photos jsonb not null default '[]'::jsonb;

update products
set short_description = case sku
  when 'GEN-MUG-001' then 'Durable 350ml ceramic mug for daily coffee and tea.'
  when 'APP-TEE-002' then 'Soft unisex tee in core store colors.'
  when 'DIG-GUIDE-003' then 'Downloadable buying guide for new store customers.'
  when 'KIT-STARTER-004' then 'Bundle-ready starter kit for curated product sets.'
  else left(description, 280)
end
where short_description = '';

update products
set photos = case sku
  when 'GEN-MUG-001' then '[{"url":"https://images.unsplash.com/photo-1514228742587-6b1558fcca3d?auto=format&fit=crop&w=1200&q=80","thumbnail_url":"https://images.unsplash.com/photo-1514228742587-6b1558fcca3d?auto=format&fit=crop&w=160&q=70","alt_text":"White ceramic mug on a table","is_primary":true}]'::jsonb
  when 'APP-TEE-002' then '[{"url":"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=1200&q=80","thumbnail_url":"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?auto=format&fit=crop&w=160&q=70","alt_text":"Folded cotton tee","is_primary":true}]'::jsonb
  when 'KIT-STARTER-004' then '[{"url":"https://images.unsplash.com/photo-1549465220-1a8b9238cd48?auto=format&fit=crop&w=1200&q=80","thumbnail_url":"https://images.unsplash.com/photo-1549465220-1a8b9238cd48?auto=format&fit=crop&w=160&q=70","alt_text":"Wrapped gift kit box","is_primary":true}]'::jsonb
  else photos
end
where photos = '[]'::jsonb;
