package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	gin "github.com/gin-gonic/gin"

	"go.vxn.dev/swis/v5/pkg/alvax"
	"go.vxn.dev/swis/v5/pkg/auth"
	"go.vxn.dev/swis/v5/pkg/backups"
	"go.vxn.dev/swis/v5/pkg/business"
	"go.vxn.dev/swis/v5/pkg/config"
	"go.vxn.dev/swis/v5/pkg/core"
	"go.vxn.dev/swis/v5/pkg/depots"
	"go.vxn.dev/swis/v5/pkg/dish"
	"go.vxn.dev/swis/v5/pkg/finance"
	"go.vxn.dev/swis/v5/pkg/infra"
	"go.vxn.dev/swis/v5/pkg/links"
	"go.vxn.dev/swis/v5/pkg/news"
	"go.vxn.dev/swis/v5/pkg/projects"
	"go.vxn.dev/swis/v5/pkg/queue"
	"go.vxn.dev/swis/v5/pkg/roles"
	"go.vxn.dev/swis/v5/pkg/system"
	"go.vxn.dev/swis/v5/pkg/users"
)

type server struct {
	listener net.Listener
	router   *gin.Engine
	srv      *http.Server

	once sync.Once
	wg   *sync.WaitGroup
}

func newServer() *server {
	return &server{}
}

func (s *server) Run() {
	s.init()
	s.handleSignals()
	s.setupRouter()
	s.setupServer()
	s.serve()

	// Shudown invoked here: wait for graceful goroutine to finish.
	s.wg.Wait()

	log.Println("the HTTP server has stopped serving new connections, program exit")
}

func (s *server) init() {
	s.once.Do(func() {
		var wg sync.WaitGroup
		s.wg = &wg
	})
}

func (s *server) handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// Wait for a signal.
		sig := <-sigs
		signal.Stop(sigs)

		log.Printf("trap signal: %s, graceful shutdown invoked...", sig.String())

		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		defer func() {
			if err := s.listener.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		// Try to gracefully shutdown the HTTP server.
		if err := s.srv.Shutdown(sctx); err != nil {
			log.Fatal(err)

			// Forcefully close the HTTP server.
			if err := s.srv.Close(); err != nil {
				log.Fatal(err)
			}

			return
		}

		log.Println("graceful shutdown completed")
	}()
}

func (s *server) setupRouter() {
	// Blank gin without any middleware.
	s.router = gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	s.router.Use(gin.Recovery())

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	s.router.Use(config.JSONLogMiddleware(), config.CORSMiddleware())

	// Serve vxn-dev internal favicon.
	s.router.StaticFile("/favicon.ico", "./favicon.ico")

	// @Summary Simple ping-pong route
	// @Description Simple ping-pong route
	// @Success 200
	// @Router /ping [get]
	// @Router /ping [head]
	// @Success 200 {string} "pong"
	// Very simple LE support --- https://github.com/gin-gonic/gin#support-lets-encrypt.
	s.router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	s.router.HEAD("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Default 404 route
	s.router.NoRoute(func(c *gin.Context) {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "unknown route, or disallowed method",
		})
	})

	// Use custom swapi Auth middlewares --- token auth
	// AuthenticationMiddleware takes care of token verification against loaded Users data structure.
	// AuthorizationMiddleware checks Access Control List (ACL) and the right for dangerous methods usage.
	s.router.Use(auth.AuthenticationMiddleware())
	s.router.Use(auth.AuthorizationMiddleware())

	// Root path
	s.router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"header_title": "sakalWebIS v5 RESTful API -- root route",
			"message":      "welcome to swis, " + auth.Params.User.Name + "!",
			"code":         http.StatusOK,
			"app_env": gin.H{
				"app_mode_environment": os.Getenv("APP_ENVIRONMENT"),
				"instance_name":        os.Getenv("HOSTNAME"),
				"alpine_version":       os.Getenv("ALPINE_VERSION"),
				"app_version":          os.Getenv("APP_VERSION"),
				"golang_version":       os.Getenv("GOLANG_VERSION"),
			},
			"timestamp": time.Now().Unix(),
			"user": gin.H{
				"acl":   auth.Params.ACL,
				"roles": auth.Params.Roles,
			},
		})
	})

	//
	// swis pkg registration
	//

	// Preregister system cache to track registered packages.
	core.MountPackage(s.router, system.Package)

	// Bulk registration and mounting of packages.
	core.MountMany(s.router, &system.Cache,
		alvax.Package,
		backups.Package,
		business.Package,
		depots.Package,
		dish.Package,
		finance.Package,
		infra.Package,
		links.Package,
		news.Package,
		projects.Package,
		queue.Package,
		roles.Package,
		users.Package,
	)

	// Initialize other components.
	dish.Dispatcher = dish.NewDispatcher()

	// Attach router to http.Server and start it, check for SERVER_PORT env variable.
	if os.Getenv("SERVER_PORT") == "" {
		log.Fatal("SERVER_PORT environment variable not provided! refusing to start the server...")
	}
}

func (s *server) setupServer() {
	if s.router == nil {
		log.Fatal("router is not initialized")
	}

	var err error
	if s.listener, err = net.Listen("tcp", ":"+os.Getenv("SERVER_PORT")); err != nil {
		panic(err)
	}

	s.srv = &http.Server{
		Addr:        s.listener.Addr().String(),
		Handler:     s.router,
		ReadTimeout: 10 * time.Second,
		// WriteTimeout: 10 * time.Second,
		// = 1 * 2^23 = 1,048,576 * 8
		MaxHeaderBytes: 1 << 23,
		// Use config.CORSMiddleware()
		// DisableGeneralOptionsHandler: true,
	}
}

func (s *server) serve() {
	log.Printf("init done, starting the HTTP server (v%s)", os.Getenv("APP_VERSION"))

	if err := s.srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
