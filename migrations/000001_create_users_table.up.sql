create extension if not exists  "uuid-ossp";

create table if not exists pr.users(
id uuid default uuid_generate_v4() primary key,
name varchar(255)not null, 
chat_id BIGINT not null UNIQUE
);

create table if not exists pr.reports(
id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
user_id BIGINT NOT NULL REFERENCES pr.users (chat_id) ON DELETE cascade,
date timestamp not null default now(),
request jsonb,
response jsonb 
);
