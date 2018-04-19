package users

import (
	"fmt"
	"regexp"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// _ allows us to import the MYSQL driver without creating a local name
// for the package.
// This ensures the package gets into your built executable,
// but avoids the compile error you'd normally get from
// not calling any functions within that package.

func TestMySQLStore(t *testing.T) {
	//create a new sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sql mock: %v", err)
	}
	//ensure it's closed at the end of the test
	defer db.Close()

	newUser := CreateNewUser()

	expectedUser, _ := newUser.ToUser()
	//construct a new MySQLStore using the mock db
	store := NewMySQLStore(db)

	rows := sqlmock.NewRows([]string{"id", "email", "passhash", "username", "firstname", "lastname", "photourl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	// test insert function
	mock.ExpectExec(regexp.QuoteMeta(sqlInsertUser)).
		WithArgs(expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = store.Insert(expectedUser)
	if err != nil {
		t.Errorf("unexpected error occurs when inserting new user: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlInsertUser)).
		WithArgs(expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL).
		WillReturnError(fmt.Errorf("test DMBS error"))

	_, err = store.Insert(expectedUser)
	if err == nil {
		t.Errorf("expected error does not occurs when inserting new user: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	// test get function
	// test get by id
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUserByID)).
		WithArgs(expectedUser.ID).
		WillReturnRows(rows)

	_, err = store.GetByID(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error occurs when get user by ID: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUserByID)).
		WithArgs(expectedUser.ID).
		WillReturnError(fmt.Errorf("test DMBS error"))

	_, err = store.GetByID(expectedUser.ID)
	if err == nil {
		t.Errorf("expected does not error occurs when getting user by id: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	// test get by email
	rows = sqlmock.NewRows([]string{"id", "email", "passhash", "username", "firstname", "lastname", "photourl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUserByEmail)).
		WithArgs(expectedUser.Email).
		WillReturnRows(rows)

	_, err = store.GetByEmail(expectedUser.Email)
	if err != nil {
		t.Errorf("unexpected error occurs when get user by email: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUserByEmail)).
		WithArgs(expectedUser.Email).
		WillReturnError(fmt.Errorf("test DMBS error"))

	_, err = store.GetByEmail(expectedUser.Email)
	if err == nil {
		t.Errorf("expected does not error occurs when getting user by email: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	// test get by user name
	rows = sqlmock.NewRows([]string{"id", "email", "passhash", "username", "firstname", "lastname", "photourl"})
	rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, expectedUser.FirstName, expectedUser.LastName, expectedUser.PhotoURL)

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUserByUserName)).
		WithArgs(expectedUser.UserName).
		WillReturnRows(rows)

	_, err = store.GetByUserName(expectedUser.UserName)
	if err != nil {
		t.Errorf("unexpected error occurs when get user by UserName: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectUserByUserName)).
		WithArgs(expectedUser.UserName).
		WillReturnError(fmt.Errorf("test DMBS error"))

	_, err = store.GetByUserName(expectedUser.UserName)
	if err == nil {
		t.Errorf("expected does not error occurs when getting user by UserName: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	// Test update
	update := &Updates{
		FirstName: "updatedFirstName",
		LastName:  "updatedLastName",
	}
	// updatedRows := sqlmock.NewRows([]string{"id", "email", "passhash", "UserName", "FirstName", "LastName", "photourl"})
	// rows.AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PassHash, expectedUser.UserName, update.FirstName, update.LastName, expectedUser.PhotoURL)

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(update.FirstName, update.LastName, expectedUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.Update(expectedUser.ID, update)

	if err != nil {
		t.Errorf("unexpected error occurs when update user: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(update.FirstName, update.LastName, expectedUser.ID).
		WillReturnError(fmt.Errorf("test DMBS error"))

	err = store.Update(expectedUser.ID, update)
	if err == nil {
		t.Errorf("expected does not error occurs when update user: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	//Test delete

	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(expectedUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.Delete(expectedUser.ID)
	if err != nil {
		t.Errorf("unexpected error occurs when deleting user: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(sqlDelete)).
		WithArgs(expectedUser.ID).
		WillReturnError(fmt.Errorf("test DMBS error"))

	err = store.Delete(expectedUser.ID)
	if err == nil {
		t.Errorf("expected error does not occurs when deleting user: %v", err)
	}

	//ensure we didn't have any unmet expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}
