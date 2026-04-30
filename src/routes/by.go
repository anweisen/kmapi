package routes

import (
  "errors"
  "github.com/gofiber/fiber/v3"
  "go.mongodb.org/mongo-driver/v2/mongo"
  "kmapi/src/api"
  "kmapi/src/api/cache"
  "kmapi/src/api/database"
  "kmapi/src/scraper"
  "strconv"
  "time"
)

func (app AppHandlerEmbed) HandleGetByGymAbiYear(ctx fiber.Ctx) error {
  // get year from path param
  yearParam := ctx.Params("year")
  year, err := strconv.Atoi(yearParam)
  if err != nil {
    return ctx.SendStatus(fiber.StatusBadRequest)
  }

  data, err := GetBavariaAbiDataForYear(app, year)
  if err != nil {
    println("Error getting data for year", year, ":", err.Error())
    return ctx.SendStatus(fiber.StatusInternalServerError)
  }
  if data == nil {
    return ctx.SendStatus(fiber.StatusNotFound)
  }

  return ctx.JSON(data)
}

func (app AppHandlerEmbed) HandleGetByGymAbiNext(ctx fiber.Ctx) error {
  currentYear := time.Now().Year()

  currentYearData, err := GetBavariaAbiDataForYear(app, currentYear)
  if err != nil {
    println("Error getting data for current year", currentYear, ":", err.Error)
    return ctx.SendStatus(fiber.StatusInternalServerError)
  }
  if currentYearData != nil {
    currentYearGraduationDate, err := time.Parse(time.DateOnly, currentYearData.GraduationDate.Date)
    if err != nil {
      println("Error parsing graduation date for current year", currentYear, ":", err.Error())
      return ctx.SendStatus(fiber.StatusInternalServerError)
    }

    // hasn't passed until the day after graduation date
    if time.Now().Before(currentYearGraduationDate.AddDate(0, 0, 1)) {
      return ctx.JSON(currentYearData)
    }
  }

  // if no data for current year or graduation date has passed, try next year
  nextYear := currentYear + 1
  nextYearData, err := GetBavariaAbiDataForYear(app, nextYear)
  if err != nil {
    println("Error getting data for next year", nextYear, ":", err.Error())
    return ctx.SendStatus(fiber.StatusInternalServerError)
  }
  if nextYearData == nil {
    return ctx.SendStatus(fiber.StatusNotFound)
  }

  return ctx.JSON(nextYearData)
}

func GetBavariaAbiDataForYear(app AppHandlerEmbed, year int) (*api.ByAbiYearData, error) {
  // 1. try cache
  // 2. scrape new data (and cache it)
  // 3. return archived data from db (and cache it)

  var cachedData *api.ByAbiYearData
  yearStr := strconv.Itoa(year)
  cacheKey := cache.ConstructCacheKey(api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, yearStr)

  err := app.Cache.GetJson(cacheKey, &cachedData)
  if err != nil {
    return nil, err
  }
  if cachedData != nil {
    println("Cache hit for year", year)
    return cachedData, nil
  }

  println("Cache miss for year", year)

  // perf: singleflight to prevent concurrent scrapes
  result, err, _ := app.Singleflight.Do(api.KeyBavaria, func() (interface{}, error) {
    return ScrapeBavariaAndCacheArchive(app) // data, err
  })

  if err != nil {
    return nil, err
  }

  data := result.(*scraper.ByScrapeData)

  abiData, exists := data.GymAbiYearData[year]
  if exists {
    return &abiData, nil
  }

  // check db for archived data
  var archiveEntry database.ArchiveEntry[api.ByAbiYearData]
  err = app.Database.FindArchiveEntry(cacheKey, &archiveEntry)
  if errors.Is(err, mongo.ErrNoDocuments) {
    return nil, nil // 404
  }
  if err != nil {
    return nil, err
  }

  err = app.Cache.SetJson(cacheKey, &archiveEntry.Data) // cache archived data -> reduce database load
  if err != nil {
    return nil, err
  }

  return &archiveEntry.Data, nil
}

func ScrapeBavariaAndCacheArchive(app AppHandlerEmbed) (*scraper.ByScrapeData, error) {
  data, err := scraper.ScrapeBavaria()
  if err != nil {
    return nil, err
  }

  for year, abiData := range data.GymAbiYearData {
    _ = app.Cache.SetJson(cache.ConstructCacheKey(api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, strconv.Itoa(year)), abiData)
    _ = app.Database.UpsertArchiveEntry(database.ConstructArchiveEntry[any](api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, year, abiData))
  }
  return data, nil
}
