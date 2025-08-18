package server

import (
	"context"
	"europm/internal/bank"
	"europm/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var srv *http.Server

func Start() (err error) {
	hostName, _ := os.Hostname()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(otelgin.Middleware(hostName))

	// Chuẩn hoá: bỏ "/" cuối nếu có → "/api/v1"
	prefix := strings.TrimSuffix(config.GetString("api_prefix"), "/")

	// Health cho chính prefix → "/api/v1/"
	r.GET(prefix+"/", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.HEAD(prefix+"/", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Group các API con → "/api/v1/..."
	rg := r.Group(prefix + "/")
	{
		rg.GET("account", bank.CallAccount)
		rg.POST("order-update", bank.OrderUpdateHandler)
	}

	srv = &http.Server{
		Addr:    ":" + config.GetString("server.port"), // 8433
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()
	return nil
}

func Stop() {
	// The context is used to inform the server it has 90 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	fmt.Println("defer cancel()")
	defer cancel()
	fmt.Println("srv.Shutdown(ctx)")
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown:")
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exiting")
}
