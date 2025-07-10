-- +goose Up
ALTER TABLE workout_entries
ALTER COLUMN weight DROP NOT NULL;

-- +goose Down
ALTER TABLE workout_entries
ALTER COLUMN weight SET NOT NULL;
