package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func CreateTable(pool *pgxpool.Pool) error {
	stmt := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY ,
			email VARCHAR(100) UNIQUE NOT NULL,
            username VARCHAR(50) NOT NULL,
            mobile VARCHAR(50) NOT NULL,
            password VARCHAR(100) NOT NULL,
			birthdate VARCHAR(50),
			gender VARCHAR(50),
			height VARCHAR(50),
			weight VARCHAR(50),
			activity VARCHAR(50),
			notifications BOOLEAN DEFAULT FALSE, 
			role_id INTEGER,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
        );

		CREATE TABLE IF NOT EXISTS points (
			id SERIAL PRIMARY KEY, 
			description VARCHAR(100) NOT NULL, 
			point INTEGER NOT NULL, 
			created_at TIMESTAMP, 
			expired_at TIMESTAMP, 
			updated_at TIMESTAMP,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			is_available BOOLEAN NOT NULL,
			available_point numeric not null,
			is_used boolean default false,
			promocode_id bigint default 0

        );


		CREATE TABLE IF NOT EXISTS promocodes (
			id SERIAL PRIMARY KEY, 
			promocode VARCHAR(100) NOT NULL,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			count INTEGER NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS info_company (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			description VARCHAR(256) NOT NULL,
			image VARCHAR(256),
			created_at TIMESTAMP,
			updated_at TIMESTAMP
        );

		CREATE TABLE IF NOT EXISTS info_app (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			image VARCHAR(256),
			app_name VARCHAR(50),
			description VARCHAR(256) NOT NULL,
			version VARCHAR(50),
			release_date VARCHAR(50),
			created_at TIMESTAMP,
			updated_at TIMESTAMP
        );

		CREATE TABLE IF NOT EXISTS info_bonus (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			description VARCHAR(256) NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
        );

		CREATE TABLE IF NOT EXISTS info_contact (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
        );

		CREATE TABLE IF NOT EXISTS courses (
			id SERIAL PRIMARY KEY,
			name VARCHAR(256) NOT NULL,
			description VARCHAR(256) NOT NULL,
			image VARCHAR(256),
			points INTEGER,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_courses (
			id SERIAL PRIMARY KEY,
			is_started BOOLEAN NOT NULL,
			is_finished BOOLEAN NOT NULL,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE NOT NULL,
			created_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS products	(
			id serial primary key,
			name varchar(255) not null unique,
			image varchar(255) not null unique,
			price numeric not null,
			height varchar(255),
			size varchar(255),
			instruction varchar(255),
			description varchar(255) 
		);

		CREATE TABLE IF NOT EXISTS lessons (
			id SERIAL PRIMARY KEY,
			lesson_number INTEGER NOT NULL, 
			name VARCHAR(256) NOT NULL,
			description VARCHAR(256) NOT NULL,
			image VARCHAR(256),
			video VARCHAR(256),
			time INTERVAL,
			course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE NOT NULL,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS lesson_products (
			id SERIAL PRIMARY KEY,
			lesson_id INTEGER REFERENCES lessons(id) ON DELETE CASCADE NOT NULL,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE NOT NULL,
			created_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bookmarks (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE NOT NULL,
			lesson_id INTEGER REFERENCES lessons(id) ON DELETE CASCADE NOT NULL,
			created_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_watched_lessons (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE NOT NULL,
			lesson_id INTEGER REFERENCES lessons(id) ON DELETE CASCADE NOT NULL,
			is_watched BOOLEAN NOT NULL,
			created_at TIMESTAMP
		);

		
		
		CREATE TABLE IF NOT EXISTS recommended_products	(
			id serial primary key,
			product_id int references products (id) on delete cascade not null,
			recommended_product int references products (id) on delete cascade not null
		);

		CREATE TABLE IF NOT EXISTS carts (
			id serial primary key,
			user_id int references users (id) on delete cascade not null,
			product_id int references products (id) on delete cascade UNIQUE not null,
			quantity int default 1
		);

				
		CREATE TABLE IF NOT EXISTS orders
		(
			id serial primary key,
			user_id int references users (id) not null,
			total_amount numeric not null,
			created_date DATE DEFAULT CURRENT_DATE,
			bonus numeric default 0,
			overall numeric not null
		);


		create table IF NOT EXISTS order_items(
			id serial primary key,
			order_id  int references orders (id) on delete cascade not null,
			product_id int  references products (id) on delete cascade not null,
			quantity numeric not null,
			price numeric not null
		);

    `

	_, err := pool.Exec(context.Background(), stmt)
	if err != nil {
		return err
	}

	return nil
}

func CreateAdmin(pool *pgxpool.Pool) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), 10)
	if err != nil {
		fmt.Println("Failed to hash admin assword")
		return err
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := string(hashPassword)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	tx, err := pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	stmt := `
		INSERT INTO users (email, username, mobile, password, role_id, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO NOTHING;
		`
	_, err = tx.Exec(context.Background(), stmt, email, "admin", "", password, 1, currentTime)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
