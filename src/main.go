package main

import (
  "github.com/gofiber/fiber/v3"
  "github.com/gofiber/fiber/v3/middleware/cors"
  "github.com/gofiber/fiber/v3/middleware/static"
  "golang.org/x/sync/singleflight"
  "kmapi/src/api/cache"
  "kmapi/src/api/database"
  "kmapi/src/routes"
  "os"
)

func main() {
  var singleFlightGroup singleflight.Group
  handler := routes.AppHandlerEmbed{
    Cache:        cache.NewRedisCache(),
    Database:     database.NewMongoDatabase(),
    Singleflight: &singleFlightGroup,
  }
  server := fiber.New()

  server.Use(cors.New(cors.Config{
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"GET"},
  }))

  server.Get("/", static.New("./static"))
  server.Get("/by/gym/abi/next", handler.HandleGetByGymAbiNext)
  server.Get("/by/gym/abi/:year", handler.HandleGetByGymAbiYear)

  bind, present := os.LookupEnv("BIND_ADDR")
  if !present {
    bind = ":5000"
  }

  err := server.Listen(bind)
  if err != nil {
    panic(err)
  }
}
