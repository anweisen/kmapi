package scraper

import (
  "fmt"
  "strconv"
  "strings"
)

func FormatNumberWithLeadingZero(number int) string {
  if number < 10 {
    return "0" + strconv.Itoa(number)
  }
  return strconv.Itoa(number)
}

func IsoFormatDate(day int, month int, year int) string {
  return strconv.Itoa(year) + "-" + FormatNumberWithLeadingZero(month) + "-" + FormatNumberWithLeadingZero(day)
}

func ConvertGermanDateToIso(germanDateParts []string) (string, error) {
  // German date format "dd. MMMM yyyy", e.g. "01. Januar 2020"
  dayString := germanDateParts[0]
  monthName := germanDateParts[1]
  yearString := germanDateParts[2]

  if strings.HasSuffix(dayString, ".") {
    dayString = dayString[:len(dayString)-1]
  }
  day, err := strconv.Atoi(dayString)
  if err != nil {
    return "", err
  }

  month := ParseGermanMonth(monthName)
  if month == -1 {
    return "", fmt.Errorf("invalid month name: %s", monthName)
  }

  if strings.HasSuffix(yearString, ",") {
    yearString = yearString[:len(yearString)-1]
  }
  year, err := strconv.Atoi(yearString)
  if err != nil {
    return "", err
  }

  return IsoFormatDate(day, month, year), nil
}

func ConvertGermanWeekDateToIso(germanWeekDateParts []string) (start string, end string, err error) {
  // German week date format "DDDD, dd. MMMM [yyyy] [-/bis/mit] DDDD, dd. MMMM yyyy", e.g. "Montag, 1. Januar - Sonntag, 7. Januar 2020"
  startDayString := germanWeekDateParts[1]
  startMonthName := germanWeekDateParts[2]
  endDayString := germanWeekDateParts[len(germanWeekDateParts)-3]
  endMonthName := germanWeekDateParts[len(germanWeekDateParts)-2]
  endYearString := germanWeekDateParts[len(germanWeekDateParts)-1]

  if strings.HasSuffix(startDayString, ".") {
    startDayString = startDayString[:len(startDayString)-1]
  }
  startDay, err := strconv.Atoi(startDayString)
  if err != nil {
    return
  }

  startMonth := ParseGermanMonth(startMonthName)
  if startMonth == -1 {
    err = fmt.Errorf("invalid month name: %s", startMonthName)
    return
  }

  if strings.HasSuffix(endDayString, ".") {
    endDayString = endDayString[:len(endDayString)-1]
  }
  endDay, err := strconv.Atoi(endDayString)
  if err != nil {
    return
  }

  endMonth := ParseGermanMonth(endMonthName)
  if endMonth == -1 {
    err = fmt.Errorf("invalid month name: %s", endMonthName)
    return
  }

  endYear, err := strconv.Atoi(endYearString)
  if err != nil {
    return
  }

  start = IsoFormatDate(startDay, startMonth, endYear)
  end = IsoFormatDate(endDay, endMonth, endYear)
  return
}

func ParseGermanMonth(formattedMonthName string) int {
  switch formattedMonthName {
  case "Januar":
    return 1
  case "Februar":
    return 2
  case "März":
    return 3
  case "April":
    return 4
  case "Mai":
    return 5
  case "Juni":
    return 6
  case "Juli":
    return 7
  case "August":
    return 8
  case "September":
    return 9
  case "Oktober":
    return 10
  case "November":
    return 11
  case "Dezember":
    return 12
  default:
    return -1
  }
}
