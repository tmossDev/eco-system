package postgres

const (
	ProductProjection = `
SELECT
	id,
	sku,
	name,
	short_description,
	description,
	category,
	price_cents,
	currency,
	inventory_count,
	status,
	photos,
	labels,
	created_user,
	created_at::text,
	COALESCE(updated_user, 0),
	updated_at::text
FROM products
`
	ListProducts  = ProductProjection + "WHERE deleted_at IS NULL ORDER BY name"
	GetByID       = ProductProjection + "WHERE id = $1 AND deleted_at IS NULL"
	CreateProduct = `
INSERT INTO products (
	sku,
	name,
	short_description,
	description,
	category,
	price_cents,
	currency,
	inventory_count,
	status,
	photos,
	labels,
	created_user,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11::jsonb, $12, now(), now())
RETURNING id`
	UpdateProduct = `
UPDATE products
SET sku = $1,
	name = $2,
	short_description = $3,
	description = $4,
	category = $5,
	price_cents = $6,
	currency = $7,
	inventory_count = $8,
	status = $9,
	photos = $10::jsonb,
	labels = $11::jsonb,
	updated_user = $12,
	updated_at = now()
WHERE id = $13 AND deleted_at IS NULL`
	DeleteProduct = `
UPDATE products
SET deleted_user = $1,
	deleted_at = now(),
	updated_user = $1,
	updated_at = now(),
	status = 'Archived'
WHERE id = $2 AND deleted_at IS NULL`
)
