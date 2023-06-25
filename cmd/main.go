package main

import (
	"flag"
	"fmt"
	_ "gRPC2/docs"
	"gRPC2/router"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// @title           gRPC service
// @version         1.0
// @description     This is the main server handling the grpc major operations.
// @termsOfService http://swagger.io/terms/

// @contact.name   gRPC
// @contact.url    https://gRPC.service.com/

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @host localhost:8080
// @BasePath /
// @schemes http

func main() {

	//utility.ConvertJSONToYML()

	r := mux.NewRouter()
	router.UserRouter(r)
	router.MovieRouter(r)
	router.ImageRouter(r)

	// Serve Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL path to the generated Swagger JSON file
	))

	// set swagger ui router
	//openAPIHandler := openapiv2.NewHandler()
	//r.PathPrefix("/q/").Handler(openAPIHandler)

	// run swagger ui on browser by this url
	//Run swagger ui on browser by this url : http://localhost:8080/q/swagger-ui/
	fmt.Println("Run swagger ui on browser by this url : http://localhost:8080/swagger/index.html")

	//fmt.Println("http server client Running port:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
