package api

// Bavaria

type ByAbiYearData struct {
  WrittenDates   []ByWrittenExamDate   `json:"written"`
  OralDates      []ByOralExamDate      `json:"oral"`
  PracticalDates []ByPracticalExamDate `json:"practical"`
  GraduationDate ByGraduationDate      `json:"graduation"`
}

type ByWrittenExamDate struct {
  Date          string   `json:"date"`                        // ISO 8601 format: YYYY-MM-DD
  FormattedDate string   `json:"date_formatted"`              // Original: "DD. MMMM YYYY" ("15. März 2024")
  Subject       string   `json:"subject"`                     // Name ("Mathematik", "Deutsch", "Englisch", "Physik", ...)
  EA            bool     `json:"ea"`                          // erhöhtes Anforderungsniveau
  GA            bool     `json:"ga"`                          // grundlegendes Anforderungsniveau
  Remaining     bool     `json:"is_remaining,omitempty"`      // alle (weiteren) Abiturprüfungsfächer
  Excluded      []string `json:"excluded_subjects,omitempty"` // (mit Ausnahme von [SUBJECT_NAME,n], )
}

type ByOralExamDate struct {
  StartDate     string `json:"start_date"`     // ISO 8601 format: YYYY-MM-DD
  EndDate       string `json:"end_date"`       // ISO 8601 format: YYYY-MM-DD
  FormattedDate string `json:"date_formatted"` // Original: "DDDD, DD. MMMM mit DDDD, DD. MMMM YYYY" ("Montag, 18. Mai mit Freitag, 22. Mai 2026")
  FormattedWeek string `json:"week_formatted"` // Original: "Erste Prüfungswoche" / "Zweite Prüfungswoche"
  WeekNumber    int    `json:"week_number"`    // 1 / 2
}

type ByPracticalExamDate struct {
  StartDate     string `json:"start_date"`     // ISO 8601 format: YYYY-MM-DD
  FormattedDate string `json:"date_formatted"` // Original: "DDDD, DD. MMMM YYYY" ("Montag, den 26. Januar 2026")
  Subject       string `json:"subject"`        // Name ("Sport", "Musik")
}

type ByGraduationDate struct {
  Date          string `json:"date"`           // ISO 8601 format: YYYY-MM-DD
  FormattedDate string `json:"date_formatted"` // Original: "DDDD, den DD. MMMM YYYY" ("Freitag, den 26. Juni 2026")
}
