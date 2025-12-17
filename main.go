package main

import (
	"github.com/benjamin-larsen/NoctesChat-WebSocket/database"
	"github.com/benjamin-larsen/NoctesChat-WebSocket/ws"
	env "github.com/joho/godotenv"
)

func main() {
	env.Load()
	database.InitDB()
	defer database.DB.Close()

	ws.SetupWS()
}