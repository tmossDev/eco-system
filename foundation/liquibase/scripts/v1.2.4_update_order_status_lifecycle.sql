ALTER TABLE orders DROP CONSTRAINT IF EXISTS ck_orders_status;

ALTER TABLE orders ALTER COLUMN status TYPE varchar(40);
ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'Order Submitted';

UPDATE orders
SET status = CASE status
  WHEN 'Created' THEN 'Order Submitted'
  WHEN 'Paid' THEN 'Order Confirmed'
  WHEN 'Fulfilled' THEN 'Order Complete'
  WHEN 'Cancelled' THEN 'Order Cancelled'
  ELSE status
END;

ALTER TABLE orders
  ADD CONSTRAINT ck_orders_status
  CHECK (status in (
    'Order Submitted',
    'Order Confirmed',
    'Order Fulfillment',
    'Order Out For Delivery',
    'Order Delivered',
    'Order Complete',
    'Order Returned',
    'Order Cancelled'
  ));

