package bd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ptilotta/gambit/models"
	"github.com/ptilotta/gambit/tools"
)

func InsertProduct(p models.Product) (int64, error) {
	fmt.Println("Comienza Registro")

	err := DbConnect()
	if err != nil {
		return 0, err
	}
	defer Db.Close()

	sentencia := "INSERT INTO products (Prod_Title "

	if len(p.ProdDescription) > 0 {
		sentencia += ", Prod_Description"
	}
	if p.ProdPrice > 0 {
		sentencia += ", Prod_Price"
	}
	if p.ProdCategId > 0 {
		sentencia += ", Prod_CategoryId"
	}
	if p.ProdStock > 0 {
		sentencia += ", Prod_Stock"
	}
	if len(p.ProdPath) > 0 {
		sentencia += ", Prod_Path"
	}

	sentencia += ") VALUES ('" + tools.EscapeString(p.ProdTitle) + "'"

	if len(p.ProdDescription) > 0 {
		sentencia += ",'" + tools.EscapeString(p.ProdDescription) + "'"
	}
	if p.ProdPrice > 0 {
		sentencia += ", " + strconv.FormatFloat(p.ProdPrice, 'e', -1, 64)
	}
	if p.ProdCategId > 0 {
		sentencia += ", " + strconv.Itoa(p.ProdCategId)
	}
	if p.ProdStock > 0 {
		sentencia += ", " + strconv.Itoa(p.ProdStock)
	}
	if len(p.ProdPath) > 0 {
		sentencia += ", '" + tools.EscapeString(p.ProdPath) + "'"
	}

	sentencia += ")"

	var result sql.Result
	result, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	LastInsertId, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	fmt.Println("Insert Product > Ejecución Exitosa")
	return LastInsertId, nil
}

func UpdateProduct(p models.Product) error {
	fmt.Println("Comienza Update")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "Update products SET "

	sentencia = tools.ArmoSentencia(sentencia, "Prod_Title", "S", 0, 0, p.ProdTitle)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Description", "S", 0, 0, p.ProdDescription)
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Price", "F", 0, p.ProdPrice, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_CategoryId", "N", p.ProdCategId, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Stock", "N", p.ProdStock, 0, "")
	sentencia = tools.ArmoSentencia(sentencia, "Prod_Path", "S", 0, 0, p.ProdPath)

	sentencia += " WHERE Prod_Id = " + strconv.Itoa(p.ProdId)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Update Product > Ejecución Exitosa")
	return nil
}

func DeleteProduct(id int) error {
	fmt.Println("Comienza Delete Product")

	err := DbConnect()
	if err != nil {
		return err
	}
	defer Db.Close()

	sentencia := "DELETE FROM products WHERE Prod_Id = " + strconv.Itoa(id)

	_, err = Db.Exec(sentencia)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Delete Product > Ejecución Exitosa")
	return nil
}

func SelectProduct(p models.Product, choice string, page int, pageSize int, orderType string, orderField string) (models.ProductResp, error) {
	fmt.Println("Comienza SelectProduct")
	var Resp models.ProductResp
	var Prod []models.Product

	err := DbConnect()
	if err != nil {
		return Resp, err
	}
	defer Db.Close()

	var sentencia string
	var sentenciaCount string
	var where, limit, join string

	sentencia = "SELECT Prod_Id, Prod_Title, Prod_Description, Prod_CreatedAt, Prod_Updated, Prod_Price, Prod_Path, Prod_CategoryId, Prod_Stock FROM products "
	sentenciaCount = "SELECT count(*) as registros FROM products "

	switch choice {
	case "P":
		where = " WHERE Prod_Id = " + strconv.Itoa(p.ProdId)
	case "S":
		where = " WHERE UCASE(CONCAT(Prod_Title, Prod_Description)) LIKE '%" + strings.ToUpper(p.ProdSearch) + "%' "
	case "C":
		where = " WHERE Prod_CategoryId = " + strconv.Itoa(p.ProdCategId)
	case "U":
		where = " WHERE UCASE(Prod_Path) LIKE '%" + strings.ToUpper(p.ProdPath) + "%' "
	case "K":
		join = " JOIN category ON Prod_CategoryId = Categ_Id AND Categ_Path LIKE '%" + strings.ToUpper(p.ProdCategPath) + "%' "
		sentencia += join
		sentenciaCount += join
	}

	sentenciaCount += where

	var rows *sql.Rows
	rows, err = Db.Query(sentenciaCount)
	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
		return Resp, err
	}

	rows.Next()
	var regi sql.NullInt32
	err = rows.Scan(&regi)

	registros := int(regi.Int32)

	if page > 0 {
		if registros > pageSize {
			limit = " LIMIT " + strconv.Itoa(pageSize)
			if page > 1 {
				offset := pageSize * (page - 1)
				limit += " OFFSET " + strconv.Itoa(offset)
			}
		} else {
			limit = ""
		}
	}

	var orderBy string
	if len(orderField) > 0 {
		switch orderField {
		case "I":
			orderBy = " ORDER BY Prod_Id "
		case "T":
			orderBy = " ORDER BY Prod_Title "
		case "D":
			orderBy = " ORDER BY Prod_Description "
		case "F":
			orderBy = " ORDER BY Prod_CreatedAt "
		case "P":
			orderBy = " ORDER BY Prod_Price "
		case "S":
			orderBy = " ORDER BY Prod_Stock "
		case "C":
			orderBy = " ORDER BY Prod_CategoryId "
		}
		if orderType == "D" {
			orderBy += " DESC"
		}
	}

	sentencia += where + orderBy + limit

	fmt.Println(sentencia)

	rows, err = Db.Query(sentencia)

	for rows.Next() {
		var p models.Product
		var ProdId sql.NullInt32
		var ProdTitle sql.NullString
		var ProdDescription sql.NullString
		var ProdCreatedAt sql.NullTime
		var ProdUpdated sql.NullTime
		var ProdPrice sql.NullFloat64
		var ProdPath sql.NullString
		var ProdCategoryId sql.NullInt32
		var ProdStock sql.NullInt32

		err := rows.Scan(&ProdId, &ProdTitle, &ProdDescription, &ProdCreatedAt, &ProdUpdated, &ProdPrice, &ProdPath, &ProdCategoryId, &ProdStock)
		if err != nil {
			return Resp, err
		}

		p.ProdId = int(ProdId.Int32)
		p.ProdTitle = ProdTitle.String
		p.ProdDescription = ProdDescription.String
		p.ProdCreatedAt = ProdCreatedAt.Time.String()
		p.ProdUpdated = ProdUpdated.Time.String()
		p.ProdPrice = ProdPrice.Float64
		p.ProdPath = ProdPath.String
		p.ProdCategId = int(ProdCategoryId.Int32)
		p.ProdStock = int(ProdStock.Int32)
		Prod = append(Prod, p)
	}

	Resp.TotalItems = registros
	Resp.Data = Prod

	fmt.Println("Select Product > Ejecución Exitosa")
	return Resp, nil
}
