CREATE TABLE IF NOT EXISTS egw.product (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" varchar NOT NULL,
	short_description varchar NOT NULL,
	description varchar NOT NULL,
	price numeric NULL DEFAULT 0,
	created_at timestamp NULL DEFAULT now(),
	updated_at timestamp NULL DEFAULT now(),
	CONSTRAINT product_pk PRIMARY KEY (id)
);