package cache

type Cache interface {
  SetJson(key string, value any) error
  GetJson(key string, dest any) error
}
