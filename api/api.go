package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sillyhatxu/convenient-utils/response"
	"github.com/sillyhatxu/elasticsearch-ui/config"
	"github.com/sillyhatxu/elasticsearch-ui/dto"
	"github.com/sillyhatxu/elasticsearch-ui/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func InitialAPI() error {
	logrus.Info("---------- initial internal api start ----------")
	router := SetupRouter()

	//corsConfig := cors.DefaultConfig()
	//corsConfig.AllowOrigins = []string{"http://localhost"}
	//router.Use(cors.New(corsConfig))
	//router.Use(cors.Default())

	//router.Use()

	return router.Run(config.Conf.ServerHost)
}

func CORSMiddleware() gin.HandlerFunc {
	//corsConfig := cors.DefaultConfig()
	//corsConfig.AllowAllOrigins = true
	//return cors.New(corsConfig)

	//return cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://localhost:4200"},
	//	AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		//return origin == "https://github.com"
	//		return false
	//	},
	//	MaxAge: 12 * time.Hour,
	//})
	//return cors.Middleware(cors.Config{
	//	Origins:         "*",
	//	Methods:         "GET, PUT, POST, DELETE",
	//	RequestHeaders:  "Origin, Authorization, Content-Type, Access-Control-Allow-Origin",
	//	ExposedHeaders:  "",
	//	MaxAge:          50 * time.Second,
	//	Credentials:     true,
	//	ValidateHeaders: false,
	//})
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			//return origin == "https://github.com"
			return true
		},
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"message": "OK"}) })
	//router.GET("/initial", initial).Use(CORSMiddleware())
	router.GET("/initial", initial)
	group := router.Group("/elasticsearch")
	{
		group.POST("/connect", connect)
		group.GET("/version", version)
		group.GET("/health", health)
		group.GET("/cluster-stats", clusterStats)
		group.GET("/indices", indices)
		group.GET("/mappings", mappings)
		//group.PUT("/:id", update)
		//group.DELETE("/:id", delete)
		//group.POST("", add)
	}
	return router
}

func initial(context *gin.Context) {
	context.JSON(http.StatusOK, response.ServerSuccess(config.Conf, nil))
}

func connect(context *gin.Context) {
	var requestBody dto.ConnectDTO
	err := context.ShouldBindJSON(&requestBody)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerParamsValidateError(nil, err.Error()))
		return
	}
	config.Conf.URL = requestBody.ESURL
	resp, code, err := service.Ping(config.Conf.URL)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerError(nil, err.Error(), nil))
		return
	}
	if code != http.StatusOK {
		context.JSON(http.StatusOK, response.ServerError(resp, fmt.Sprintf("connect error;http code : %v", code), nil))
		return
	}
	context.JSON(http.StatusOK, response.ServerSuccess(resp, nil))
}

func version(context *gin.Context) {
	logrus.Infof("api [/elasticsearch/version]")
	resp, err := service.Version(config.Conf.URL)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerError(nil, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, response.ServerSuccess(resp, nil))
}

func health(context *gin.Context) {
	logrus.Infof("api [/elasticsearch/health]")
	resp, err := service.Health(config.Conf.URL)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerError(nil, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, response.ServerSuccess(resp, nil))
}

func clusterStats(context *gin.Context) {
	logrus.Infof("api [/elasticsearch/cluster-stats]")
	resp, err := service.ClusterStats(config.Conf.URL)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerError(nil, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, response.ServerSuccess(resp, nil))
}

//index
//docs.count
//status
//health
//store.size
func indices(context *gin.Context) {
	logrus.Infof("api [/elasticsearch/indices]")
	resp, err := service.Indices(config.Conf.URL)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerError(nil, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, response.ServerSuccess(resp, nil))
}

func mappings(context *gin.Context) {
	logrus.Infof("api [/elasticsearch/mappings]")
	resp, err := service.GetMappings(config.Conf.URL)
	if err != nil {
		context.JSON(http.StatusOK, response.ServerError(nil, err.Error(), nil))
		return
	}
	context.JSON(http.StatusOK, response.ServerSuccess(resp, nil))
}
