package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shyinyong/go-tcp-test/api"
	"github.com/shyinyong/go-tcp-test/config"
	"github.com/shyinyong/go-tcp-test/db/mysql"
	"log"
)

func main() {
	// Config env initialize
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initialize database connections
	store := mysql.NewDB(&cfg)
	defer store.Close()

	// find db
	st, err := studentByID(store, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("studentByID record: %v \n", st)

	students, err := fetchStudents(store)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("fetchStudents count: %v \n", len(students))

	// Start the server
	//server.Start()

	// Start gate server
	//go startGateServer(":8081")
	//
	//// Start login server
	//go startLoginServer(":8082")
	//// Block forever
	select {}
}

func runGinServer(config config.Config, store *sqlx.DB) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start("localhost:8080")
	if err != nil {
		log.Fatal("cannot start server")
	}
}

func studentByID(store *sqlx.DB, id uint32) (Student, error) {
	st := Student{}

	//if err := db.QueryRowx("SELECT * FROM students WHERE id = ?", id).StructScan(&st); err != nil {
	//	if err == sql.ErrNoRows {
	//		return st, fmt.Errorf("studentById %d: no such student", id)
	//	}
	//	return st, fmt.Errorf("studentById %d: %v", id, err)
	//}
	if err := store.Get(&st, "SELECT * FROM t_battle WHERE idx = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return st, fmt.Errorf("studentById %d: no such student", id)
		}
		return st, fmt.Errorf("studentById %d: %v", id, err)
	}
	return st, nil
}

func fetchStudents(store *sqlx.DB) ([]Student, error) {
	// A slice of Students to hold data from returned rows.
	students := make([]Student, 0, 10)

	err := store.Select(&students, "SELECT * FROM t_battle LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("fetchStudents %v", err)
	}

	return students, nil
}

type Student struct {
	Idx             int    `db:"idx"`
	AccountId       string `db:"account_id"`
	BattleData      []byte `db:"battle_data"`
	KillsModifytime int    `db:"kills_modifytime"`
	Createtime      int    `db:"createtime"`
	Modifytime      int    `db:"modifytime"`
}
