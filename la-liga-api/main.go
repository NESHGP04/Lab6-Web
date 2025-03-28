package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"home_team"`
	AwayTeam    string `json:"away_team"`
	ScoreHome   int    `json:"score_home"`
	ScoreAway   int    `json:"score_away"`
	YellowCards int    `json:"yellow_cards"`
	RedCards    int    `json:"red_cards"`
	ExtraTime   int    `json:"extra_time"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=yourpassword dbname=la_liga sslmode=disable")
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/matches", createMatch).Methods("POST")
	r.HandleFunc("/api/matches/{id}", getMatch).Methods("GET")
	r.HandleFunc("/api/matches/{id}", updateMatch).Methods("PUT")
	r.HandleFunc("/api/matches/{id}", deleteMatch).Methods("DELETE")
	r.HandleFunc("/api/matches/{id}/goal", registerGoal).Methods("POST")
	r.HandleFunc("/api/matches/{id}/yellow-card", registerYellowCard).Methods("POST")
	r.HandleFunc("/api/matches/{id}/red-card", registerRedCard).Methods("POST")
	r.HandleFunc("/api/matches/{id}/extra-time", setExtraTime).Methods("POST")

	log.Println("✅ Servidor iniciado en el puerto 8080")
	http.ListenAndServe(":8080", r)
}

func createMatch(w http.ResponseWriter, r *http.Request) {
	var match Match
	json.NewDecoder(r.Body).Decode(&match)
	_, err := db.Exec("INSERT INTO matches (home_team, away_team, score_home, score_away, yellow_cards, red_cards, extra_time) VALUES ($1, $2, 0, 0, 0, 0, 0)", match.HomeTeam, match.AwayTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var match Match
	err := db.QueryRow("SELECT * FROM matches WHERE id = $1", id).Scan(&match.ID, &match.HomeTeam, &match.AwayTeam, &match.ScoreHome, &match.ScoreAway, &match.YellowCards, &match.RedCards, &match.ExtraTime)
	if err != nil {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(match)
}

func updateMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var match Match
	json.NewDecoder(r.Body).Decode(&match)
	_, err := db.Exec("UPDATE matches SET home_team=$1, away_team=$2 WHERE id=$3", match.HomeTeam, match.AwayTeam, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteMatch(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec("DELETE FROM matches WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func registerGoal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec("UPDATE matches SET score_home = score_home + 1 WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func registerYellowCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec("UPDATE matches SET yellow_cards = yellow_cards + 1 WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func registerRedCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec("UPDATE matches SET red_cards = red_cards + 1 WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func setExtraTime(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var extraTime struct {
		ExtraTime int `json:"extra_time"`
	}
	json.NewDecoder(r.Body).Decode(&extraTime)
	_, err := db.Exec("UPDATE matches SET extra_time = $1 WHERE id = $2", extraTime.ExtraTime, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
