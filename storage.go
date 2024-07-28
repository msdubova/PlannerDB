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
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllPlans() ([]Plan, error) {
	rows, err := s.db.Query("SELECT * FROM plans")

	if err != nil {
		return nil, fmt.Errorf("selecting plans: %w", err)
	}
	defer rows.Close()
	var plans []Plan

	for rows.Next() {
		var p Plan
		err := rows.Scan(&p.ID, &p.Descriptio, &p.Title, &p.Complete)

		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}

		plans = append(plans, p)
	}
	return plans, nil
}

func (s *Storage) CreatePlan(plan Plan) (int, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO plans (title, descriptio, complete) VALUES ($1, $2, $3) RETURNING id", plan.Title, plan.Descriptio, plan.Complete).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Помилка вставки плану: %w", err)
	}
	return id, nil
}

func (s *Storage) GetPlanById(planId int) (Plan, bool) {
	var plan Plan
	err := s.db.QueryRow("SELECT id, title, descriptio, complete FROM plans WHERE id = $1", planId).Scan(&plan.ID, &plan.Title, &plan.Descriptio, &plan.Complete)
	if err != nil {
		return plan, false
	}
	return plan, true
}

func (s *Storage) UpdatePlan(id int, updatedPlan Plan) bool {
	_, err := s.db.Exec("UPDATE plans SET title = $1, descriptio = $2, complete = $3 WHERE id = $4", updatedPlan.Title, updatedPlan.Descriptio, updatedPlan.Complete, id)
	if err != nil {
		fmt.Println("Помилка оновлення плану", err)
		return false
	}
	return true
}

func (s *Storage) DeletePlan(id int) {
	_, err := s.db.Exec("DELETE FROM plans WHERE id = $1", id)
	if err != nil {
		fmt.Println("Помилка видалення плану", err)
	}
}

func (s *Storage) CreateUser(user User) bool {
	_, err := s.db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", user.ID, user.Username, user.Password)
	if err != nil {
		fmt.Println("Помилка створення користувача", err)
		return false
	}
	return true
}

func (s *Storage) GetUser(username string) (User, bool) {
	var user User
	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return user, false
	}
	return user, true
}
