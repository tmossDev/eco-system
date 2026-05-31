package postgres

const (
	getOrCreateCart = `
INSERT INTO carts (user_id, created_at, updated_at)
VALUES ($1, now(), now())
ON CONFLICT (user_id) WHERE checked_out_at IS NULL
DO UPDATE SET updated_at = carts.updated_at
RETURNING id, user_id, created_at::text, updated_at::text`

	getCart = `
SELECT id, user_id, created_at::text, updated_at::text
FROM carts
WHERE user_id = $1 AND checked_out_at IS NULL`

	getCartItems = `
SELECT
	p.id,
	p.sku,
	p.name,
	ci.quantity,
	p.price_cents,
	p.currency,
	COALESCE(p.photos->0->>'thumbnail_url', '')
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
ORDER BY ci.created_at, p.id`

	addItem = `
INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
SELECT $1, p.id, $3, now(), now()
FROM products p
WHERE p.id = $2
	AND p.status = 'Active'
	AND p.deleted_at IS NULL
	AND p.inventory_count >= $3
ON CONFLICT (cart_id, product_id)
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = now()
WHERE cart_items.quantity + EXCLUDED.quantity <= (
	SELECT inventory_count FROM products WHERE id = EXCLUDED.product_id
)
RETURNING product_id`

	updateItem = `
UPDATE cart_items ci
SET quantity = $3, updated_at = now()
FROM products p
WHERE ci.cart_id = $1
	AND ci.product_id = $2
	AND p.id = ci.product_id
	AND p.status = 'Active'
	AND p.deleted_at IS NULL
	AND p.inventory_count >= $3
RETURNING ci.product_id`

	removeItem = `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	clearCart  = `DELETE FROM cart_items WHERE cart_id = $1`
	touchCart  = `UPDATE carts SET updated_at = now() WHERE id = $1`
)
