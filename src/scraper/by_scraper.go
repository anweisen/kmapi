package scraper

import (
  "fmt"
  "github.com/gocolly/colly/v2"
  "kmapi/src/api"
  "strconv"
  "strings"
)

var URL = "https://www.km.bayern.de/termine/pruefungen-und-zeugnisse"

var (
  StateNone       = 0
  StateAbiInit    = 1
  StateAbiWritten = 2
  StateAbiOral    = 3
  StateAbiReport  = 4
)

type ByScrapeData struct {
  GymAbiYearData map[int]api.ByAbiYearData
}

func ScrapeBavaria() (*ByScrapeData, error) {
  c := colly.NewCollector()

  scrapeData := &ByScrapeData{}

  // "Gymnasium" accordion id="rxContentId4"
  c.OnHTML("#rxContentId4", func(elGym *colly.HTMLElement) {
    gymAbiYearData := make(map[int]api.ByAbiYearData) // collect data into this
    // "Gymnasium"-entries class="rxModuleTextImage", ignore ".rxModuleDivider"
    // - "Abiturprüfungen 2026"
    // - "Abiturprüfungen 2027"
    // - "Termine der Jahrgangsstufentests im Schuljahr 2026/2027"
    // - "Termine der Jahrgangsstufentests im Schuljahr 2027/2028"
    // - "Termine der Besonderen Prüfung im Anschluss an das Schuljahr 2025/2026"
    elGym.ForEach(".rxModuleTextImage", func(_ int, elAbiCurrent *colly.HTMLElement) {

      var state = StateNone
      var year int
      var writtenDates []api.ByWrittenExamDate
      var oralWeeks []api.ByOralExamDate
      var practicalDates []api.ByPracticalExamDate
      var graduationDate *api.ByGraduationDate

      elAbiCurrent.ForEach("p", func(_ int, elAbiCurrentP *colly.HTMLElement) {
        var text = strings.TrimSpace(elAbiCurrentP.Text)

        // <strong>HEADLINES</strong>
        if strings.HasPrefix(text, "Abiturprüfung") {
          state = StateAbiInit
          split := strings.Split(text, " ")
          if len(split) < 2 {
            // TODO
          }

          parsedYear, err := strconv.Atoi(split[1])
          if err != nil {
            // TODO
          }
          year = parsedYear
        } else if strings.HasPrefix(text, "Schriftlicher Teil") {
          state = StateAbiWritten
          // same dom-element <p> contains title and first date (not separated in text = <br>) -> trim and process
          text = strings.TrimPrefix(text, "Schriftlicher Teil")
        } else if strings.HasPrefix(text, "Kolloquiumsprüfungen") {
          state = StateAbiOral
          // same dom-element <p> contains title and first date (not separated in text = <br>) -> trim and process
          text = strings.TrimPrefix(text, "Kolloquiumsprüfungen")
        } else if strings.HasPrefix(text, "Abiturzeugnis") {
          state = StateAbiReport
        }
        // trim leading and trailing whitespace: some date entries <p> have leading spaces for no clear reason
        text = strings.TrimSpace(text)

        // TODO(refactor): extract each to own func
        switch state {
        case StateAbiWritten:
          println("Written", text)
          // 2026:
          // - [DATE,3] [SUBJECT_NAME,1] (grundlegendes und erhöhtes Anforderungsniveau)
          // - [DATE,3] [Mathematik/Deutsch,1]
          // - [DATE,3] alle weiteren Abiturprüfungsfächer auf grundlegendem/erhöhtem Anforderungsniveau
          // 2027:
          // - [DATE,3] [SUBJECT_NAME,1] (grundlegendes und erhöhtes Anforderungsniveau)
          // - [DATE,3] [Mathematik/Deutsch,1]
          // ! [DATE,3 alle Prüfungsfächer auf erhöhtem/grundlegendem Anforderungsniveau (mit Ausnahme von [SUBJECT_NAME,n], )

          dateParts := strings.SplitN(text, " ", 4)
          if len(dateParts) < 4 { // changed date format?
            // TODO
          }
          isoDate, err := ConvertGermanDateToIso(dateParts)
          if err != nil { // changed date format?
            // TODO
          }

          isRemaining := strings.Contains(dateParts[3], "alle") || strings.Contains(dateParts[3], "weitere")

          var subjectName, descriptionPart string
          if isRemaining {
            // - alle weiteren Abiturprüfungsfächer auf grundlegendem/erhöhtem Anforderungsniveau
            // - alle Abiturprüfungsfächer auf erhöhtem/grundlegendem Anforderungsniveau (mit Ausnahme von [SUBJECT_NAME,n], )
            endNameIndex := strings.Index(dateParts[3], "Abiturprüfungsfächer") + len("Abiturprüfungsfächer")
            subjectName = dateParts[3][:endNameIndex]
            if len(dateParts[3]) > endNameIndex { // the current format requires description for remaining subjects, but just in case
              descriptionPart = dateParts[3][1+endNameIndex:] // trim leading space
            } else {
              descriptionPart = ""
            }
          } else {
            // - [SUBJECT_NAME,1] (grundlegendes und erhöhtes Anforderungsniveau)
            // - [Mathematik/Deutsch,1]
            remainingParts := strings.SplitN(dateParts[3], " ", 2)
            subjectName = remainingParts[0]
            if len(remainingParts) > 1 {
              descriptionPart = remainingParts[1]
            } else {
              descriptionPart = ""
            }
          }

          ea, ga := true, false // für Mathematik/Deutsch: nicht spezifiziert, immer eA
          if len(descriptionPart) > 1 {
            ea = strings.Contains(descriptionPart, "erhöht")
            ga = strings.Contains(descriptionPart, "grundlegend")
          }
          var excludedSubjects []string
          if isRemaining && strings.Contains(descriptionPart, "mit Ausnahme von") {
            // - alle Abiturprüfungsfächer auf erhöhtem/grundlegendem Anforderungsniveau (mit Ausnahme von [SUBJECT_NAME,n], )
            startIndex := strings.Index(descriptionPart, "mit Ausnahme von") + len("mit Ausnahme von")
            excludedPart := strings.TrimSpace(strings.TrimSuffix(descriptionPart[startIndex:], ")")) // trim leading space and trailing parenthesis
            excludedSubjects = strings.Split(excludedPart, ", ")
          }

          writtenDates = append(writtenDates, api.ByWrittenExamDate{
            Date:          isoDate,
            FormattedDate: strings.Join(dateParts[:3], " "),
            Subject:       subjectName,
            EA:            ea,
            GA:            ga,
            Remaining:     isRemaining,
            Excluded:      excludedSubjects,
          })
          break
        case StateAbiOral:
          println("Oral", text)

          if strings.HasPrefix(text, "Erste Prüfungswoche") {
            // same dom-element <p> contains both weeks (not seperated in text = <br>)
            separationIndex := strings.Index(text, "Zweite Prüfungswoche")
            firstWeekText := strings.TrimSpace(text[:separationIndex])
            secondWeekText := strings.TrimSpace(text[separationIndex:])

            // - Erste/Zweite Prüfungswoche: [DDDD, DD. MMMM mit DDDD, DD. MMMM YYYY] (= Montag, 18. Mai mit Freitag, 22. Mai 2026)
            for i, weekText := range []string{firstWeekText, secondWeekText} {
              weekParts := strings.SplitN(weekText, ": ", 2)
              if len(weekParts) < 2 {
                // TODO
              }

              weekStartDate, weekEndDate, err := ConvertGermanWeekDateToIso(strings.Split(weekParts[1], " "))
              if err != nil {
                // TODO
              }

              oralWeeks = append(oralWeeks, api.ByOralExamDate{
                StartDate:     weekStartDate,
                EndDate:       weekEndDate,
                FormattedDate: weekParts[1],
                FormattedWeek: weekParts[0],
                WeekNumber:    i + 1,
              })
            }
          } else if strings.HasPrefix(text, "Die praktischen Prüfungen") {
            // 2026/2027:
            // - Die praktischen Prüfungen im Fach Sport werden nicht vor Montag, den 26. Januar 2026,
            //   die praktischen Prüfungen im Fach Musik nicht vor Montag, den 9. März 2026 durchgeführt.
            for {
              subjectStartIndex := strings.Index(text, "im Fach")
              if subjectStartIndex < 0 {
                break
              }

              text = text[subjectStartIndex+len("im Fach "):]
              println("Practical", text)
              subjectNameIndex := strings.Index(text, " ")
              if subjectNameIndex < 0 {
                // TODO
              }

              subjectName := text[:subjectNameIndex]
              dateStartIndex := strings.Index(text, "nicht vor")
              if dateStartIndex < 0 {
                // TODO
              }
              println("- Subject:", subjectName)

              dateTextAndRemaining := text[dateStartIndex+len("nicht vor "):]
              dateParts := strings.SplitN(dateTextAndRemaining, " ", 6) // DDDD, den DD., MMMM, YYYY, x
              if len(dateParts) < 6 {
                // TODO
              }
              println("- DateTextAndRemaining:", dateTextAndRemaining)

              usableDateParts := dateParts[2:]
              println("- UsableDateParts:", strings.Join(usableDateParts, " "))
              isoDate, err := ConvertGermanDateToIso(usableDateParts)
              if err != nil {
                // TODO
              }

              practicalDates = append(practicalDates, api.ByPracticalExamDate{
                StartDate:     isoDate,
                FormattedDate: strings.Join(dateParts[:5], " "),
                Subject:       subjectName,
              })

            }
          }

          break
        case StateAbiReport:
          println("Report", text)
          // - [...] findet am Freitag, den 26. Juni 2026 statt [...]
          dateStartIndex := strings.Index(text, "am ") + len("am ")
          dateTextAndRemaining := text[dateStartIndex:]
          dateParts := strings.SplitN(dateTextAndRemaining, " ", 6) // DDDD, den DD., MMMM, YYYY, x
          if len(dateParts) < 6 {
            // TODO
          }
          usableDateParts := dateParts[2:]
          isoDate, err := ConvertGermanDateToIso(usableDateParts)
          if err != nil {
            // TODO
          }
          graduationDate = &api.ByGraduationDate{
            Date:          isoDate,
            FormattedDate: strings.Join(dateParts[:5], " "),
          }
          break
        }
      })

      if state >= StateAbiInit && state <= StateAbiReport {
        for _, date := range writtenDates {
          fmt.Printf("Written Exam Date: %+v\n", date)
        }
        for _, date := range oralWeeks {
          fmt.Printf("Oral Exam Week: %+v\n", date)
        }
        for _, date := range practicalDates {
          fmt.Printf("Practical Exam Date: %+v\n", date)
        }

        println("State: ", state, " Year: ", year)
        println("-------------------------------------------------")
        gymAbiYearData[year] = api.ByAbiYearData{
          WrittenDates:   writtenDates,
          OralDates:      oralWeeks,
          PracticalDates: practicalDates,
          GraduationDate: *graduationDate,
        }
      }

    })

    scrapeData.GymAbiYearData = gymAbiYearData
  })

  c.OnScraped(func(response *colly.Response) {
    fmt.Printf("Scraped: %s\n", response.Request.URL)
  })

  err := c.Visit(URL)

  if err != nil {
    return nil, err
  }
  return scrapeData, nil
}
