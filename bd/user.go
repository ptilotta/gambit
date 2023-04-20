package bd

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ptilotta/gambit/models"
	"github.com/ptilotta/gambit/tools"
)

func UpdateUser(UField models.User, User string) error {
	fmt.Println("Comienza UpdateUser")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE users SET "

	coma := ""
	if len(UField.UserFirstName) > 0 {
		coma = ","
		sentencia += "User_FirstName = '" + UField.UserFirstName + "'"
	}

	if len(UField.UserLastName) > 0 {
		sentencia += coma + "User_LastName = '" + UField.UserLastName + "'"
	}

	sentencia += ", User_DateUpg = '" + tools.FechaMySQL() + "' WHERE User_UUID='" + User + "'"

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update User > Ejecución Exitosa")
	return nil
}

func SelectUser(UserId string) (models.User, error) {
	fmt.Println("Comienza SelectUser")
	User := models.User{}

	err := DbConnect()
	if err != nil {
		return User, err
	}
	defer Db.Close()

	sentencia := "SELECT * FROM users WHERE User_UUID = '" + UserId + "'"

	var rows *sql.Rows
	rows, err = Db.Query(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return User, err
	}
	defer rows.Close()

	rows.Next()

	var firstName sql.NullString
	var lastName sql.NullString
	var dateUpg sql.NullTime

	rows.Scan(&User.UserUUID, &User.UserEmail, &firstName, &lastName, &User.UserStatus, &User.UserDateAdd, &dateUpg)

	User.UserFirstName = firstName.String
	User.UserLastName = lastName.String
	User.UserDateUpd = dateUpg.Time.String()

	fmt.Println("Select User > Ejecución Exitosa")
	return User, nil

}

func SelectUsers(Page int) (models.ListUsers, error) {
	fmt.Println("Comienza SelectUsers")

	var lu models.ListUsers
	User := []models.User{}

	err := DbConnect()
	if err != nil {
		return lu, err
	}
	defer Db.Close()

	var offset int = (Page * 10) - 10
	var sentencia string
	var sentenciaCount string = "SELECT count(*) as registros FROM users"

	sentencia = "select * from users LIMIT 10"
	if offset > 0 {
		sentencia += " OFFSET " + strconv.Itoa(offset)
	}

	var rowsCount *sql.Rows

	rowsCount, err = Db.Query(sentenciaCount)
	if err != nil {
		return lu, err
	}

	defer rowsCount.Close()

	rowsCount.Next()

	var registros int
	rowsCount.Scan(&registros)
	lu.TotalItems = registros

	var rows *sql.Rows
	rows, err = Db.Query(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return lu, err
	}

	for rows.Next() {
		var u models.User
		var firstName sql.NullString
		var lastName sql.NullString
		var dateUpg sql.NullTime

		rows.Scan(&u.UserUUID, &u.UserEmail, &firstName, &lastName, &u.UserStatus, &u.UserDateAdd, &dateUpg)

		u.UserFirstName = firstName.String
		u.UserLastName = lastName.String
		u.UserDateUpd = dateUpg.Time.String()
		User = append(User, u)
	}

	fmt.Println("Select Users > Ejecución Exitosa")

	lu.Data = User
	return lu, nil
}
