package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todo "github.com/serg495/go-to-do/entity"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (repo *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := repo.db.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (repo *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("SELECT tl.* FROM %s AS tl INNER JOIN %s AS ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	err := repo.db.Select(&lists, query, userId)

	return lists, err
}

func (repo *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf("SELECT tl.* FROM %s AS tl INNER JOIN %s AS ul on tl.id = ul.list_id "+
		"WHERE ul.user_id = $1 AND tl.id = $2",
		todoListsTable, usersListsTable)
	err := repo.db.Get(&list, query, userId, listId)

	return list, err
}

func (repo *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setSubQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s AS tl SET %s FROM %s AS ul WHERE tl.id = ul.list_id AND tl.id=$%d AND ul.user_id=$%d",
		todoListsTable, setSubQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("update query: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := repo.db.Exec(query, args...)

	return err
}

func (repo *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s AS tl USING %s AS ul WHERE tl.id = ul.list_id "+
		"AND ul.user_id = $1 AND tl.id = $2",
		todoListsTable, usersListsTable)
	_, err := repo.db.Exec(query, userId, listId)

	return err
}
