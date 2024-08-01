package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Info().Msg("Імпортовано головний пакет")
}

func main() {
	storage, err := newStorage(os.Getenv("POSTGRES_CONN_STR"))
	if err != nil {
		log.Fatal().Msgf("Помилка створення бази даних: %v", err)
	}

	auth := &Auth{s: storage}
	plans := &PlanResource{s: storage}
	users := &UserResource{s: storage}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("GET /users", users.GetAllUsers)

	mux.HandleFunc("GET /plans", auth.checkAuth(plans.GetAllPlans))
	mux.HandleFunc("POST /plans", auth.checkAuth(plans.CreatePlan))
	mux.HandleFunc("DELETE /plans/{id}", auth.checkAuth(plans.DeletePlan))
	mux.HandleFunc("PUT /plans/{id}", auth.checkAuth(plans.UpdatePlan))

	fmt.Println("Слухаєм :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Невдала спроба створити та прослухати 8080", err)
	}
}

type PlanResource struct {
	s *Storage
}

func (p *PlanResource) GetAllPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := p.s.GetAllPlans()
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка отримання планів з бази даних:, %v", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(plans)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка кодування планів у JSON: %v", err), http.StatusInternalServerError)
	}
}

func (p *PlanResource) GetPlanByID(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	planID, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Неправильний ID плану: %v", err), http.StatusBadRequest)
		return
	}

	plan, err := p.s.GetPlanByID(planID)
	if err != nil {

		http.Error(w, fmt.Sprintf("Помилка отримання плану з бази даних: %v", err), http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(plan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка кодування плану у JSON: %v", err), http.StatusInternalServerError)
	}
}

func (p *PlanResource) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var plan Plan

	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка декодування запиту: %v", err), http.StatusBadRequest)
		return
	}

	planID, err := p.s.CreatePlan(plan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка створення плану в базі даних: %v", err), http.StatusInternalServerError)
		return
	}

	plan.ID = planID

	err = json.NewEncoder(w).Encode(plan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка кодування плану у JSON`: %v", err), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (p *PlanResource) DeletePlan(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	planID, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Неправильний ID плану: %v", err), http.StatusBadRequest)
		return
	}

	_, err = p.s.GetPlanByID(planID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка перевірки наявності плану: %v", err), http.StatusInternalServerError)
		return
	}

	err = p.s.DeletePlan(planID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка видалення плану: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *PlanResource) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")
	planID, err := strconv.Atoi(idValue)
	if err != nil {
		http.Error(w, "Неправильний ID плану", http.StatusBadRequest)
		return
	}

	var updatedPlan Plan
	err = json.NewDecoder(r.Body).Decode(&updatedPlan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка декодування запиту: %v", err), http.StatusBadRequest)
		return
	}
	_, err = p.s.GetPlanByID(planID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка перевірки наявності плану: %v", err), http.StatusInternalServerError)
		return
	}

	err = p.s.UpdatePlan(planID, updatedPlan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка оновлення плану: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type UserResource struct {
	s *Storage
}

func (ur *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка декодування запиту: %v", err), http.StatusBadRequest)
		return
	}

	exists, err := ur.s.CheckUsernameExists(user.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка перевірки користувача в базі даних: %v", err), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Користувач з таким ім'ям вже існує", http.StatusConflict)
		return
	}

	userID, err := ur.s.CreateUser(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка створення користувача в базі даних: %v", err), http.StatusInternalServerError)
		return
	}

	user.ID = userID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка кодування користувача у JSON: %v", err), http.StatusInternalServerError)
	}
}

func (ur *UserResource) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ur.s.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка отримання юзерів з бази даних: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Помилка кодування юзерів у JSON: %v", err), http.StatusInternalServerError)
	}
}
