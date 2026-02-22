CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    order_uid UUID REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR NOT NULL,
    status INTEGER NOT NULL );