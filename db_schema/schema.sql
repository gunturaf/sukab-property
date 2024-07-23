-- minimum version: postgres 14

create table properties
(
    property_id     bigserial
        constraint properties_pk
            primary key,
    prefecture      varchar(500) collate "ja-JP-x-icu",
    city            varchar(500) collate "ja-JP-x-icu",
    town            varchar(500) collate "ja-JP-x-icu",
    chome           integer,
    banchi          integer,
    go              integer,
    building        varchar(500) collate "ja-JP-x-icu",
    price           bigint,
    nearest_station varchar(500) collate "ja-JP-x-icu",
    property_type   varchar(500) collate "ja-JP-x-icu",
    land_area       varchar(500) collate "ja-JP-x-icu",
    created_at      timestamp with time zone default NOW()
);
