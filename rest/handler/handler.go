package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/lib/pq"

	"encoding/base64"
)

type handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

const (
	bearer = "bearer"
	username = "apidesign" //test basic authen purpose
	password = "45678" //test basic authen purpose
)

func checkAuthorization(c echo.Context) error {


	auth := c.Request().Header.Get(echo.HeaderAuthorization)

	l := len(bearer)

	if len(auth) > l+1 {
		// Invalid base64 shouldn't be treated as error
		// instead should be treated as invalid client input
		b, err := base64.StdEncoding.DecodeString(auth[l+1:])
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}

		cred := string(b)
		for i := 0; i < len(cred); i++ {
			if cred[i] == ':' {
				// Verify credentials
				u := cred[:i]
				p := cred[i+1:]

				log.Printf("Username: password=%s:%s", username, password)
				//Just testing purpose, hardcode username and password 
				if u != username && p != password {

					return c.JSON(http.StatusForbidden, Err{Message: "You are not authorized to use this path"}) 
				}

			}
		}
	}

	return nil

}

func (h *handler) CreateExpense(c echo.Context) error {

	checkAuthorization(c)

	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	pgTagsarr := pq.Array(exp.Tags)

	row := h.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id;", exp.Title, exp.Amount, exp.Note, pgTagsarr)
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	//return c.JSON(http.StatusCreated, exp)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(exp)
}
