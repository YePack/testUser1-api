package users

import (
	"errors"
	"fmt"
	"github.com/yepack/testUser1-api/datasources/mysql/users_db"
	"github.com/yepack/testUtils-api/logger"

	"github.com/yepack/testUser1-api/utils/mysql_utils"
	"github.com/yepack/testUtils-api/rest_errors"
	"strings"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, " +
		"last_name, email, date_created, status, password)VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser = " SELECT ID, first_name, last_name, " +
		"email, date_created, status FROM users WHERE id =?;"
	queryUpdateUser = "UPDATE users SET first_name=?," +
		" last_name=?, email=? WHERE id=?;"
	queryDeleteUser   = "DELETE FROM users WHERE id=?;"
	queryFindByStatus = "SELECT id,first_name, last_name," +
		"email, date_created, status  FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id,first_name, last_name, email," +
		"date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *User) Save() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user", err)
		return rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName,
		user.Email, user.DateCreate, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating "+
			"a new user", err)
		return rest_errors.NewInternalServerError("error when trying to get last insert "+
			"id after creating a new user", errors.New("database error"))
	}
	user.Id = userId
	return nil
}

func (user *User) Get() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user", err)
		return rest_errors.NewInternalServerError("error when trying to prepare get", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.DateCreate, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error when trying to get user by id", errors.New("database error"))
	}
	return nil
}

func (user *User) Update() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare update user statement", errors.New("database error"))
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare delete user statement", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users "+
			"by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to prepare find users "+
			"by status statement", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find users by status", errors.New("database error"))
	}

	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName,
			&user.Email, &user.DateCreate, &user.Status); err != nil {
			logger.Error("error when scan user row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when scan user row into user struct", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(
			fmt.Sprintf("No users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email&pass", err)
		return rest_errors.NewInternalServerError("error when trying to prepare get user by email&pass", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.DateCreate, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email&pass", getErr)
		return rest_errors.NewInternalServerError("error when trying to get user by email&pass", errors.New("database error"))
	}
	return nil
}
