create table migrations (
    id integer constraint pk_migrations primary key,
    applied_at date not null default now()
);

create table tags (
    id uuid default gen_random_uuid() constraint pk_tags primary key,
    name varchar(50) not null,
    name_normalised varchar(50) generated always as ( upper(name) ) stored,
    user_id uuid not null,

    constraint ix_tags_user_id_name_normalised unique (name_normalised, user_id)
);

create table links (
    id uuid constraint pk_links primary key,
    user_id uuid not null
);

create table link_tags (
    id uuid default gen_random_uuid() constraint pk_link_tags primary key,
    user_id uuid not null,
    link_id uuid not null references links(id) on delete cascade,
    tag_id uuid not null references tags(id) on delete cascade,

    constraint ix_link_tags_link_tag unique (link_id, tag_id)
);