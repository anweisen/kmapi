package main

import (
  "github.com/gofiber/fiber/v3"
  "kmapi/src/api/cache"
  "os"
)

func main() {
  handler := routes.AppHandlerEmbed{
    Cache: cache.NewRedisCache(),
  }
  server := fiber.New()

  bind, present := os.LookupEnv("BIND_ADDR")
  if !present {
    bind = ":5000"
  }

  err := server.Listen(bind)
  if err != nil {
    panic(err)
  }
}
