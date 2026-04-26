package e2e_test

import (
	"bookstore-api/ent"
	"bookstore-api/internal/application"
	"bookstore-api/internal/infrastructure/persistence"
	apihandler "bookstore-api/internal/interface/api"
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

var (
	testServer *httptest.Server
	testClient *ent.Client
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

// run is split from TestMain so that deferred cleanups execute before os.Exit.
func run(m *testing.M) int {
	ctx := context.Background()

	container, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("bookstore_test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
	)
	if err != nil {
		panic("failed to start MySQL container: " + err.Error())
	}
	defer container.Terminate(ctx) //nolint:errcheck

	dsn, err := container.ConnectionString(ctx, "charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to get connection string: " + err.Error())
	}

	testClient, err = ent.Open("mysql", dsn)
	if err != nil {
		panic("failed to open ent client: " + err.Error())
	}
	defer testClient.Close()

	if err = testClient.Schema.Create(ctx); err != nil {
		panic("failed to run schema migration: " + err.Error())
	}

	testServer = newServer()
	defer testServer.Close()

	return m.Run()
}

func newServer() *httptest.Server {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(apihandler.CorrelationID(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	config := huma.DefaultConfig("Bookstore API", "1.0.0")
	config.SchemasPath = ""
	config.CreateHooks = nil
	humaAPI := humagin.New(r, config)

	repo := persistence.NewEntBookRepository(testClient)
	service := application.NewBookService(repo)
	bookHandler := apihandler.NewBookHandler(service)
	recordHandler := apihandler.NewRecordHandler()
	bookHandler.Register(humaAPI)
	recordHandler.Register(humaAPI)

	return httptest.NewServer(r)
}

// cleanBooks deletes all books between tests to keep them independent.
func cleanBooks(t *testing.T) {
	t.Helper()
	_, err := testClient.Book.Delete().Exec(context.Background())
	if err != nil {
		t.Fatalf("cleanBooks: %v", err)
	}
}
