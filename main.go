package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// func main() {
// 	mux := http.NewServeMux()
// 	s := NewStorage()
// 	plans := PlanResource{
// 		s: NewStorage(),
// 	}

// 	users := UserResource{
// 		s: s,
// 	}
// 	auth := Auth{
// 		s: s,
// 	}
// 	mux.HandleFunc("POST /users", users.CreateUser)
// 	mux.HandleFunc("GET /plans", auth.checkAuth(plans.GetAllPlans))
// 	mux.HandleFunc("POST /plans", auth.checkAuth(plans.CreatePlan))
// 	mux.HandleFunc("DELETE /plans/{id}", auth.checkAuth(plans.DeletePlan))
// 	mux.HandleFunc("PUT /plans/{id}", auth.checkAuth(plans.UpdatePlan))

// 	fmt.Println("Слухаєм :8080")
// 	if err := http.ListenAndServe(":8080", mux); err != nil {
// 		fmt.Println("Невдала спроба створити та прослухати 8080", err)
// // 	}
// // }

// type PlanResource struct {
// 	s *Storage
// }

// func (p *PlanResource) GetAllPlans(w http.ResponseWriter, r *http.Request) {
// 	plans := p.s.GetAllPlans()

// 	err := json.NewEncoder(w).Encode(plans)
// 	if err != nil {
// 		fmt.Println("ПОмилка кодування в JSON", err)
// 		return
// 	}
// }

// func (p *PlanResource) CreatePlan(w http.ResponseWriter, r *http.Request) {
// 	var plan Plan

// 	err := json.NewDecoder(r.Body).Decode(&plan)
// 	if err != nil {
// 		fmt.Println("ПОмилка декодування", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	plan.ID = p.s.CreatePlan(plan)

// 	err = json.NewEncoder(w).Encode(plan)
// 	if err != nil {
// 		fmt.Println("ПОмилка кодування в JSON", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// }

// func (p *PlanResource) DeletePlan(w http.ResponseWriter, r *http.Request) {
// 	idValue := r.PathValue("id")
// 	planId, err := strconv.Atoi(idValue)
// 	if err != nil {
// 		fmt.Println("Не існує нічого з таким id")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return

// 	}
// 	_, ok := p.s.GetPlanById(planId)
// 	if !ok {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	p.s.DeletePlanById(planId)
// }

// func (p *PlanResource) UpdatePlan(w http.ResponseWriter, r *http.Request) {
// 	idValue := r.PathValue("id")
// 	planId, err := strconv.Atoi(idValue)
// 	if err != nil {
// 		fmt.Println("Не існує нічого з таким id")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return

// 	}
// 	_, ok := p.s.GetPlanById(planId)
// 	if !ok {
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}

// 	var UpdatedPlan Plan
// 	err = json.NewDecoder(r.Body).Decode(&UpdatedPlan)

// 	if err != nil {
// 		fmt.Println("ПОмилка декодування JSON", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	success := p.s.ChangePlan(planId, UpdatedPlan)
// 	if !success {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)

// }

// type UserResource struct {
// 	s *Storage
// }

// func (ur *UserResource) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var user User

// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		fmt.Println("ПОмилка декодування", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	ok := ur.s.CreateUser(user)
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// }

func init() {
	log.Info().Msgf("Imported main package")
}
func main() {
	// db, err := sql.Open("postgres", "postgres://admin:donotcrackplease@localhost:5433/planner?sslmode=disable")
	// if err != nil {
	// 	log.Fatal().Msgf("Opening database %v", err)
	// }

	// if err := db.Ping(); err != nil {
	// 	log.Fatal().Msgf("Pinging database %v", err)

	// }
	// rows, err := db.Query("SELECT id, name FROM users")

	// if err != nil {
	// 	log.Fatal().Msgf("Selecting users")
	// }

	// var users []user

	// for rows.Next() {
	// 	var u user
	// 	err := rows.Scan(&u.id, &u.name)
	// 	if err != nil {
	// 		log.Fatal().Msgf("Scanning rows : %v", err)
	// 	}

	// 	users = append(users, u)
	// }
	storage, err := newStorage("postgres://admin:donotcrackplease@localhost:5433/planner?sslmode=disable")
	if err != nil {
		log.Fatal().Msgf("Creating storage: %v", err)
	}

	users, err := storage.getAllUsers()

	if err != nil {
		log.Fatal().Msgf("getting all users: %v", err)

	}

	log.Info().Msgf("Got users: %v", users)

	err = storage.insertUser("Taras")

	if err != nil {
		log.Fatal().Msgf("Getting all users: %v", err)
	}

	log.Info().Msgf("Got users; %v", users)
}

type user struct {
	id   int
	name string
}

type storage struct {
	db *sql.DB
}

func newStorage(connString string) (*storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database %w", err)
	}

	return &storage{db: db}, nil
}
func (s *storage) getAllUsers() ([]user, error) {
	rows, err := s.db.Query("SELECT id, name FROM users")

	if err != nil {
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	var users []user

	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}

		users = append(users, u)
	}
	return users, nil
}

func (s *storage) insertUser(name string) error {
	_, err := s.db.Exec("INSERT INTO users(name) VALUES($1)", name)

	if err != nil {
		return fmt.Errorf("inserting user: %w", err)

	}
	return nil
}
