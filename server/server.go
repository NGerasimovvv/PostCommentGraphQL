package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NGerasimovvv/GraphQL/graph"
	"github.com/NGerasimovvv/GraphQL/internal/gateway"
	"github.com/NGerasimovvv/GraphQL/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
)

func graphqlHandler(storage storage.Storage) gin.HandlerFunc {
	postGateway := gateway.NewPostGateway(storage)
	commentGateway := gateway.NewCommentGateway(storage)

	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			PostGateway:    postGateway,
			CommentGateway: commentGateway,
		},
	}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func InitServer(storage storage.Storage) {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.POST("/graphql", graphqlHandler(storage))
	r.GET("/", playgroundHandler())
	log.Println("connect to http://localhost:8000/ for GraphQL playground")
	log.Fatal(r.Run(":8000"))
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
