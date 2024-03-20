package main

import (
	"fmt"

	"bryce-ziemer/github.com/lenslocked/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("connected")

	// added to show that code in models.user.go is working
	us := models.UserService{
		DB: db,
	}

	user, err := us.Create("bob4@bob.com", "bob123")
	if err != nil {
		panic(err)

	}
	fmt.Println(user)

	// Create a table....
	// _, err = db.Exec(`
	// CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	email TEXT UNIQUE NOT NULL
	// );

	// CREATE TABLE IF NOT EXISTS orders (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	amount INT,
	// 	description TEXT
	// );

	// `)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Tables created!")

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

	// userID := 1
	// for i:=1; i<=5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)
	// 	_, err = db.Exec(`
	// 	INSERT INTO orders(user_id, amount, description)
	// 	VALUES ($1, $2, $3);
	// 	`, userID, amount, desc)
	// }

	// if err != nil{
	// 	panic(err)
	// }

	// fmt.Println("Created fake orders.")

	// 	type Order struct {
	// 		ID          int
	// 		UserId      int
	// 		Amount      int
	// 		Description string
	// 	}

	// 	var orders []Order
	// 	userId := 1
	// 	rows, err := db.Query(`
	// 		SELECT id, amount, description
	// 		FROM orders
	// 		WHERE user_id=$1;
	// 	`, userId)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer rows.Close()

	// 	for rows.Next() {
	// 		var order Order
	// 		order.UserId = userId
	// 		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		orders = append(orders, order)

	// 	}

	// 	// check for an error
	// 	err = rows.Err()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("Orders:", orders)

}
