package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/lib/pq"

	"encoding/base64"
	"errors"
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

func checkAuthorization(c echo.Context)  error {

	auth := c.Request().Header.Get(echo.HeaderAuthorization)

	l := len(bearer)

	if len(auth) > l+1 {
		// Invalid base64 shouldn't be treated as error
		// instead should be treated as invalid client input
		b, err := base64.StdEncoding.DecodeString(auth[l+1:])
		if err != nil {
			c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
			return errors.New("No have Authorization header")
		}

		cred := string(b)
		for i := 0; i < len(cred); i++ {
			if cred[i] == ':' {
				// Verify credentials
				u := cred[:i]
				p := cred[i+1:]

				//log.Printf("Username: password=%s:%s", username, password)
				//Just testing purpose, hardcode username and password 
				if u != username && p != password {

					c.JSON(http.StatusForbidden, Err{Message: "You are not authorized to use this path"})
					return errors.New("No have Authorization header")
				}

			}
		}

		//pass check
	}else{
		//error. No have Authorization header
		c.JSON(http.StatusForbidden, Err{Message: "You are not authorized to use this path. Please input token in http header"})
		return errors.New("No have Authorization header")

	}

	//pass check
	log.Println("PASS check")
	return nil

}

func (h *handler) CreateExpense(c echo.Context) error {

	resultcheck := checkAuthorization(c)

	if resultcheck != nil {
		log.Println("FOUND ERROR. EXIT")
		return nil
	}

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

func (h *handler) GetExpenses(c echo.Context) error {

	resultcheck := checkAuthorization(c)
	if resultcheck != nil {
		log.Println("FOUND ERROR. EXIT")
		return nil
	}

	uid := c.Param("id")

	//id MUST BE INTEGER
	_, err_int := strconv.Atoi(uid)
	if err_int != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err_int.Error()})	
	}

	rows, err := h.DB.Query("SELECT * FROM expenses WHERE id=" + uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	defer rows.Close()

	//var nn = []Expense{}
	var n = Expense{}

	var id, amount int
	var title, note string
	var tags []string	

	if rows.Next() {
		err := rows.Scan(&id, &title, &amount, &note, pq.Array(&tags))
		if err != nil {
			log.Fatal(err)
		}


		n.ID = id
		n.Title = title
		n.Amount = amount
		n.Note = note
		n.Tags = tags

		//nn = append(nn, n)
	}

	return c.JSON(http.StatusOK, n)
}

func (h *handler) UpdateExpense(c echo.Context) error {
	
	resultcheck := checkAuthorization(c)
	if resultcheck != nil {
		log.Println("FOUND ERROR. EXIT")
		return nil
	}

	eid := c.Param("id")

	//id MUST BE INTEGER
	eid_int, err_int := strconv.Atoi(eid)
	if err_int != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err_int.Error()})	
	}

	exp := Expense{}
	err := c.Bind(&exp); if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	pgTagsarr := pq.Array(exp.Tags)

	stmtupdate, errupdate := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")

	if errupdate != nil {
		log.Fatal("can't prepare statment update", errupdate)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	exp.ID = eid_int

	if _, err := stmtupdate.Exec(exp.ID, exp.Title, exp.Amount, exp.Note, pgTagsarr); err != nil {
		log.Fatal("error execute update ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}


	return c.JSON(http.StatusOK, exp)
}