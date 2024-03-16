package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgressConfig struct {
	Host string
	Port string
	User string
	Password string
	Database string
	SSLModel string
}

func (cfg PostgressConfig) toString() string{
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLModel)

}

func main() {
	cfg := PostgressConfig{
		Host: "localhost",
		Port: "5432",
		User: "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLModel: "disable",
	}

	db, err := sql.Open("pgx", cfg.toString())
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("connected")

	// Create a table....
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		description TEXT
	);
		
	`)

	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created!")


	//Insert some data
	// name:= "Test User"
	// email := "test@email.com"

	// row := db.QueryRow(`
	//  INSERT INTO users(name, email)
	//  VALUES ($1, $2) RETURNING id;
	// `, name, email)

	// var id int
	// // pass id memory address in and row.scan will set value there
	// err = row.Scan(&id) // if did not have this scan, would need to check for errors another way (i.e row.Err())
	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Println("User created! id = ", id)

	// id := 1
	// row := db.QueryRow(`
	// 	SELECT name, email
	// 	FROM users
	// 	WHERE id=$1;
	// 	`, id)

	// var name, email string
	// err = row.Scan(&name, &email)
	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Printf("Usr information: name=%s, email=%s\n", name, email)

	userID := 1
	for i:=1; i<=5; i++ {
		amount := i * 100
		desc := fmt.Sprintf("Fake order #%d", i)
		_, err = db.Exec(`
		INSERT INTO orders(user_id, amount, description)
		VALUES ($1, $2, $3);
		`, userID, amount, desc)
	}

	if err != nil{
		panic(err)
	}

	fmt.Println("Created fake orders.")
}
