package postgres

import (
	"context"
	"fmt"
	"strings"

	gohtmx "github.com/falagansoftware/go-htmx/internal"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) FindUserByUid(ctx context.Context, uid string) (*gohtmx.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := findUserByUid(ctx, tx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) FindUsers(ctx context.Context, filters *gohtmx.UserFilters) ([]*gohtmx.User, error) {
	tx, err := u.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users, _, err := findUsers(ctx, tx, filters)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Helpers

func findUserByUid(ctx context.Context, tx *Tx, uid string) (*gohtmx.User, error) {
	users, _, err := findUsers(ctx, tx, &gohtmx.UserFilters{Uid: &uid})

	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, &gohtmx.Error{Code: gohtmx.ENOTFOUND, Message: "User not found"}
	}
	return users[0], nil

}

func findUsers(ctx context.Context, tx *Tx, filter *gohtmx.UserFilters) (u []*gohtmx.User, n int, e error) {
	// Where clause based on filters props
	where, args := []string{"1=1"}, []interface{}{}
	orderBy := "name"
	direction := "ASC"

	if v := filter.Uid; v != nil {
		condition := fmt.Sprintf("uid = '%v'", *v)
		where, args = append(where, condition), append(args, *v)
	}

	if v := filter.Name; v != nil {
		condition := fmt.Sprintf("name = '%v'", *v)
		where, args = append(where, condition), append(args, *v)
	}

	if v := filter.Surname; v != nil {
		condition := fmt.Sprintf("surname = '%v'", *v)
		where, args = append(where, condition), append(args, *v)
	}

	if v := filter.Email; v != nil {
		condition := fmt.Sprintf("email = '%v'", *v)
		where, args = append(where, condition), append(args, *v)
	}

	if v := filter.Active; v {
		condition := fmt.Sprintf("active = '%v'", v)
		where, args = append(where, condition), append(args, true)
	}

	if v := filter.Sort; v != "" {
		orderBy = fmt.Sprintf("ORDER BY %v", filter.Sort)
	}

	if v := filter.Order; v != "" {
		direction = v
	}
	// Execute query
	query := `SELECT uid, name, surname, email, active, created_at, updated_at FROM users WHERE ` + strings.Join(where, " AND ") + ` ` + orderBy + ` ` + direction + ` ` + FormatLimitOffset(filter.Limit, filter.Offset)
	rows, err := tx.QueryContext(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	// Map rows to struct
	users := make([]*gohtmx.User, 0)

	for rows.Next() {
		var user gohtmx.User
		err := rows.Scan(&user.Uid, &user.Name, &user.Surname, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	// Check rows error
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, len(users), nil

}
