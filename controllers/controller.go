package main

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, request *http.Request)  {
	fmt.Fprint(w, "<h1>Hello World!</h1>")
}
