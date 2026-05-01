-- +goose Up

ALTER TABLE users
ADD COLUMN is_deleted BOOLEAN DEFAULT false;




-- +goose Down
ALTER TABLE users drop COLUMN is_deleted;
