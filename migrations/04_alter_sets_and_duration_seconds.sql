-- +goose Up
ALTER TABLE workout_entries
ALTER COLUMN reps DROP NOT NULL,
ALTER COLUMN duration_seconds DROP NOT NULL;

-- +goose Down
ALTER TABLE workout_entries
ALTER COLUMN reps SET NOT NULL,
ALTER COLUMN duration_seconds SET NOT NULL;
