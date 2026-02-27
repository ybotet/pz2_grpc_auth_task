package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/ybotet/pz2_grpc_auth_task/services/tasks/internal/clients"
	"github.com/ybotet/pz2_grpc_auth_task/services/tasks/internal/handlers"
	"github.com/ybotet/pz2_grpc_auth_task/services/tasks/internal/middleware"
)

func main() {
    tasksPort := os.Getenv("TASKS_PORT")
    if tasksPort == "" {
        tasksPort = "8082"
    }

    authAddr := os.Getenv("AUTH_GRPC_ADDR")
    if authAddr == "" {
        authAddr = "localhost:50051"
    }

    // Conectar a Auth service
    authClient, err := clients.NewAuthClient(authAddr)
    if err != nil {
        log.Fatalf("Error conectando a Auth service: %v", err)
    }
    defer authClient.Close()

    // Crear middleware y handlers
    authMiddleware := middleware.NewAuthMiddleware(authClient.GetClient())
    taskHandler := handlers.NewTaskHandler()

    // Configurar rutas
    r := mux.NewRouter()
    r.HandleFunc("/tasks", authMiddleware.Authenticate(taskHandler.GetTasks)).Methods("GET")

    log.Printf("Servidor Tasks escuchando en puerto %s", tasksPort)
    log.Fatal(http.ListenAndServe(":"+tasksPort, r))
}
