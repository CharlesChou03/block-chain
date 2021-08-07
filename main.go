package main

import (
	"github.com/CharlesChou03/_git/block-chain.git/config"
	_ "github.com/CharlesChou03/_git/block-chain.git/docs"
	"github.com/CharlesChou03/_git/block-chain.git/internal/db"
	"github.com/CharlesChou03/_git/block-chain.git/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setup() {
	config.Setup()
	db.MySQLDB = db.SetupMySQLDB()
	db.MySQLDB.CreateBlockChainTable()
	db.RedisDB = db.SetupRedisDB()
}

func setupRouter() *gin.Engine {
	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", handlers.HealthHandler)
	r.GET("/version", handlers.VersionHandler)

	r.GET("/blocks", handlers.GetLatestNBlockHandler)
	r.GET("/blocks/:id", handlers.GetBlockByNumHandler)
	r.GET("/transaction/:txHash", handlers.GetTransactionByHashHandler)

	return r
}

// @title Swagger
// @version 0.0.1
func main() {
	setup()
	defer db.MySQLDB.Close()
	defer db.RedisDB.Close()
	r := setupRouter()
	r.Run(":9999")
}
