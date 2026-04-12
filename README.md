# kmapi

###### Unofficial REST API for German Abitur exam dates. Scrapes official state ministry website, built in Go. <br>

**Disclaimer**: This project is not affiliated with any official government entity and is intended for educational and informational purposes only.

> **Motivation**: Im Rahmen der Entwicklung meiner **[G9 Notenapp](https://github.com/anweisen/G9)** für das neue G9 Abitur an bayerischen Gymnasien,
> benötigte ich eine zuverlässige Quelle für die aktuellen Abiturprüfungstermine. Da es keine offizielle API gibt, habe ich diese inoffizielle API entwickelt,
> um die Daten automatisiert zu sammeln und bereitzustellen.
> Daher ist die API derzeit auf bayerische Gymnasien beschränkt, jedoch so konzipiert, dass sie leicht auf andere Schularten und Bundesländer erweitert werden kann. <br>
> **Name**: *KM* ist die Abkürzung für *Kultusministerium*, dessen Website diese API scraped, daher *kmapi*.

**Demo: [kmapi.anweisen.net](https://kmapi.anweisen.net)**

This project is written in Golang and uses Redis for caching to reduce the scraping load.

Because the API scrapes data from the official state ministry website, it will be subject to changes in the website's structure,
which will temporarily affect the API's functionality, requiring updates to the scraping logic.

## API Endpoints
1. [Bavaria (Bayern)](#bavaria-bayern)
   1. [Abitur (Gymnasium)](#abitur-gymnasium-bayern)

`🚧` *Currently, only Bavaria is supported, but the API is designed to be extended to other states in the future.*

All endpoints return data in JSON format. <br>
All endpoint paths are structured as: ``/{state}/{school_type}/{exam_type}/{year}``, where:
- `{state}`: The state abbreviation (e.g., `by` for Bavaria).
- `{school_type}`: The type of school (e.g., `gym` for Gymnasium).
- `{exam_type}`: The type of exam (e.g., `abi` for Abitur).
- `{year}`: The year for which the exam dates are requested (e.g., `2026`).

### Bavaria (Bayern)

Source: *[www.km.bayern.de/termine/pruefungen-und-zeugnisse](https://www.km.bayern.de/termine/pruefungen-und-zeugnisse)*

### Abitur (Gymnasium, Bayern)
**GET** ``/by/gym/abi/{year}``

| Field      | Type                                     | Description |
|------------|------------------------------------------|-------------|
| written    | array of **written exam date** objects   | Enthält     |
| oral       | array of **oral exam week** objects      |             |
| practical  | array of **practical exam date** objects |             |
| graduation | **graduation date** object               |             |


**Written Exam Date Object**

| Field              | Type      | Description                                                                                                                                                                                                                                                              |
|--------------------|-----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| subject            | string    | Name des Abiturprüfungsfachs (z.B. "Mathematik") oder "alle (weiteren) Prüfungsfächer" für `is_remaining=true`                                                                                                                                                           |
| date               | string    | Datum der Abiturprüfung in ISO-8601 (YYYY-MM-DD) (z.B. "2026-05-20").                                                                                                                                                                                                    |
| date_formatted     | string    | Original formatiertes Datum (z.B. "23. April 2026")                                                                                                                                                                                                                      |
| ea                 | boolean   | Prüfung für Abiturprüfungsfach auf erhöhtem Anforderungsniveau                                                                                                                                                                                                           |
| ga                 | boolean   | Prüfung für Abiturprüfungsfach auf grundlegendem Anforderungsniveau                                                                                                                                                                                                      |
| is_remaining?      | boolean?  | Falls, `true` werden an diesem Tag alle weiteren Fächer auf dem gegebenen Anforderungsniveau geprüft                                                                                                                                                                     |
| excluded_subjects? | []string? | Array mit Namen von an diesem Termin ausgenommenen Prüfungsfächer. Handelt es sich bei diesem Termin um "alle Fächer" (`is_remaining=true`) werden u.U. bestimmte Fächer explizit ausgenommen, da sie *zu einem späteren Zeitpunkt* einen separaten Prüfungstermin haben |

**Oral Exam Week Object**

| Field          | Type   | Description                                                                          |
|----------------|--------|--------------------------------------------------------------------------------------|
| week_number    | int    | Nummer der mündlichen Prüfungswoche (1 oder 2)                                       |
| start_date     | string | Startdatum der mündlichen Prüfungswoche in ISO-8601 (YYYY-MM-DD) (z.B. "2026-06-15") |
| end_date       | string | Enddatum der mündlichen Prüfungswoche in ISO-8601 (YYYY-MM-DD) (z.B. "2026-06-19")   |
| date_formatted | string | Original formatiertes Datum (z.B. "Montag, 18. Mai mit Freitag, 22. Mai 2026")       |
| week_formatted | string | Original formatierter Name der mündlichen Prüfungswoche (z.B. "Erste Prüfungswoche") |

**Practical Exam Date Object**

| Field          | Type   | Description                                                                                                          |
|----------------|--------|----------------------------------------------------------------------------------------------------------------------|
| subject        | string | Name des praktischen Abiturprüfungsfachs (z.B. "Sport", "Musik")                                                     |
| start_date     | string | Startdatum in ISO-8601 (YYYY-MM-DD) der praktischen Abiturprüfungen in diesem Fach (z.B. "2026-06-01")               |
| date_formatted | string | Original formatiertes Startdatum der praktischen Abiturprüfungen in diesem Fach (z.B. "Montag, den 26. Januar 2026") |

**Graduation Date Object**

| Field          | Type   | Description                                                                                                     |
|----------------|--------|-----------------------------------------------------------------------------------------------------------------|
| date           | string | Datum in ISO-8601 (YYYY-MM-DD) der Entlassung der Abiturienten und Abiturzeugnisausstellung (z.B. "2026-07-15") |
| date_formatted | string | Original formatiertes Datum (z.B. "Freitag, den 26. Juni 2026")                                                 |

