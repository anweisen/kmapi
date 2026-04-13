package routes

import (
  "github.com/gofiber/fiber/v3"
  "kmapi/src/api"
  "kmapi/src/api/cache"
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
  // (3. return historical data from db)

  var cachedData *api.ByAbiYearData
  yearStr := strconv.Itoa(year)
  err := app.Cache.GetJson(cache.ConstructCacheKey(api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, yearStr), &cachedData)
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
    return DoScrapeBavariaAndCache(app) // data, err
  })

  if err != nil {
    return nil, err
  }

  data := result.(*scraper.ByScrapeData)

  abiData, exists := data.GymAbiYearData[year]
  if exists {
    return &abiData, nil
  }

  // TODO(feat): should first check db for historical data
  return nil, nil // TODO better error handling?
}

func DoScrapeBavariaAndCache(app AppHandlerEmbed) (*scraper.ByScrapeData, error) {
  data, err := scraper.ScrapeBavaria()
  if err != nil {
    return nil, err
  }

  for year, abiData := range data.GymAbiYearData {
    err = app.Cache.SetJson(cache.ConstructCacheKey(api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, strconv.Itoa(year)), abiData)
    if err != nil {
      println("Error caching data for year", year, ":", err.Error())
      return nil, err
    } else {
      println("Cached data for year", year)
    }
  }
  return data, nil
}
