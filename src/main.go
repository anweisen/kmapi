package main

import (
  "github.com/gofiber/fiber/v3"
  "kmapi/src/api/cache"
  "kmapi/src/routes"
  "os"
)

func main() {
  var singleFlightGroup singleflight.Group
  handler := routes.AppHandlerEmbed{
    Cache:        cache.NewRedisCache(),
    Singleflight: &singleFlightGroup,
  }
  server := fiber.New()

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
