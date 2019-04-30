package f3api

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func RunServer(impl RestApi) {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/payments", impl.GetAllPayments),
		rest.Post("/payments", impl.PostPayment),
		rest.Get("/payments/:id", impl.GetPayment),
		rest.Put("/payments/:id", impl.PutPayment),
		rest.Delete("/payments/:id", impl.DeletePayment),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
