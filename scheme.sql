create table history
(
    id INTEGER
        primary key,
    symbol CHARACTER,
    value DECIMAL(10,2),
    created_d DATE default current_date,
    created_dt datetime default current_timestamp
);

create table portfolio
(
    id INTEGER
        primary key autoincrement,
    symbol CHARACTER not null,
    open_price DECIMAL(16,8) not null,
    volume DECIMAL(16,8) not null,
    open_dt DATETIME default current_timestamp,
    close_dt DATETIME,
    is_close BOOLEAN default 0,
    close_price DECIMAL(16,8),
    update_dt DATETIME,
    currency SMALLINT default 1
);

create table symbols
(
    id INTEGER
        primary key autoincrement,
    symbol CHARACTER,
    currency SMALLINT
);

