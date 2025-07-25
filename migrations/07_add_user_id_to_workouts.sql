-- +goose Up
ALTER TABLE workouts
ADD COLUMN user_id BIGINT NOT NULL,
ADD CONSTRAINT fk_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE workouts
DROP COLUMN user_id;
DROP CONSTRAINT fk_user;
