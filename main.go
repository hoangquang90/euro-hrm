package main

import (
	"context"
	"europm/internal/config"
	"europm/internal/db/dbhrm"
	"europm/internal/logging"
	"europm/internal/server"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

//	@title		Swagger FDM API
//	@version	2.0

//	@securityDefinitions.apikey	X-USER-ID
//	@in							header
//	@name						X-USER-ID

// @securityDefinitions.apikey	X-USER-NAME
// @in							header
// @name						X-USER-NAME
func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	fmt.Println("starting ...")
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	err := config.Init()
	if err != nil { // Handle errors reading the config file
		log.Fatalf("can't initialize config: %v", err)
	}

	err = logging.Init()
	if err != nil {
		log.Fatalf("can't initialize logger: %v", err)
	}
	//DB =======================
	err = dbhrm.Init()
	if err != nil {
		log.Fatalf("can't initialize db %v", err)
	}
	//==========================
	// HTTP Server ======================================
	err = server.Start()
	if err != nil {
		log.Fatalf("can't start http server %v", err)
	}

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	server.Stop()

	logging.Destroy()

	fmt.Println("stopped!")
	// Block main thread
	// select {}
}

// func main() {
// 	// 1. Tạo private key 2048-bit
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 2. Ghi file private key (PKCS#1)
// 	privFile, err := os.Create("mb.key")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer privFile.Close()

// 	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
// 	privPem := &pem.Block{
// 		Type:  "RSA PRIVATE KEY",
// 		Bytes: privBytes,
// 	}
// 	if err := pem.Encode(privFile, privPem); err != nil {
// 		panic(err)
// 	}
// 	println(" Đã tạo file mb.key")

// 	// 3. Ghi file public key
// 	pubFile, err := os.Create("mb.pub")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer pubFile.Close()

// 	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	pubPem := &pem.Block{
// 		Type:  "PUBLIC KEY",
// 		Bytes: pubBytes,
// 	}
// 	if err := pem.Encode(pubFile, pubPem); err != nil {
// 		panic(err)
// 	}
// 	println(" Đã tạo file mb.pub")
// }
