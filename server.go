package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	//"net/http"

	"os"
	"os/signal"
	"time"

	"context"
	"net/http"
	//"strings"

	//"context"
	//"os/signal"
	//"syscall"

	//"rest/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/sutthiphongjun/assessment/rest/handler"
)

var db *sql.DB

func main() {

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	fmt.Println("DATABASE_URL:", os.Getenv("DATABASE_URL"))

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	} else {
		log.Println(db)
	}

	defer db.Close()

	log.Println("okay")

	createTb := `
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		);	 
	 `

	rs, err2 := db.Exec(createTb)

	if err2 != nil {
		log.Fatal("Create table error", err2)
	}

	rowseffected, _ := rs.RowsAffected()
	if rowseffected == 0 {
		fmt.Println("Success create table expenses")
	}

	h := handler.NewApplication(db)

	e := echo.New()

	//app log
	e.Use(middleware.Logger())

	e.GET("/expenses", h.ListExpenses)
	e.GET("/expenses/:id", h.GetExpenses)
	e.POST("/expenses", h.CreateExpense)
	e.PUT("/expenses/:id", h.UpdateExpense)

	// Intentionally, not setup database at this moment so we ignore feature to access database
	// e.GET("/news", h.ListNews)
	serverPort := os.Getenv("PORT")
	//e.Logger.Fatal(e.Start(serverPort))

	//graceful shutdown
	go func() {
		if err := e.Start(serverPort); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
