package users

import (
	"database/sql"
	"log"

	database "github.com/Muchogoc/hackernews/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

//User represents user from database
type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

//Create inserts user object into database
func (user *User) Create() {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username,Password) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	hashedPassword, err := HashPassword(user.Password)

	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

//GetUserIDByUsername check if a user exists in database by given username
func GetUserIDByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}

	row := statement.QueryRow(username)

	var ID int

	err = row.Scan(&ID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}

		return 0, err
	}

	return ID, nil
}

//Authenticate check if user exists in database and generate token
func (user *User) Authenticate() bool {
	statement, err := database.Db.Prepare("select Password from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}

	row := statement.QueryRow(user.Username)

	var hashedPassword string

	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		log.Fatal(err)
	}

	return CheckPasswordHash(user.Password, hashedPassword)
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

//CheckPasswordHash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
