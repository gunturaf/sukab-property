package property

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type PropertyRepository interface {
	Insert(ctx context.Context, property PropertyData) error
	ListAll(ctx context.Context) ([]PropertyData, error)
}

func NewRepo(db *sqlx.DB) *PropertyRepo {
	return &PropertyRepo{
		db: db,
	}
}

type PropertyRepo struct {
	db *sqlx.DB
}

func (pr *PropertyRepo) Insert(ctx context.Context, property PropertyData) error {
	query := `INSERT INTO properties (
								property_id, prefecture, city, town, chome, banchi, go, building,
                price, nearest_station, property_type, land_area
						) VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	values := []any{
		property.Prefecture,
		property.City,
		property.Town,
		property.Chome,
		property.Banchi,
		property.Go,
		property.Building,
		property.Price,
		property.NearestStation,
		property.PropertyType,
		property.LandArea,
	}

	tx := pr.db.MustBeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()

	_, errExec := tx.ExecContext(ctx, query, values...)
	if errExec != nil {
		return errExec
	}

	return tx.Commit()
}

func (pr *PropertyRepo) ListAll(ctx context.Context) ([]PropertyData, error) {
	var results []PropertyData

	// at a later date, we might need to set some filters and/or data sorting in this query
	query := `SELECT 
								property_id, prefecture, city, town, chome, banchi, go, building,
                price, nearest_station, property_type, land_area
	          FROM properties`
	err := pr.db.SelectContext(ctx, &results, query)
	if err != nil {
		return nil, err
	}

	return results, nil
}
