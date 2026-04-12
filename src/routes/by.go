package routes

import (
  "github.com/gofiber/fiber/v3"
  "kmapi/src/api"
  "kmapi/src/api/cache"
  "kmapi/src/scraper"
  "strconv"
)

func (app AppHandlerEmbed) HandleGetByGymAbiYear(ctx fiber.Ctx) error {
  // get year from path param
  yearParam := ctx.Params("year")
  year, err := strconv.Atoi(yearParam)
  if err != nil {
    return ctx.SendStatus(fiber.StatusBadRequest)
  }

  var cachedData *api.ByAbiYearData
  err = app.Cache.GetJson(cache.ConstructCacheKey(api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, yearParam), &cachedData)
  if err != nil {
    return ctx.SendStatus(fiber.StatusInternalServerError)
  }
  if cachedData == nil {
    println("Cache miss for year", year)
  } else {
    println("Cache hit for year", year)
    return ctx.JSON(cachedData)
  }

  // TODO(perf): single flight
  data := DoScrapeBavariaAndCache(app)

  abiData, exists := data.GymAbiYearData[year]
  if !exists {
    // TODO(feat): should first check db for historical data
    return ctx.SendStatus(fiber.StatusNotFound)
  }

  return ctx.JSON(abiData)
}

func DoScrapeBavariaAndCache(app AppHandlerEmbed) *scraper.ByScrapeData {
  data := scraper.ScrapeBavaria()
  for year, abiData := range data.GymAbiYearData {
    err := app.Cache.SetJson(cache.ConstructCacheKey(api.KeyBavaria, api.KeyGymnasium, api.KeyAbi, strconv.Itoa(year)), abiData)
    if err != nil {
      println("Error caching data for year", year, ":", err.Error())
    } else {
      println("Cached data for year", year)
    }
  }
  return data
}
