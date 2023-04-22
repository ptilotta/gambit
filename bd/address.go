package bd

import (
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ptilotta/gambit/models"
)

func InsertAddress(addr models.Address, User string) error {
	fmt.Println("Comienza el Registro InsertAddress")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "INSERT INTO addresses (Add_UserId, Add_Address, Add_City, Add_State, Add_PostalCode, Add_Phone, Add_Title, Add_Name )"
	sentencia += " VALUES ('" + User + "','" + addr.AddAddress + "','" + addr.AddCity + "','" + addr.AddState + "','"
	sentencia += addr.AddPostalCode + "','" + addr.AddPhone + "','" + addr.AddTitle + "','" + addr.AddName + "')"

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sentencia)
	fmt.Println("Insert Address > Ejecución Exitosa")
	return nil
}

func AddressExists(User string, id int) (error, bool) {
	fmt.Println("Comienza AddressExists")

	err := DbConnect()
	if err != nil {
		return err, false
	}
	defer Db.Close()

	sentencia := "SELECT 1 FROM addresses WHERE Add_Id = " + strconv.Itoa(id) + " AND Add_UserId = '" + User + "'"
	fmt.Println(sentencia)

	rows, err := Db.Query(sentencia)
	if err != nil {
		return err, false
	}

	var valor string
	rows.Next()
	rows.Scan(&valor)

	fmt.Println("AddressExists > Ejecución Exitosa - valor devuelto " + valor)
	if valor == "1" {
		return nil, true
	}
	return nil, false
}

func UpdateAddress(addr models.Address) error {
	fmt.Println("Comienza UpdateAddress")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "UPDATE addresses SET "

	if addr.AddAddress != "" {
		sentencia += "Add_Address = '" + addr.AddAddress + "', "
	}
	if addr.AddCity != "" {
		sentencia += "Add_City = '" + addr.AddCity + "', "
	}
	if addr.AddName != "" {
		sentencia += "Add_Name = '" + addr.AddName + "', "
	}
	if addr.AddPhone != "" {
		sentencia += "Add_Phone = '" + addr.AddPhone + "', "
	}
	if addr.AddPostalCode != "" {
		sentencia += "Add_PostalCode = '" + addr.AddPostalCode + "', "
	}
	if addr.AddState != "" {
		sentencia += "Add_State = '" + addr.AddState + "', "
	}
	if addr.AddTitle != "" {
		sentencia += "Add_Title = '" + addr.AddTitle + "', "
	}

	sentencia, _ = strings.CutSuffix(sentencia, ", ")
	sentencia += " WHERE Add_Id = " + strconv.Itoa(addr.AddId)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sentencia)
	fmt.Println("Update Address > Ejecución Exitosa")
	return nil

}

func DeleteAddress(id int) error {
	fmt.Println("Comienza DeleteAddress")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "DELETE FROM addresses WHERE Add_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sentencia)
	fmt.Println("Delete Address > Ejecución Exitosa")
	return nil
}
