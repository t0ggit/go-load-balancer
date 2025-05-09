-- +goose Up
-- +goose StatementBegin
create table bucket_settings (
    key_hash char(64) primary key,
    capacity integer not null,
    refill integer not null,
    refill_interval text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists bucket_settings;
-- +goose StatementEnd
