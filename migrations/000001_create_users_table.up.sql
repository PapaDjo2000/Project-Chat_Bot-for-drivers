create extension if not exists  "uuid-ossp";

create table if not exists pr.users(
id uuid default uuid_generate_v4() primary key,
name varchar(255)not null, 
chat_id int not null,
active bool not null
);

create table if not exists pr.reports(
id uuid default uuid_generate_v4()primary key,
user_id uuid references pr.users (id) on delete cascade,
date timestamp not null default now(),
request jsonb,
response jsonb 
);

insert into pr.users (name) value ('maks')