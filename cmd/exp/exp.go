package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type Post struct {
		ID         int
		UserID     int
		Content    string
		ParentPost bool
	}

	type User struct {
		ID       int
		Username string
	}

	// _, err = db.Exec(`
	// CREATE TABLE IF NOT EXISTS posts (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	content TEXT NOT NULL,
	// 	parent_post BOOLEAN NOT NULL
	// );

	// CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	email TEXT NOT NULL
	// );`)
	// if err != nil {
	// 	panic(err)
	// }

	userID := 1
	var posts []Post

	row := db.QueryRow(`Select id, username from users WHERE id=$1`, userID)
	var user User
	user.ID = userID
	err = row.Scan(&user.ID, &user.Username)
	if err != nil {
		panic(err)
	}

	fmt.Printf("user %v", user)

	rows, err := db.Query(`Select id, user_id, content, parent_post from posts`)
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ParentPost)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	fmt.Println("posts: ", posts)

	// _, err = db.Exec(`INSERT INTO posts(user_id, content, parent_post) VALUES($1, $2, $3);`, userId, "hello world!", true)
	// if err != nil {
	// 	panic(err)
	// }

	// var orders []Order
	// userID := 1

	// rows, err := db.Query(`
	// 	SELECT id, amount, description
	// 	FROM orders
	// 	WHERE user_id=$1`, userID)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var order Order
	// 	order.UserID = userID
	// 	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("order:", orders)

	// userId := 1

	// for i := 1; i < 6; i++ {
	// 	amount := i * 100
	// 	description := fmt.Sprintf("a sled #%d", i)
	// 	_, err = db.Exec(`INSERT INTO orders(user_id, amount, description)
	// 		VALUES($1, $2, $3);
	// 	`, userId, amount, description)
	// }
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	email TEXT NOT NULL
	//   );

	//   CREATE TABLE IF NOT EXISTS orders (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	amount INT,
	// 	description TEXT
	//   );`)
	// if err != nil {
	// 	panic(err)
	// }

	// name := "Jon Calhoun"
	// email := "jon@calhoun.io"
	// _, err = db.Exec(`
	// INSERT INTO users(name, email)
	// VALUES($1, $2);`, name, email)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("User created.")
	// fmt.Println("Tables created.")
}
