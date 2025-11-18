package repoutils

import (
	"context"
	"fmt"
	"online-courses/internal/database"

	"github.com/google/uuid"
)

func Exists(ctx context.Context, db database.DB, table, column string, id uuid.UUID) (bool, error) {

	query := fmt.Sprintf("select exists (select 1 from %s where %s = $1);", table, column)
	// query :="SELECT EXISTS(SELECT 1 FROM " + table + " WHERE " + column + " = $3);"

	var exists bool
	err := db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
