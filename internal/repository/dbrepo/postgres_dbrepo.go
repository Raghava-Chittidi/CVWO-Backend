package dbrepo

// import (
// 	"context"
// 	"database/sql"
// 	"time"

// 	"github.com/CVWO-Backend/internal/models"
// )

// type PostgresDBRepo struct {
// 	DB *sql.DB
// }

// func (m *PostgresDBRepo) Connection() *sql.DB {
// 	return m.DB
// }

// func (m *PostgresDBRepo) AllThreads() ([]*models.Thread, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
// 	defer cancel()

// 	query := `SELECT * FROM threads ORDER BY date`
// 	rows, err := m.DB.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var threads []*models.Thread
// 	for rows.Next() {
// 		var thread models.Thread
// 		err := rows.Scan(
// 			&thread.ID,
// 			&thread.Topic,
// 			&thread.Description,
// 			&thread.Creator,
// 			&thread.Comments,
// 			&thread.CreatedAt,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}
// 		threads = append(threads, &thread)
// 	}

// 	return threads, nil
// }