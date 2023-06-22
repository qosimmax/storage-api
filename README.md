<div align="center">
    <h1>storage-api</h1>
</div>

## Description

## Пошли вы нахуй со своей задачей...

## Requirements

* [Go](https://golang.org) 
* [Postgres](https://www.postgresql.org/)
* [File Server API](https://github.com/qosimmax/file-server-api)

## Setup

In order to set up this app, run the following commands:

1. `git clone git@github.com:qosimmax/storage-api.git`
2. `cp .env.example .env`
3. Set correct env vars in .env
4. `make build`

## Run

In order to run this app, run the following commands:

1. `make run`

## DB tables

```sql
create table server
(
    id      varchar(36)  not null
        primary key,
    name    varchar(100) not null,
    address varchar(50)  not null
);

create unique index server_address_uindex
    on server (address)

```

```sql
create table file_info
(
    id         varchar(36) not null
        primary key,
    name       text,
    size       bigint,
    created_at timestamp
)
```

```sql
create table server_files
(
    file_id    varchar(36),
    server_id  varchar(36),
    part_size  bigint,
    "order"    integer,
    created_at timestamp
)
```
