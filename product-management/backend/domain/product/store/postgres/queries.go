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
	COALESCE((
		SELECT jsonb_agg(
			jsonb_build_object(
				'id', d.id,
				'name', d.name,
				'description', d.description,
				'discount_type', d.discount_type,
				'scope', d.scope,
				'percentage_basis_points', d.percentage_basis_points,
				'amount_cents', d.amount_cents,
				'currency', COALESCE(d.currency, ''),
				'min_product_count', d.min_product_count,
				'starts_at', COALESCE(d.starts_at::text, ''),
				'ends_at', COALESCE(d.ends_at::text, ''),
				'status', d.status,
				'product_ids', COALESCE(dp.product_ids, '[]'::jsonb),
				'created_user', d.created_user,
				'created_at', d.created_at::text,
				'updated_user', COALESCE(d.updated_user, 0),
				'updated_at', d.updated_at::text
			)
			ORDER BY d.name
		)
		FROM discounts d
		LEFT JOIN LATERAL (
			SELECT jsonb_agg(discount_products.product_id ORDER BY discount_products.product_id) AS product_ids
			FROM discount_products
			WHERE discount_products.discount_id = d.id
		) dp ON true
		WHERE d.deleted_at IS NULL
		  AND (
			d.scope = 'Global'
			OR EXISTS (
				SELECT 1
				FROM discount_products product_match
				WHERE product_match.discount_id = d.id
				  AND product_match.product_id = products.id
			)
		  )
	), '[]'::jsonb) AS discounts,
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
	created_user,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::jsonb, $11, now(), now())
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
	updated_user = $11,
	updated_at = now()
WHERE id = $12 AND deleted_at IS NULL`
	DeleteProduct = `
UPDATE products
SET deleted_user = $1,
	deleted_at = now(),
	updated_user = $1,
	updated_at = now(),
	status = 'Archived'
WHERE id = $2 AND deleted_at IS NULL`
	DiscountProjection = `
SELECT
	d.id,
	d.name,
	d.description,
	d.discount_type,
	d.scope,
	d.percentage_basis_points,
	d.amount_cents,
	COALESCE(d.currency, ''),
	d.min_product_count,
	COALESCE(d.starts_at::text, ''),
	COALESCE(d.ends_at::text, ''),
	d.status,
	COALESCE((
		SELECT jsonb_agg(discount_products.product_id ORDER BY discount_products.product_id)
		FROM discount_products
		WHERE discount_products.discount_id = d.id
	), '[]'::jsonb),
	d.created_user,
	d.created_at::text,
	COALESCE(d.updated_user, 0),
	d.updated_at::text
FROM discounts d
`
	ListDiscounts   = DiscountProjection + "WHERE d.deleted_at IS NULL ORDER BY d.name"
	GetDiscountByID = DiscountProjection + "WHERE d.id = $1 AND d.deleted_at IS NULL"
	CreateDiscount  = `
INSERT INTO discounts (
	name,
	description,
	discount_type,
	scope,
	percentage_basis_points,
	amount_cents,
	currency,
	min_product_count,
	starts_at,
	ends_at,
	status,
	created_user,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8, NULLIF($9, '')::timestamp, NULLIF($10, '')::timestamp, $11, $12, now(), now())
RETURNING id`
	UpdateDiscount = `
UPDATE discounts
SET name = $1,
	description = $2,
	discount_type = $3,
	scope = $4,
	percentage_basis_points = $5,
	amount_cents = $6,
	currency = NULLIF($7, ''),
	min_product_count = $8,
	starts_at = NULLIF($9, '')::timestamp,
	ends_at = NULLIF($10, '')::timestamp,
	status = $11,
	updated_user = $12,
	updated_at = now()
WHERE id = $13 AND deleted_at IS NULL`
	DeleteDiscount = `
UPDATE discounts
SET deleted_user = $1,
	deleted_at = now(),
	updated_user = $1,
	updated_at = now(),
	status = 'Archived'
WHERE id = $2 AND deleted_at IS NULL`
	ReplaceDiscountProducts = "DELETE FROM discount_products WHERE discount_id = $1"
	InsertDiscountProduct   = "INSERT INTO discount_products (discount_id, product_id) VALUES ($1, $2)"
)
