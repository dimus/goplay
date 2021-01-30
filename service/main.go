package main

import "log"

func main() {
	// tidy up how main works with errors
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, dbtidy, err := setupDatabase()
	if err != nil {
		return error.Wrap(err, "setup database")
	}
	defer dbtidy()
	srv := &server{
		db: db,
	}
	return nil
}
