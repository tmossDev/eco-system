package postgres

const (
	getCartForCheckout = `
SELECT id
FROM carts
WHERE user_id = $1 AND checked_out_at IS NULL
FOR UPDATE`

	countCartItems = `SELECT COUNT(*) FROM cart_items WHERE cart_id = $1`

	decrementProductInventory = `
WITH updated AS (
	UPDATE products p
	SET inventory_count = p.inventory_count - ci.quantity,
		updated_at = now()
	FROM cart_items ci
	WHERE ci.cart_id = $1
		AND p.id = ci.product_id
		AND p.status = 'Active'
		AND p.deleted_at IS NULL
		AND p.inventory_count >= ci.quantity
	RETURNING p.id
)
SELECT COUNT(*) FROM updated`

	createOrder = `
INSERT INTO orders (user_id, cart_id, status, item_count, subtotal_cents, currency, created_at, updated_at)
SELECT $2, $1, 'Created', SUM(ci.quantity), SUM(p.price_cents * ci.quantity), COALESCE(MIN(p.currency), 'USD'), now(), now()
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1
GROUP BY ci.cart_id
RETURNING id, user_id, cart_id, status, item_count, subtotal_cents, currency, created_at::text, updated_at::text`

	createOrderItems = `
INSERT INTO order_items (order_id, product_id, sku, name, quantity, price_cents, currency, line_total_cents, thumbnail_url, created_at)
SELECT
	$2,
	p.id,
	p.sku,
	p.name,
	ci.quantity,
	p.price_cents,
	p.currency,
	p.price_cents * ci.quantity,
	COALESCE(p.photos->0->>'thumbnail_url', ''),
	now()
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.cart_id = $1`

	checkoutCart = `UPDATE carts SET checked_out_at = now(), updated_at = now() WHERE id = $1`

	getOrder = `
SELECT id, user_id, cart_id, status, item_count, subtotal_cents, currency, created_at::text, updated_at::text
FROM orders
WHERE id = $1`

	listOrders = `
SELECT id, user_id, cart_id, status, item_count, subtotal_cents, currency, created_at::text, updated_at::text
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC, id DESC`

	getOrderItems = `
SELECT product_id, sku, name, quantity, price_cents, currency, line_total_cents, thumbnail_url
FROM order_items
WHERE order_id = $1
ORDER BY created_at, product_id`
)
