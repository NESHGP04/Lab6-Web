package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"home_team"`
	AwayTeam    string `json:"away_team"`
	MatchDate   string `json:"match_date"`
	ScoreHome   int    `json:"score_home"`
	ScoreAway   int    `json:"score_away"`
	YellowCards int    `json:"yellow_cards"`
	RedCards    int    `json:"red_cards"`
	ExtraTime   int    `json:"extra_time"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "74200474M"
	dbname   = "LigaDB"
)

var db *sql.DB

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// func getUser(c echo.Context) error {
// 	id := c.Param("id")
// 	return c.String(http.StatusOK, id)
// }

func main() {

	e := echo.New()

	// Habilitar CORS
	e.Use(middleware.CORS())

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e.GET("/api/matches", getAllMatches)
	e.GET("/api/matches/:id", getMatchByID)
	e.POST("/api/matches", createMatch)
	e.PUT("/api/matches/:id", updateMatch)
	e.DELETE("/api/matches/:id", deleteMatch)

	e.PATCH("/api/matches/:id/goals", updateGoals)
	e.PATCH("/api/matches/:id/yellowcards", addYellowCard)
	e.PATCH("/api/matches/:id/redcards", addRedCard)
	e.PATCH("/api/matches/:id/extratime", updateExtraTime)

	// e.GET("/ping", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "PONG!")
	// })

	// e.GET("user/:id", getUser)

	e.Logger.Fatal(e.Start("0.0.0.0:8070"))
}

func getAllMatches(c echo.Context) error {
	rows, err := db.Query(`SELECT id, home_team, away_team, match_date, score_home, score_away, yellow_cards, red_cards, extra_time FROM "Matches"`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate, &m.ScoreHome, &m.ScoreAway, &m.YellowCards, &m.RedCards, &m.ExtraTime); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		matches = append(matches, m)
	}

	return c.JSON(http.StatusOK, matches)
}

func getMatchByID(c echo.Context) error {
	id := c.Param("id")
	var m Match
	err := db.QueryRow(`SELECT id, home_team, away_team, match_date, score_home, score_away, yellow_cards, red_cards, extra_time FROM "Matches" WHERE id=$1`, id).
		Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate, &m.ScoreHome, &m.ScoreAway, &m.YellowCards, &m.RedCards, &m.ExtraTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, "Partido no encontrado")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, m)
}

// Crear un nuevo partido
func createMatch(c echo.Context) error {
	var m Match
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := `INSERT INTO "Matches" (home_team, away_team, match_date, score_home, score_away, yellow_cards, red_cards, extra_time)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err := db.QueryRow(query, m.HomeTeam, m.AwayTeam, m.MatchDate, m.ScoreHome, m.ScoreAway, m.YellowCards, m.RedCards, m.ExtraTime).Scan(&m.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, m)
}

// Actualizar un partido existente
func updateMatch(c echo.Context) error {
	id := c.Param("id")
	var m Match
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := `UPDATE "Matches" SET home_team=$1, away_team=$2, match_date=$3, score_home=$4, score_away=$5, yellow_cards=$6, red_cards=$7, extra_time=$8 WHERE id=$9`
	res, err := db.Exec(query, m.HomeTeam, m.AwayTeam, m.MatchDate, m.ScoreHome, m.ScoreAway, m.YellowCards, m.RedCards, m.ExtraTime, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "Partido no encontrado")
	}

	return c.JSON(http.StatusOK, "Partido actualizado correctamente")
}

// Eliminar un partido
func deleteMatch(c echo.Context) error {
	id := c.Param("id")

	res, err := db.Exec(`DELETE FROM "Matches" WHERE id=$1`, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "Partido no encontrado")
	}

	return c.JSON(http.StatusOK, "Partido eliminado correctamente")
}

// Actualizar goles de un partido
func updateGoals(c echo.Context) error {
	id := c.Param("id")
	var match Match
	if err := c.Bind(&match); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := `UPDATE "Matches" SET score_home=$1, score_away=$2 WHERE id=$3`
	res, err := db.Exec(query, match.ScoreHome, match.ScoreAway, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return checkRowsAffected(res, c, "Goles actualizados correctamente")
}

// Registrar una tarjeta amarilla
func addYellowCard(c echo.Context) error {
	id := c.Param("id")

	query := `UPDATE "Matches" SET yellow_cards = yellow_cards + 1 WHERE id=$1`
	res, err := db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return checkRowsAffected(res, c, "Tarjeta amarilla registrada correctamente")
}

// Registrar una tarjeta roja
func addRedCard(c echo.Context) error {
	id := c.Param("id")

	query := `UPDATE "Matches" SET red_cards = red_cards + 1 WHERE id=$1`
	res, err := db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return checkRowsAffected(res, c, "Tarjeta roja registrada correctamente")
}

// Registrar tiempo extra
func updateExtraTime(c echo.Context) error {
	id := c.Param("id")
	var match Match
	if err := c.Bind(&match); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	query := `UPDATE "Matches" SET extra_time=$1 WHERE id=$2`
	res, err := db.Exec(query, match.ExtraTime, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return checkRowsAffected(res, c, "Tiempo extra actualizado correctamente")
}

// Función para verificar si se actualizó correctamente
func checkRowsAffected(res sql.Result, c echo.Context, successMsg string) error {
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, "Partido no encontrado")
	}
	return c.JSON(http.StatusOK, successMsg)
}
