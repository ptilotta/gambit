package handlers

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ptilotta/gambit/auth"
	"github.com/ptilotta/gambit/routers"
)

func Manejadores(path string, method string, body string, headers map[string]string, request events.APIGatewayV2HTTPRequest) (int, string) {

	fmt.Println("Voy a procesar " + path + " > " + method)

	id := request.PathParameters["id"]
	idn, _ := strconv.Atoi(id)

	isOk, statusCode, user := validoAuthorization(path, method, headers)
	if !isOk {
		return statusCode, user
	}

	fmt.Println("path[0:4] = " + path[0:4])

	switch path[0:4] {
	case "user":
		return ProcesoUsers(body, path, method, user, id, request)
	case "prod":
		return ProcesoProducts(body, path, method, user, idn, request)
	case "stoc":
		return ProcesoStock(body, path, method, user, idn, request)
	case "addr":
		return ProcesoAddress(body, path, method, user, idn, request)
	case "cate":
		return ProcesoCategory(body, path, method, user, idn, request)
	case "orde":
		return ProcesoOrder(body, path, method, user, idn, request)
	}

	return 400, "Method Invalid"
}

func validoAuthorization(path string, method string, headers map[string]string) (bool, int, string) {
	if (path == "product" && method == "GET") ||
		(path == "category" && method == "GET") {
		return true, 200, ""
	}

	token := headers["authorization"]
	if len(token) == 0 {
		return false, 401, "Token requerido"
	}

	todoOK, err, msg := auth.ValidoToken(token)
	if !todoOK {
		if err != nil {
			fmt.Println("Error en el token " + err.Error())
			return false, 401, err.Error()
		} else {
			fmt.Println("Error en el token " + msg)
			return false, 401, msg
		}
	}

	fmt.Println("Token OK")
	return true, 200, msg
}

func ProcesoUsers(body string, path string, method string, user string, id string, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method Invalid"
}

func ProcesoProducts(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	switch method {
	case "POST":
		return routers.InsertProduct(body, user)
	case "PUT":
		return routers.UpdateProduct(body, user, id)
	case "DELETE":
		return routers.DeleteProduct(user, id)
	case "GET":
		return routers.SelectProduct(request)
	}

	return 400, "Method Invalid"
}

func ProcesoCategory(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	switch method {
	case "POST":
		return routers.InsertCategory(body, user)
	case "PUT":
		return routers.UpdateCategory(body, user, id)
	case "DELETE":
		return routers.DeleteCategory(body, user, id)
	case "GET":
		return routers.SelectCategories(body, request)
	}
	return 400, "Method Invalid"
}

func ProcesoStock(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method Invalid"
}

func ProcesoAddress(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return routers.UpdateStock(body, user, id)
}

func ProcesoOrder(body string, path string, method string, user string, id int, request events.APIGatewayV2HTTPRequest) (int, string) {
	return 400, "Method Invalid"
}
