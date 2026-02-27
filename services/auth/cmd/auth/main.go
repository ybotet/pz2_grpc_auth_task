package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/ybotet/pz2_grpc_auth_task/services/auth/internal/auth"
	grpcserver "github.com/ybotet/pz2_grpc_auth_task/services/auth/internal/grpc"
)

func main() {
    port := os.Getenv("AUTH_GRPC_PORT")
    if port == "" {
        port = "50051"
        log.Printf("AUTH_GRPC_PORT no configurado, usando puerto por defecto: %s", port)
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "your-secret-key-change-in-production"
        log.Printf("JWT_SECRET no configurado, usando clave por defecto (NO USAR EN PRODUCCIÓN)")
    }

    lis, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("Error al escuchar en puerto %s: %v", port, err)
    }

    s := grpc.NewServer()
    authService := auth.NewService(jwtSecret)
    grpcserver.Register(s, authService)

    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        log.Println("Apagando servidor gRPC...")
        s.GracefulStop()
    }()

    log.Printf("Servidor Auth gRPC escuchando en puerto %s", port)
    
    // Generar token de ejemplo para pruebas
    testToken, _ := authService.GenerateToken("usuario_prueba", 24*time.Hour)
    log.Printf("Token de prueba (válido 24h): %s", testToken)
    
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Error al servir: %v", err)
    }
}