CREATE TABLE delivery (
    order_uid UUID PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    zip VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    region VARCHAR NOT NULL,
    email VARCHAR NOT NULL );