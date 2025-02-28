package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type todo struct {
	Item string
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var todos []string
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res)
		todos = append(todos, res)
	}
	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}
	fmt.Printf("%v", newTodo)
	if newTodo.Item != "" {
		_, err := db.Exec("INSERT into todos VALUES ($1)", newTodo.Item)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	olditem := c.Query("olditem")
	newitem := c.Query("newitem")
	db.Exec("UPDATE todos SET item=$1 WHERE item=$2", newitem, olditem)
	return c.Redirect("/")
}

// func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
// 	todoToDelete := c.Query("item")
// 	db.Exec("DELETE from todos WHERE item=$1", todoToDelete)
// 	return c.SendString("deleted")
// }

func main() {

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	endpoint := os.Getenv("DB_ENDPOINT")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/my-db?sslmode=disable", username, password, endpoint, dbPort)

	// connStr := "postgresql://arsh:@localhost:5432/todos?sslmode=disable"

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	// app.Delete("/delete", func(c *fiber.Ctx) error {
	// 	return deleteHandler(c, db)
	// })

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Static("/", "./public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
