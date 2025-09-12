create table product_info
(
    id         serial primary key,
    key        text not null unique,
    value      text not null,
    created_at timestamp with time zone default now()
);
