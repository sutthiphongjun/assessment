//go:build integration
// +build integration

package handler

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

// EXP01
func TestITCreateExpense(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/expensedb?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.POST("/expenses", h.CreateExpense)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
			log.Println("Just Loop wait port open.....")
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath","tags": ["food", "beverage"]}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITGetExpenses(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/expensedb?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.GET("/expenses/:id", h.GetExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
			log.Println("Just Loop wait port open.....")
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")

	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}


//EXP03
func TestITPutExpense(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/expensedb?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.PUT("/expenses/:id", h.UpdateExpense)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
			log.Println("Just Loop wait port open.....")
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := `{"title": "apple smoothie","amount": 89,"note": "no discount","tags": ["beverage"]}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")	
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}



func TestITCreateExpenseID2(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/expensedb?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.POST("/expenses", h.CreateExpense)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
			log.Println("Just Loop wait port open.....")
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	reqBody := `{"title": "iPhone 14 Pro Max 1TB","amount": 66900,"note": "birthday gift from my love","tags": ["gadget"]}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")	
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "{\"id\":2,\"title\":\"iPhone 14 Pro Max 1TB\",\"amount\":66900,\"note\":\"birthday gift from my love\",\"tags\":[\"gadget\"]}"


	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITListExpenses(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/expensedb?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := NewApplication(db)

		e.GET("/expenses", h.ListExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
			log.Println("Just Loop wait port open.....")
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//apidesign:45678
	req.Header.Set(echo.HeaderAuthorization, "Bearer YXBpZGVzaWduOjQ1Njc4")		
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "[{\"id\":1,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]},{\"id\":2,\"title\":\"iPhone 14 Pro Max 1TB\",\"amount\":66900,\"note\":\"birthday gift from my love\",\"tags\":[\"gadget\"]}]"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}


