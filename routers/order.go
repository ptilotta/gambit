package routers

import (
	"encoding/json"
	"strconv"

	//	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/gambit/bd"
	"github.com/ptilotta/gambit/models"
)

func InsertOrder(body string, User string) (int, string) {
	var o models.Orders
	err := json.Unmarshal([]byte(body), &o)
	if err != nil {
		return 400, "Error en los datos recibidos " + err.Error()
	}

	o.Order_UserUUID = User

	OK, message := ValidOrder(o)
	if !OK {
		return 400, message
	}

	result, err2 := bd.InsertOrder(o)
	if err2 != nil {
		return 400, "Ocurrió un error al intentar realizar el registro de la orden " + err2.Error()
	}

	return 200, "{ OrderID :" + strconv.Itoa(int(result)) + "}"

}

func ValidOrder(o models.Orders) (bool, string) {
	if o.Order_Total == 0 {
		return false, "Debe indicar el total de la orden"
	}

	count := 0
	for _, od := range o.OrderDetails {
		if od.OD_ProdId == 0 {
			return false, "Debe indicar el ID del producto en el detalle de la Orden"
		}
		if od.OD_Quantity == 0 {
			return false, "Debe indicar la cantidad del producto en el detalle de la Orden"
		}
		count++
	}
	if count == 0 {
		return false, "Debe indicar items en la orden"
	}
	return true, ""
}
