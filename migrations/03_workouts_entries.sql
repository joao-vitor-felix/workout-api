-- +goose Up
CREATE TABLE IF NOT EXISTS workout_entries (
  id BIGSERIAL PRIMARY KEY,
  workout_id BIGSERIAL NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
  exercise_name VARCHAR(255) NOT NULL,
  sets INT NOT NULL,
  reps INT NOT NULL,
  duration_seconds INT NOT NULL,
  weight DECIMAL(5, 2) NOT NULL,
  notes TEXT,
  order_index INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT valid_workout_entry CHECK (
    (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND
    (reps IS NULL OR duration_seconds IS NULL)
  )
);

-- +goose Down
DROP TABLE IF EXISTS workout_entries;
