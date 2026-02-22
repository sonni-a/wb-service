CREATE TABLE payment (
    order_uid UUID PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR NOT NULL,
    request_id VARCHAR NOT NULL,
    currency VARCHAR NOT NULL,
    provider VARCHAR NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL );