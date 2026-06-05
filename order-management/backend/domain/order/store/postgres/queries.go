package postgres

const (
	listOrders = `
SELECT id, user_id, cart_id, status, item_count, subtotal_cents, currency, created_at::text, updated_at::text
FROM orders
ORDER BY created_at DESC, id DESC`

	getOrder = `
SELECT id, user_id, cart_id, status, item_count, subtotal_cents, currency, created_at::text, updated_at::text
FROM orders
WHERE id = $1`

	getOrderItems = `
SELECT product_id, sku, name, quantity, price_cents, currency, line_total_cents, thumbnail_url
FROM order_items
WHERE order_id = $1
ORDER BY created_at, product_id`

	updateOrderStatus = `
UPDATE orders
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, user_id, cart_id, status, item_count, subtotal_cents, currency, created_at::text, updated_at::text`
)
