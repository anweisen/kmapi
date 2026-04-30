package routes

import (
  "golang.org/x/sync/singleflight"
  "kmapi/src/api/cache"
  "kmapi/src/api/database"
)

type AppHandlerEmbed struct {
  Cache        cache.Cache
  Database     database.Database
  Singleflight *singleflight.Group
}
