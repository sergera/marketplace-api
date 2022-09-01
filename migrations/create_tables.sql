CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
	price INT NOT NULL,
	status CITEXT NOT NULL,
	date_created TIMESTAMP NOT NULL,
	CONSTRAINT order_status_length CHECK (LENGTH(address) <= 11),
);

CREATE INDEX order_status ON orders(status);
CREATE INDEX order_price ON orders(price);
