package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lib/pq"
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

func (h *handler) CreateExpense(c echo.Context) error {

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
