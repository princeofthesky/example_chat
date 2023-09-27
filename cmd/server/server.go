package main

import (
	"github.com/princeofthesky/example_chat/middlewares"
	"github.com/princeofthesky/example_chat/repository"
	"github.com/princeofthesky/example_chat/token"
	"github.com/princeofthesky/example_chat/trace_log"
	"github.com/princeofthesky/example_chat/transport"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	enableLog := os.Getenv("APP_SHOW_LOG")
	fileNameLog := os.Getenv("FILE_NAME_LOG")
	trace_log.Init(fileNameLog, enableLog)
	token.InitJWTSECRETKEY()
}

func main() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			f, _ := os.Create("cpu.prof")
			pprof.StartCPUProfile(f)
			time.Sleep(1 * time.Minute)
			pprof.StopCPUProfile()
			f.Close()

			f, _ = os.Create("mem.prof")
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}()

	router := gin.New()
	// router.Use(middlewares.JwtAuthMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// skConnect := connection.NewSocketInstance()
	// go skConnect.Initialize()
	// handler := transport.NewHandler(systemredis.Redis, mongodb.MongoClient, skConnect)

	var cprefix = os.Getenv("COMPANY_PREFIX")
	var liveprefix = os.Getenv("LIVE_PREFIX")
	var instanceName = os.Getenv("INSTANCE_NAME")
	// repo := repository.NewSocketSky(SystemRedisClient, MongoConnectClient)
	// Handler := skytransport.NewSkyHandler(repo)

	// skyRepo := repository.NewSky(SystemRedisClient, MongoConnectClient)

	skyRepo := repository.NewSocketLive(cprefix, liveprefix, instanceName)

	// skyRepo := repository.NewSocketLive(SystemRedisClient, MongoConnectClient, SystemRedisClientSync, cprefix, liveprefix)

	skyHandler := transport.NewSkyHandler(skyRepo)
	// vao chi doc tin nhan
	authorized := router.Group("/ws")
	authorized.Use(middlewares.JwtAuthMiddleware())
	{

		authorized.GET("/:userId/:userConnectionId", skyHandler.OpenConnection)
		// router.GET("/ws/:userId/:userConnectionId", skyHandler.OpenConnection)
	}

	//Lee check token

	router.Run(":" + os.Getenv("APP_PORT_SOCKET"))
}
