package repository

const (
	InsertOrderQuery = `
INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
                    delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	InsertDeliveryQuery = `
INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	InsertPaymentQuery = `
INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt,
                     bank, delivery_cost, goods_total, custom_fee)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	InsertItemQuery = `
INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale,
                   size, total_price, nm_id, brand, status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	GetOrderWithJoinsQuery = `
SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
       o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
       d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
       p.transaction, p.request_id, p.currency, p.provider, p.amount,
       p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
FROM orders o
LEFT JOIN delivery d ON o.order_uid = d.order_uid
LEFT JOIN payment p ON o.order_uid = p.order_uid
WHERE o.order_uid = $1`

	GetItemsQuery = `
SELECT chrt_id, track_number, price, rid, name, sale, size,
       total_price, nm_id, brand, status
FROM items
WHERE order_uid = $1`

	GetAllOrdersQuery = `
SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
       delivery_service, shardkey, sm_id, date_created, oof_shard
FROM orders`
)
