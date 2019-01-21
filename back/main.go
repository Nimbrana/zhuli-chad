package main

import (
	"fmt"
	"net/http"

	. "./signaling"
)

func main() {

	http.HandleFunc("/", HandleConnections)

	fmt.Println("Server listen to port: 8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}

}
