package store

import (
	"database/sql"
	"fmt"
)

type Workout struct {
	Id              int64          `json:"id"`
	Title           string         `json:"title"`
	UserId          int            `json:"user_id"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	Id              int64    `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
	Unit            string   `json:"unit"`
}

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkoutById(int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(int64) error
	GetWorkouts(take int, skip int) ([]Workout, error)
	GetWorkoutOwner(id int64) (int64, error)
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

func (p *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback() // Rollback if something goes wrong

	query :=
		`INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	err = tx.QueryRow(
		query,
		workout.UserId,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
	).Scan(&workout.Id)
	if err != nil {
		return nil, err
	}

	for index := range workout.Entries {
		err := createWorkoutEntry(tx, workout.Id, &workout.Entries[index])
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", workout)

	return workout, nil
}

func (p *PostgresWorkoutStore) GetWorkouts(take int32, skip int32) ([]Workout, error) {
	workouts := []Workout{}

	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := p.db.Query(query, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workout Workout
		err := rows.Scan(
			&workout.Id,
			&workout.Title,
			&workout.Description,
			&workout.DurationMinutes,
			&workout.CaloriesBurned,
		)
		if err != nil {
			return nil, err
		}

		workout.Entries = []WorkoutEntry{}
		workouts = append(workouts, workout)
	}

	return workouts, nil
}

func (p *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	workout := &Workout{}

	query :=
		`SELECT id, title, description, duration_minutes, calories_burned
	FROM workouts
	WHERE id = $1
	`

	err := p.db.QueryRow(query, id).Scan(&workout.Id, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if err != nil {
		return nil, err
	}

	entryQuery := `
		SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, unit, order_index
		FROM workout_entries
		WHERE workout_id = $1
		ORDER BY order_index
	`
	rows, err := p.db.Query(entryQuery, workout.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err := rows.Scan(
			&entry.Id,
			&entry.ExerciseName,
			&entry.Sets, &entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.Unit,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (p *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE workouts
		SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
		WHERE id = $5
	`

	result, err := tx.Exec(
		query,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
		workout.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Updating entries
	_, err = tx.Exec("DELETE FROM workout_entries WHERE workout_id = $1", workout.Id)
	if err != nil {
		return err
	}

	for _, entry := range workout.Entries {
		err := insertWorkoutEntry(tx, workout.Id, entry)
		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}

func (p *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `
	DELETE FROM workouts
	WHERE id = $1
	`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *PostgresWorkoutStore) GetWorkoutOwner(id int64) (int64, error) {
	var userId int64

	query := `
		SELECT user_id
		FROM workouts
		WHERE id = $1
	`

	err := p.db.QueryRow(query, id).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func createWorkoutEntry(tx *sql.Tx, workoutId int64, entry *WorkoutEntry) error {
	query :=
		`INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, unit, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`

	err := tx.QueryRow(query,
		workoutId,
		entry.ExerciseName,
		entry.Sets,
		entry.Reps,
		entry.DurationSeconds,
		entry.Weight,
		entry.Notes,
		entry.Unit,
		entry.OrderIndex,
	).Scan(&entry.Id)
	if err != nil {
		return err
	}

	return nil
}

func insertWorkoutEntry(tx *sql.Tx, workoutId int64, entry WorkoutEntry) error {
	query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, unit, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
	_, err := tx.Exec(
		query,
		workoutId,
		entry.ExerciseName,
		entry.Sets,
		entry.Reps,
		entry.DurationSeconds,
		entry.Weight,
		entry.Notes,
		entry.Unit,
		entry.OrderIndex,
	)

	if err != nil {
		return err
	}

	return nil
}
