CREATE TABLE egw.order_item (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	order_id uuid NOT NULL,
	product_id uuid NOT NULL,
	product_name varchar NOT NULL,
	quantity integer NOT NULL,
	CONSTRAINT order_item_pk PRIMARY KEY (id),
	CONSTRAINT order_item_fk FOREIGN KEY (product_id) REFERENCES egw.product(id),
	CONSTRAINT order_item_fk_1 FOREIGN KEY (order_id) REFERENCES egw."order"(id)
);