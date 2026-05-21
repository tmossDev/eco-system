package postgres

const (
	ProductProjection = `
SELECT
	id,
	sku,
	name,
	description,
	category,
	price_cents,
	currency,
	inventory_count,
	status,
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
	description,
	category,
	price_cents,
	currency,
	inventory_count,
	status,
	created_user,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now())
RETURNING id`
	UpdateProduct = `
UPDATE products
SET sku = $1,
	name = $2,
	description = $3,
	category = $4,
	price_cents = $5,
	currency = $6,
	inventory_count = $7,
	status = $8,
	updated_user = $9,
	updated_at = now()
WHERE id = $10 AND deleted_at IS NULL`
	DeleteProduct = `
UPDATE products
SET deleted_user = $1,
	deleted_at = now(),
	updated_user = $1,
	updated_at = now(),
	status = 'Archived'
WHERE id = $2 AND deleted_at IS NULL`
)
