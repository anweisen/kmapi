package api

import (
  "os"
  "strconv"
)

func GetEnvInt(key string, defaultValue int) int {
  envString, present := os.LookupEnv(key)
  if !present {
    return defaultValue
  }
  parsedInt, err := strconv.Atoi(envString)
  if err != nil {
    return defaultValue
  }
  return parsedInt
}
