package sqlstorage

//nolint:depguard
import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/configs"
	"github.com/BingaBonga/otus_hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" // Postgres driver.
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, config *configs.DBConfig) (err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username, config.Password, config.Host, config.Port, config.Dbname)

	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	err = s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	isExist, err := s.exists(ctx, event.ID)
	if err != nil {
		return err
	}

	if isExist {
		return storage.ErrEventAlreadyExist
	}

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO event (id, title, start_date, duration, description,  owner,  remind_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ID,
		event.Title,
		event.StartDate,
		event.Duration,
		event.Description,
		event.Owner,
		event.RemindAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	isExist, err := s.exists(ctx, event.ID)
	if err != nil {
		return err
	}

	if !isExist {
		return storage.ErrEventDoesNotExist
	}

	_, err = s.db.ExecContext(
		ctx,
		`UPDATE event
			SET title=$1,
    		    start_date=$2, 
    		    duration=$3, 
    		    description=$4, 
    		    owner=$5, 
    		    remind_at=$6
			WHERE id=$7`,
		event.Title,
		event.StartDate,
		event.Duration,
		event.Description,
		event.Owner,
		event.RemindAt,
		event.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	isExist, err := s.exists(ctx, id)
	if err != nil {
		return err
	}

	if !isExist {
		return storage.ErrEventDoesNotExist
	}

	_, err = s.db.ExecContext(ctx, "DELETE FROM event WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

//nolint:lll
func (s *Storage) GetEventsByPeriod(ctx context.Context, owner string, startTime time.Time, endTime time.Time) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, 
       			title, 
       			start_date, 
    		    duration, 
    		    description, 
    		    owner, 
    		    remind_at
			FROM event
			WHERE owner = $1 AND start_date >=$2 AND start_date < $3`,
		owner, startTime, endTime,
	)
	if err != nil {
		return nil, err
	}

	events := make([]storage.Event, 0)
	for rows.Next() {
		var ev storage.Event
		if err = rows.Scan(
			&ev.ID,
			&ev.Title,
			&ev.StartDate,
			&ev.Duration,
			&ev.Description,
			&ev.Owner,
			&ev.RemindAt,
		); err != nil {
			return nil, err
		}
		events = append(events, ev)
	}

	return events, nil
}

func (s *Storage) exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	row, err := s.db.QueryContext(ctx, "SELECT EXISTS(SELECT * FROM event WHERE id = $1)", id)
	if err != nil {
		return false, err
	}

	for row.Next() {
		if err := row.Scan(&exists); err != nil {
			return false, err
		}
	}

	return exists, nil
}
