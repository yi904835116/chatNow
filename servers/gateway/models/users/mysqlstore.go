package users

import (
	"database/sql"
	"fmt"

	"github.com/info344-s18/challenges-yi904835116/servers/gateway/indexes"
)

// Various SQL statements we will need to execute.

// SQL to select a particular user by ID.
// Use `?` for column values that we will get at runtime.
const sqlSelectUserByID = "select * from user where id=?"

const sqlSelectAll = "select * from user"

// SQL to select a particular user by email.
const sqlSelectUserByEmail = "select * from user where email=?"

// SQL to select a particular user by username.
const sqlSelectUserByUserName = "select * from user where username=?"

// SQL to insert a new user row.
const sqlInsertUser = "insert into user(email,passhash,username,firstname,lastname,photourl) values (?,?,?,?,?,?)"

// SQL to update user.
const sqlUpdate = "update user set firstname=?, lastname=? where id=?"

// SQL to delete user.
const sqlDelete = "delete from user where id=?"

type userRow struct {
	id        int64
	email     string
	passhash  []byte
	username  string
	firstname string
	lastname  string
	photourl  string
}

// MySQLStore implements Store for a MySQL database.
type MySQLStore struct {
	// a live reference to the database.
	db *sql.DB
}

// NewMySQLStore constructs a MySQLStore.
func NewMySQLStore(db *sql.DB) *MySQLStore {
	if db == nil {
		panic("nil pointer passed to NewMySQLStore")
	}

	return &MySQLStore{
		db: db,
	}
}

// GetByID returns the User with the given ID.
func (store *MySQLStore) GetByID(id int64) (*User, error) {
	rows, err := store.db.Query(sqlSelectUserByID, id)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %v", err)
	}

	users, err := handleResult(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %s", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found")
	}

	// Return the first (and only) element from the slice.
	return users[0], nil
}

// GetByEmail returns the User with the given email.
func (store *MySQLStore) GetByEmail(email string) (*User, error) {
	rows, err := store.db.Query(sqlSelectUserByEmail, email)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %v", err)
	}

	users, err := handleResult(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %s", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found")
	}

	// Return the first (and only) element from the slice.
	return users[0], nil
}

// GetByUserName returns the User with the given Username.
func (store *MySQLStore) GetByUserName(username string) (*User, error) {
	rows, err := store.db.Query(sqlSelectUserByUserName, username)
	if err != nil {
		return nil, fmt.Errorf("error selecting user: %v", err)
	}

	users, err := handleResult(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %s", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user found")
	}

	// Return the first (and only) element from the slice.
	return users[0], nil
}

// Insert converts the NewUser to a User, inserts
// it into the database, and returns it.
func (store *MySQLStore) Insert(user *User) (*User, error) {

	// Use transaction to make sure inserts to be atomic (all or nothing).
	database := store.db

	res, err := database.Exec(sqlInsertUser, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)

	if err != nil {
		fmt.Errorf("error inserting new row: %v", err)
	} else {
		//get the auto-assigned ID for the new row
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Errorf("error getting new ID: %v", id)
		} else {
			user.ID = id
		}
	}

	if err != nil {
		// Rollback the transaction if there's an error.
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	return user, nil
}

// Update applies UserUpdates to the given user ID.
func (store *MySQLStore) Update(userID int64, updates *Updates) (*User, error) {
	if updates == nil {
		return nil, fmt.Errorf("Updates is nil")
	}

	_, err := store.db.Exec(sqlUpdate, updates.FirstName, updates.LastName, userID)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}

	return store.GetByID(userID)
}

// Delete deletes the user with the given ID.
func (store *MySQLStore) Delete(userID int64) error {

	_, err := store.db.Exec(sqlDelete, userID)
	if err != nil {
		return fmt.Errorf("error deleting data: %v", err)
	}

	return nil
}

// Trie returns a trie tree thtat stores existing users info
func (store *MySQLStore) Trie() (*indexes.Trie, error) {
	trie := indexes.NewTrie()
	rows, err := store.db.Query(sqlSelectAll)

	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("error selecting user: %v", err)
	}

	users, err := handleResult(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %s", err)
	}

	if len(users) == 0 {
		return trie, nil
	}

	for _, user := range users {
		trie.Insert(user.UserName, user.ID)
		trie.Insert(user.LastName, user.ID)
		trie.Insert(user.FirstName, user.ID)
	}

	return trie, nil
}

// ConvertIDToUsers converts all keys(User IDs) in a given map to a slice of User.
func (store *MySQLStore) ConvertIDToUsers(userIDs map[int64]bool) ([]*User, error) {
	users := []*User{}
	for userID := range userIDs {
		user, err := store.GetByID(userID)
		if err != nil {
			return nil, fmt.Errorf("error getting user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// scanUsers scans query result rows into a []*User.
func handleResult(rows *sql.Rows) ([]*User, error) {
	// Ensure the rows are closed regardless of how
	// we exit this function.
	defer rows.Close()

	users := []*User{}

	row := userRow{}

	for rows.Next() {

		// Scan each record into User struct.
		err := rows.Scan(&row.id, &row.email, &row.passhash, &row.username, &row.firstname, &row.lastname, &row.photourl)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		user := &User{
			ID:        row.id,
			Email:     row.email,
			PassHash:  row.passhash,
			UserName:  row.username,
			FirstName: row.firstname,
			LastName:  row.lastname,
			PhotoURL:  row.photourl,
		}

		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	return users, nil
}
