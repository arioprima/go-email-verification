create table users(
    id varchar(36) primary key,
    user_name varchar(10) not null,
    email varchar(20) not null,
    password varchar(255) not null,
    is_verified boolean default false,
    verification_token varchar(255),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);