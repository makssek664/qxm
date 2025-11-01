package main;

import
(
	"log"
)

func main() {
	h := httpInitServer(nil)
	err := h.httpListen()	
	if err != nil {
		log.Fatal(err)
	}
}
