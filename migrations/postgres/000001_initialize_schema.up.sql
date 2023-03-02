CREATE TABLE users
(
    id            uuid UNIQUE         NOT NULL DEFAULT gen_random_uuid(),
    login         varchar(255) UNIQUE NOT NULL,
    password_hash bytea               NOT NULL,
    password_salt varchar(255)        NOT NULL,
    is_admin      boolean             NOT NULL DEFAULT FALSE,
    last_online   timestamp           NOT NULL DEFAULT to_timestamp(0)
);

CREATE TABLE companies
(
    id          uuid UNIQUE  NOT NULL DEFAULT gen_random_uuid(),
    owner_id    uuid REFERENCES users (id),
    name        varchar(255) NOT NULL,
    time_offset interval     NOT NULL DEFAULT '00:00:00'
);

CREATE TABLE devices
(
    id          uuid UNIQUE                    NOT NULL DEFAULT gen_random_uuid(),
    token       varchar(255)                   NOT NULL DEFAULT md5(random()::text),
    company_id  uuid REFERENCES companies (id) NOT NULL,
    name        varchar(255)                   NOT NULL,
    resource    integer                        NOT NULL,
    last_online timestamp                      NOT NULL DEFAULT to_timestamp(0)
);

CREATE TABLE schedules
(
    id         uuid UNIQUE                  NOT NULL DEFAULT gen_random_uuid(),
    device_id  uuid REFERENCES devices (id) NOT NULL,
    week_day   integer                      NOT NULL,
    time_start integer                      NOT NULL,
    time_stop  integer                      NOT NULL
);

CREATE TABLE journals
(
    id              uuid UNIQUE                  NOT NULL DEFAULT gen_random_uuid(),
    device_id       uuid REFERENCES devices (id) NOT NULL,
    timestamp_start timestamp                    NOT NULL,
    timestamp_end   timestamp                    NOT NULL,
    done            bool                         NOT NULL DEFAULT FALSE
);