package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=workout_user password=workout_password dbname=workout_db sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	_, err = db.Exec("TRUNCATE workouts, workout_entries RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("failed to clean test database: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push-up day",
				Description:     "upper body training",
				DurationMinutes: 60,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench press",
						Sets:         10,
						Reps:         IntPtr(20),
						Weight:       FloatPtr(50),
						Notes:        "warn up correctly",
						OrderIndex:   1,
						Unit:         "kg",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid workout",
			workout: &Workout{
				Title:           "full body day",
				Description:     "full body training",
				DurationMinutes: 60,
				CaloriesBurned:  100,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         10,
						Reps:         IntPtr(20),
						Notes:        "warn up correctly",
						OrderIndex:   1,
						Unit:         "kg",
					},
					{
						ExerciseName:    "Abdominal crunch",
						Sets:            4,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(30),
						OrderIndex:      2,
						Weight:          FloatPtr(90),
						Notes:           "use a mat",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)

			retrieved, err := store.GetWorkoutById(createdWorkout.Id)
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.Title, retrieved.Title)
			assert.Equal(t, createdWorkout.Description, retrieved.Description)
			assert.Equal(t, len(createdWorkout.Entries), len(retrieved.Entries))

			for i, entry := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, entry.ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, entry.Sets)
				assert.Equal(t, tt.workout.Entries[i].Reps, entry.Reps)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}
