package output

import (
	"fmt"
	"time"
)

// FormatRequestDate renders the request date string according to dateStyle and calendar options.
// Supported styles:
// - gce-verbose (default): Natural language sentence; calendar-sensitive
// - iso-date: "2025-09-09" (calendar-neutral; Gregorian)
// - rfc3339: RFC 3339 timestamp (UTC midnight) (Gregorian)
// - unix: Unix epoch seconds (UTC midnight) (Gregorian)
// - iso-week: ISO week date, e.g., "2025-W37-2" (Gregorian)
// Supported calendars for gce-verbose:
// - gregorian (default)
// - julian
// - buddhist (Thai solar)
// - minguo (ROC)
// - japanese (era-based)
// - islamic (Hijri, tabular civil)
// - hebrew (planned)
func FormatRequestDate(t time.Time) string {
	switch dateStyle {
	case "iso-date":
		return t.Format("2006-01-02")
	case "rfc3339":
		utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return utc.Format(time.RFC3339)
	case "unix":
		utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return fmt.Sprintf("%d", utc.Unix())
	case "iso-week":
		y, w := t.ISOWeek()
		isoWeekday := int(t.Weekday())
		if isoWeekday == 0 {
			isoWeekday = 7
		}
		return fmt.Sprintf("%04d-W%02d-%d", y, w, isoWeekday)
	case "gce-verbose":
		fallthrough
	default:
		return verboseCalendarSentence(t)
	}
}

func verboseCalendarSentence(t time.Time) string {
	weekday := t.Weekday().String()
	gregMonth := t.Month().String()
	day := t.Day()
	sfx := ordinalSuffix(day)
	y := t.Year()
	switch calendar {
	case "gregorian":
		return fmt.Sprintf("This request was made on %s, %s %d%s of the year %d of the common era.", weekday, gregMonth, day, sfx, y)
	case "buddhist":
		by := y + 543
		return fmt.Sprintf("This request was made on %s, %s %d%s of the year %d of the Buddhist Era.", weekday, gregMonth, day, sfx, by)
	case "minguo":
		ry := y - 1911
		return fmt.Sprintf("This request was made on %s, %s %d%s of the year %d of the Minguo calendar.", weekday, gregMonth, day, sfx, ry)
	case "julian":
		jy, jm, jd := julianFromGregorian(y, int(t.Month()), day)
		jMonth := monthName(jm)
		jsfx := ordinalSuffix(jd)
		return fmt.Sprintf("This request was made on %s, %s %d%s of the year %d of the Julian calendar.", weekday, jMonth, jd, jsfx, jy)
	case "japanese":
		era, eraYear := japaneseEra(y, int(t.Month()), day)
		return fmt.Sprintf("This request was made on %s, %s %d%s in %s %d of the Japanese calendar.", weekday, gregMonth, day, sfx, era, eraYear)
	case "islamic":
		iy, im, id := islamicCivilFromGregorian(y, int(t.Month()), day)
		iMonth := islamicMonthName(im)
		isfx := ordinalSuffix(id)
		return fmt.Sprintf("This request was made on %s, %s %d%s in year %d AH of the Islamic (Hijri) calendar.", weekday, iMonth, id, isfx, iy)
	case "hebrew":
		// Approximate (tabular) Hebrew year mapping: Hebrew year increments around Sep/Oct.
		// We use a coarse threshold of Sep 20 for increment; this avoids early-year misclassification.
		hy := hebrewYearApprox(y, int(t.Month()), day)
		return fmt.Sprintf("This request was made on %s, %s %d%s in year %d AM of the Hebrew calendar (tabular approximation).", weekday, gregMonth, day, sfx, hy)
	default:
		return fmt.Sprintf("This request was made on %s, %s %d%s of the year %d of the common era.", weekday, gregMonth, day, sfx, y)
	}
}

func ordinalSuffix(day int) string {
	switch day % 100 {
	case 11, 12, 13:
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func monthName(m int) string {
	switch m {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "December"
	default:
		return ""
	}
}

// Julian calendar conversion via JDN
func julianFromGregorian(y, m, d int) (jy, jm, jd int) {
	jdn := gregorianToJDN(y, m, d)
	jy, jm, jd = jdnToJulian(jdn)
	return
}

func gregorianToJDN(y, m, d int) int {
	a := (14 - m) / 12
	y2 := y + 4800 - a
	m2 := m + 12*a - 3
	jdn := d + (153*m2+2)/5 + 365*y2 + y2/4 - y2/100 + y2/400 - 32045
	return jdn
}

func jdnToJulian(jdn int) (year, month, day int) {
	c := jdn + 32082
	d := (4*c + 3) / 1461
	e := c - (1461*d)/4
	m := (5*e + 2) / 153
	day = e - (153*m+2)/5 + 1
	month = m + 3 - 12*(m/10)
	year = d - 4800 + m/10
	return
}

// Japanese era mapping
func japaneseEra(y, m, d int) (string, int) {
	// Define era boundaries (inclusive start dates)
	type eraDef struct {
		name    string
		y, m, d int
	}
	eras := []eraDef{
		{"Reiwa", 2019, 5, 1},
		{"Heisei", 1989, 1, 8},
		{"Showa", 1926, 12, 25},
		{"Taisho", 1912, 7, 30},
		{"Meiji", 1868, 10, 23},
	}
	for _, e := range eras {
		if afterOrEqual(y, m, d, e.y, e.m, e.d) {
			year := y - e.y + 1
			return e.name, year
		}
	}
	return "Pre-Meiji", y
}

func afterOrEqual(y, m, d, y2, m2, d2 int) bool {
	if y != y2 {
		return y > y2
	}
	if m != m2 {
		return m > m2
	}
	return d >= d2
}

// Islamic (Hijri) civil (tabular) conversion via JDN
func islamicCivilFromGregorian(y, m, d int) (iy, im, id int) {
	jdn := gregorianToJDN(y, m, d)
	l := jdn - 1948440 + 10632
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := ((10985 - l) / 5316) * ((50 * l) / 17719)
	j += (l / 5670) * ((43 * l) / 15238)
	l = l - ((30-j)/15)*((17719*j)/50) - (j/16)*((15238*j)/43) + 29
	im = (24 * l) / 709
	id = l - (709*im)/24
	iy = 30*n + j - 30
	return
}

func islamicMonthName(m int) string {
	switch m {
	case 1:
		return "Muharram"
	case 2:
		return "Safar"
	case 3:
		return "Rabi' al-awwal"
	case 4:
		return "Rabi' al-thani"
	case 5:
		return "Jumada al-awwal"
	case 6:
		return "Jumada al-thani"
	case 7:
		return "Rajab"
	case 8:
		return "Sha'ban"
	case 9:
		return "Ramadan"
	case 10:
		return "Shawwal"
	case 11:
		return "Dhu al-Qi'dah"
	case 12:
		return "Dhu al-Hijjah"
	default:
		return ""
	}
}

// hebrewYearApprox computes a tabular approximation of the Hebrew year for a given Gregorian date.
// Hebrew year = Gregorian year + 3760, increments near Rosh Hashanah (Sep/Oct). We use Sep 20 as threshold.
func hebrewYearApprox(gy, gm, gd int) int {
	if gm > 9 {
		return gy + 3761
	}
	if gm < 9 {
		return gy + 3760
	}
	// gm == 9 (September)
	if gd >= 20 {
		return gy + 3761
	}
	return gy + 3760
}
