package main

import (
	"cloudace/progcon/test/gae/handler/router"
	"net/http"

	"google.golang.org/appengine"
)

// func init() {
// 	m := router.NewRouter()
// 	http.Handle("/", m)
// }

func main() {
	m := router.NewRouter()
	http.Handle("/", m)
	appengine.Main()
}
