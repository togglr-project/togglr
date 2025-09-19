create table product_info
(
    id         serial primary key,
    key        varchar(255) not null unique,
    value      text not null,
    created_at timestamp with time zone default now()
);

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM product_info WHERE key = 'client_id'
        ) THEN
            INSERT INTO product_info (key, value)
            VALUES ('client_id', gen_random_uuid()::TEXT);
        END IF;
    END $$;
