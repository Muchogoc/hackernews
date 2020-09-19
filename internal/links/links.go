package links

import (
	"log"

	database "github.com/Muchogoc/hackernews/internal/pkg/db/mysql"
	"github.com/Muchogoc/hackernews/internal/users"
)

//Link represent users we get from database
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

//Save inserts link object into database and returns It's ID
func (link Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	log.Print("Row inserted!")

	return id
}

//GetAll queries database and returns all links
func GetAll() []Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID") // changed
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var links []Link
	var username string
	var id string

	for rows.Next() {
		var link Link

		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username) // changed
		if err != nil {
			log.Fatal(err)
		}

		link.User = &users.User{
			ID:       id,
			Username: username,
		}

		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return links
}
