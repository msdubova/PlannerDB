package main

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func newStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("помилка підключення до бази даних: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("помилка перевірки підключення до бази даних: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllPlans() ([]Plan, error) {
	rows, err := s.db.Query("SELECT id, title, descriptio, complete FROM plans")
	if err != nil {
		return nil, fmt.Errorf("помилка отримання планів: %w", err)
	}

	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var p Plan
		err := rows.Scan(&p.ID, &p.Title, &p.Descriptio, &p.Complete)
		if err != nil {
			return nil, fmt.Errorf("помилка зчитування планів: %w", err)
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (s *Storage) GetPlanByID(id int) (Plan, error) {
	var plan Plan
	err := s.db.QueryRow("SELECT id, title, descriptio, complete FROM plans WHERE id = $1", id).Scan(&plan.ID, &plan.Title, &plan.Descriptio, &plan.Complete)
	if err != nil {
		return plan, fmt.Errorf("помилка отримання плану: %w", err)
	}
	return plan, nil
}

func (s *Storage) CreatePlan(plan Plan) (int, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO plans (title, descriptio, complete) VALUES ($1, $2, $3) RETURNING id", plan.Title, plan.Descriptio, plan.Complete).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("помилка вставки плану: %w", err)
	}
	return id, nil
}

func (s *Storage) UpdatePlan(id int, updatedPlan Plan) error {
	_, err := s.db.Exec("UPDATE plans SET title = $1, descriptio = $2, complete = $3 WHERE id = $4", updatedPlan.Title, updatedPlan.Descriptio, updatedPlan.Complete, id)
	if err != nil {
		return fmt.Errorf("помилка оновлення плану: %w", err)
	}

	return nil
}

func (s *Storage) DeletePlan(id int) error {
	result, err := s.db.Exec("DELETE FROM plans WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("помилка видалення плану: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("помилка перевірки кількості видалених рядків: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("план з ID %d не знайдено", id)
	}

	return nil
}

func (s *Storage) CreateUser(user User) (int, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("помилка створення користувача: %w", err)
	}
	return id, nil
}

func (s *Storage) GetUser(username string) (User, bool) {
	var user User
	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return user, false
	}
	return user, true
}

func (s *Storage) GetAllUsers() ([]User, error) {
	rows, err := s.db.Query("SELECT id, username, password FROM users")
	if err != nil {
		return nil, fmt.Errorf("помилка отримання юзерів: %v", err)
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("помилка зчитування юзерів: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("помилка при ітерації: %v", err)
	}

	return users, nil
}

func (s *Storage) CheckUsernameExists(username string) (bool, error) {
	var exists bool
	query := "SELECT exists (SELECT 1 FROM users WHERE username=$1)"
	err := s.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
