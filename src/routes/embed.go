package routes

import (
  "golang.org/x/sync/singleflight"
  "kmapi/src/api/cache"
)

type AppHandlerEmbed struct {
  Cache        cache.Cache
  Singleflight *singleflight.Group
}
