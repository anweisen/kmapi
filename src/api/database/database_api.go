package database

import (
  "kmapi/src/api/cache"
  "strconv"
  "time"
)

type Database interface {
  UpsertArchiveEntry(entry ArchiveEntry[any]) error
  FindArchiveEntry(id string, dest any) error
}

type ArchiveEntry[T any] struct {
  Id     string    `json:"-" bson:"_id,omitempty"` // = cache key
  Stored time.Time `json:"stored" bson:"stored"`
  // unique cache key: {state}:{school}:{category}:{year} split into fields; future proofing for more complex queries
  State    string `json:"state" bson:"state"`
  School   string `json:"school" bson:"school"`
  Category string `json:"category" bson:"category"`
  Year     int    `json:"year" bson:"year"`

  Data T `json:"data" bson:"data"`
}

func ConstructArchiveEntry[T any](state string, school string, category string, year int, data T) ArchiveEntry[T] {
  return ArchiveEntry[T]{
    Id:       cache.ConstructCacheKey(state, school, category, strconv.Itoa(year)),
    Stored:   time.Now(),
    State:    state,
    School:   school,
    Category: category,
    Year:     year,
    Data:     data,
  }
}
