package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Проверяем аргументы командной строки
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go api [flags]     - Start the API server")
		fmt.Println("  go run main.go migrate <cmd>   - Run database migrations")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run main.go api")
		fmt.Println("  go run main.go api -port=3000 -db=./custom.db")
		fmt.Println("  go run main.go migrate up")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "api":
		// Запуск API сервера
		cmdArgs := append([]string{"run", "cmd/api/main.go"}, args...)
		cmd := exec.Command("go", cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}

	case "migrate":
		// Запуск миграций
		if len(args) == 0 {
			log.Fatal("Migration command required: up, down, status, reset")
		}
		cmdArgs := append([]string{"run", "cmd/migrate/main.go"}, args...)
		cmd := exec.Command("go", cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to run migration: %v", err)
		}

	default:
		log.Fatalf("Unknown command: %s\nAvailable commands: api, migrate", command)
	}
}
