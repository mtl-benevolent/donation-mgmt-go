-- +goose Up
-- +goose StatementBegin
alter default privileges grant select, insert, update, delete on tables to "donations_rw";
alter default privileges grant select, update on sequences to "donations_rw";

alter default privileges grant drop on tables to "donations_maintenance";
alter default privileges grant drop on sequences to "donations_maintenance";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'No down implementation';
-- +goose StatementEnd
