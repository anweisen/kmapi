package cache

type Cache interface {
  SetJson(key string, value any) error
  GetJson(key string, dest any) error
}

func ConstructCacheKey(state string, school string, category string, year string) string {
  return state + ":" + school + ":" + category + ":" + year
}
