package users

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// NewUser creates a NewUser with all valid fields.
func CreateNewUser() *NewUser {
	return &NewUser{
		Email:        "test@gmail.com",
		Password:     "password",
		PasswordConf: "password",
		UserName:     "testUserName",
		FirstName:    "firstname",
		LastName:     "lastname",
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name              string
		invalidField      string
		invalidFieldValue string
		hint              string
	}{
		{
			"valid new user",
			"",
			"",
			"this is a valid NewUser",
		},
		{
			"invalid email",
			"Email",
			"invalid",
			"error parsing email",
		},
		{
			"invalid email",
			"Email",
			"invalid@",
			"error parsing email",
		},
		{
			"invalid email",
			"Email",
			"@invalid",
			"error parsing email",
		},
		{
			"invalid password",
			"Password",
			"12345",
			"Password needs to be at least 6 characters",
		},
		{
			"invalid password confirmation",
			"PasswordConf",
			"wordpass",
			"password needs to match password confirmation",
		},
		{
			"empty username",
			"UserName",
			"",
			"username must be non-zero length",
		}, {
			"username with space",
			"UserName",
			"ab  cd",
			"username may not contain spaces",
		},
	}

	for _, each := range cases {
		nu := CreateNewUser()

		validate := reflect.ValueOf(nu).Elem().FieldByName(each.invalidField)

		if validate.IsValid() {
			validate.SetString(each.invalidFieldValue)
		}

		err := nu.Validate()

		// Test valid cases.
		if each.invalidField == "" && each.invalidFieldValue == "" {
			if err != nil {
				t.Errorf("case: %sinvalid field: {%s: %s} hint: %s", each.name, each.invalidField, each.invalidFieldValue, each.hint)
			}
		} else {
			// Test invalid cases.
			if err == nil {
				t.Errorf("case: %sinvalid field: {%s: %s} hint: %s", each.name, each.invalidField, each.invalidFieldValue, each.hint)
			}
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		name         string
		validEmail   string
		invalidEmail string
		hint         string
	}{
		{
			"email contains leading or trailing whitespace",
			"test@gmail.com",
			" test@gmail.com   ",
			"make sure to trim leading and trailing whitespace from an email address",
		},
		{
			"email contains uppercase characters",
			"apple@uw.edu",
			"apple@UW.edu",
			"Should convert all characters to lower-case",
		},
	}

	for _, c := range cases {
		nu := CreateNewUser()

		email := c.invalidEmail
		// Trim leading and trailing whitespace from an email address.
		email = strings.TrimSpace(email)

		// Force all characters in the email to be lower-case.
		email = strings.ToLower(email)

		nu.Email = email

		// Convert NewUser to User.
		usr, err := nu.ToUser()
		if err != nil {
			t.Errorf("error converting NewUser to User\n")
		}

		if usr == nil {
			t.Errorf("ToUser() returned nil\n")
		}

		// Test Email.
		if usr.Email != c.validEmail {
			t.Errorf("case: %sgot: %s should be : %s hint: %s", c.name, usr.Email, c.validEmail, c.hint)
		}

		// Test PhotoURL.
		h := md5.New()
		io.WriteString(h, nu.Email)
		src := h.Sum(nil)
		result := hex.EncodeToString(src)
		photoURL := gravatarBasePhotoURL + result

		if len(usr.PhotoURL) == 0 {
			t.Errorf("PhotoURL field is empty\n")
		}

		if usr.PhotoURL != photoURL {
			t.Errorf("invalid PhotoURL")
		}

		// Test PassHash.
		if len(usr.PassHash) == 0 {
			t.Errorf("password hash is empty")
		}

		err = bcrypt.CompareHashAndPassword(usr.PassHash, []byte(nu.Password))
		if err != nil {
			t.Errorf("invalid password: %v", err)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name           string
		firstName      string
		lastName       string
		expectedOutput string
	}{
		{
			"both first last name are non-empty",
			"Apple",
			"Juice",
			"Apple Juice",
		},
		{
			"first name is empty",
			"",
			"Juice",
			"Juice",
		},
		{
			"last name is empty",
			"Apple",
			"",
			"Apple",
		},
		{
			"both first last name are empty",
			"",
			"",
			"",
		},
	}

	for _, each := range cases {
		nu := CreateNewUser()
		nu.FirstName = each.firstName
		nu.LastName = each.lastName

		usr, err := nu.ToUser()
		if err != nil {
			t.Errorf("error converting NewUser to User\n")
		}

		fullName := usr.FullName()
		if fullName != each.expectedOutput {
			t.Errorf("case: %sgot: %s should be : %s", each.name, fullName, each.expectedOutput)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	usr := &User{}

	if err := usr.SetPassword("password"); err != nil {
		t.Errorf("error when set the password: %v", err)
	}

	if err := usr.Authenticate("password"); err != nil {
		t.Errorf("the password is valid")
	}

	if err := usr.Authenticate(""); err == nil {
		t.Errorf("empth password")
	}

	if err := usr.Authenticate("ppppp"); err == nil {
		t.Errorf("invalid password")
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name    string
		updates *Updates
		valid   bool
	}{
		{
			"valid updates",
			&Updates{
				"Apple",
				"Juice",
			},
			false,
		},
		{
			"first name is empty",
			&Updates{
				"",
				"Juice",
			},
			true,
		},
		{
			"last name is empty",
			&Updates{
				"Apple",
				"",
			},
			true,
		},
		{
			"both first name and last name are empty",
			&Updates{
				"",
				"",
			},
			true,
		},
	}

	for _, each := range cases {
		usr := &User{}

		err := usr.ApplyUpdates(each.updates)
		if (!each.valid && err != nil) || (each.valid && err == nil) {
			t.Errorf("case: %sexpect error: %v actual error: %s", each.name, each.valid, err)
		}
	}
}
