package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE workouts, workout_entries RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate test database: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "Valid Workout",
			workout: &Workout{
				Title:           "Morning Run",
				Description:     "Run for the morning",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Running",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(70.5),
						Notes:        "Warm up",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid Workout with invalid Entries",
			workout: &Workout{
				Title:           "Morning Run",
				Description:     "Run for the morning",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "push-ups",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(70.5),
						Notes:        "Warm up",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Running",
						Sets:            3,
						Reps:            IntPtr(10),
						DurationSeconds: IntPtr(30),
						Weight:          FloatPtr(70.5),
						Notes:           "Warm up",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.Create(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
