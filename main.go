package main

import "ledger/service"

func main() {
	r := service.SetupRouter()
	r.Run(":8282")
}