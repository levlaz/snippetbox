package examples

import (
	"database/sql"
)

// We need somewhere to store the prepared statement for the lifetime of our // web application. A neat way is to embed it in the model alongside the
// connection pool.
type ExampleModel struct {
	DB         *sql.DB
	InsertStmt *sql.Stmt
}

// Create a constructor for the model, in which we set up the prepared // statement.
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
	// Use the Prepare method to create a new prepared statement for the
	// current connection pool. This returns a sql.Stmt object which represents // the prepared statement.
	insertStmt, err := db.Prepare("INSERT INTO ...")
	if err != nil {
		return nil, err
	}
	// Store it in our ExampleModel struct, alongside the connection pool.
	return &ExampleModel{DB: db, InsertStmt: insertStmt}, nil
}

// Any methods implemented against the ExampleModel struct will have access to // the prepared statement.
func (m *ExampleModel) Insert(args string) error {
	// We then need to call Exec directly against the prepared statement, rather // than against the connection pool. Prepared statements also support the
	// Query and QueryRow methods.
	_, err := m.InsertStmt.Exec(args)
	return err
}

// // In the web application's main function we will need to initialize a new // ExampleModel struct using the constructor function.
// func main() {
// 	db, err := sql.Open("", "")
// 	if err != nil {
// 		logger.Error(err.Error())
// 		os.Exit(1)
// 	}
// 	defer db.Close()
// 	// Use the constructor function to create a new ExampleModel struct.
// 	exampleModel, err := NewExampleModel(db)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		os.Exit(1)
// 	}
// 	// Defer a call to Close() on the prepared statement to ensure that it is // properly closed before our main function terminates.
// 	defer exampleModel.InsertStmt.Close()
// }
