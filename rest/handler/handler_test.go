//go:build unit
// +build unit

package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	//"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/lib/pq"
)

func TestCreateExpense(t *testing.T) {

	//json input
	expenseJSON := `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath", "tags": ["food", "beverage"]}`

	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expenseJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")
	rec := httptest.NewRecorder()

	tags1 := []string{"food", "beverage"}

	mockedRow := sqlmock.NewRows([]string{"id"}).AddRow("1")

	db, mock, err := sqlmock.New()

	mockedSql := "INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id"
	mock.ExpectQuery(regexp.QuoteMeta(mockedSql)).WithArgs("strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array(tags1)).WillReturnRows((mockedRow))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	h := handler{db}
	c := e.NewContext(req, rec)

	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.CreateExpense(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}


func TestGetExpenses(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")

	rec := httptest.NewRecorder()

	tags1 := []string{"food","beverage"}
	//tagss := fmt.Sprintf("%v", pq.Array(tags))

	//fmt.Printf("type of a is %T\n", tagss)
	//fmt.Printf("%s", tagss)

	expsMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note","tags"}).
		AddRow("1", "strawberry smoothie", "79", "night market promotion discount 10 bath", pq.Array(tags1))

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id=1").WillReturnRows(expsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	// Act
	err = h.GetExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}