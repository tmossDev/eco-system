package postgres

const (
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
	COALESCE(d.buy_quantity, 0),
	COALESCE(d.free_quantity, 0),
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
	buy_quantity,
	free_quantity,
	min_product_count,
	starts_at,
	ends_at,
	status,
	created_user,
	created_at,
	updated_at
) VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), NULLIF($8, 0), NULLIF($9, 0), $10, NULLIF($11, '')::timestamp, NULLIF($12, '')::timestamp, $13, $14, now(), now())
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
	buy_quantity = NULLIF($8, 0),
	free_quantity = NULLIF($9, 0),
	min_product_count = $10,
	starts_at = NULLIF($11, '')::timestamp,
	ends_at = NULLIF($12, '')::timestamp,
	status = $13,
	updated_user = $14,
	updated_at = now()
WHERE id = $15 AND deleted_at IS NULL`
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
	GetPromotionSettings    = `
SELECT
	promotions_enabled,
	COALESCE(updated_user, 0),
	updated_at::text
FROM promotion_settings
WHERE id = 1`
	UpdatePromotionSettings = `
UPDATE promotion_settings
SET promotions_enabled = $1,
	updated_user = $2,
	updated_at = now()
WHERE id = 1`
)
