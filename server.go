package main

import (
	"context"
	"log"

	"github.com/Amitmahato/file-server/infrastructures"
	"github.com/Amitmahato/file-server/routes"
	"github.com/Amitmahato/file-server/services"
	"github.com/Amitmahato/file-server/utils"
	"go.uber.org/fx"
)



func main() {
	utils.LoadEnv()

	fx.New(
		fx.Options(infrastructures.Module),
		fx.Provide(services.NewGCPStorageService),
		fx.Provide(routes.NewRoutes),
		fx.Invoke(func (
			lifecycle fx.Lifecycle,
			router infrastructures.GinRouter,
			routes routes.Routes,
			){
				lifecycle.Append(fx.Hook{
					OnStart: func(c context.Context) error {
						log.Println("Starting Application")
						log.Println("Setting Up Routes")
						routes.Setup()
						go func ()  {
							log.Println("Starting Server")
							port := utils.GetEnvWithKey("SERVER_PORT")
							router.Gin.Run(":"+port)						
						}()
						return nil
					},
					OnStop: func(c context.Context) error {
						log.Println("Stopping Application")
						return nil
					},
				})
			}),
	).Run()
}
