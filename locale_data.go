package chronogo

// registerDefaultLocales registers all supported locales
func registerDefaultLocales() {
	RegisterLocale(createEnUSLocale())
	RegisterLocale(createEsESLocale())
	RegisterLocale(createFrFRLocale())
	RegisterLocale(createDeDELocale())
	RegisterLocale(createZhHansLocale())
	RegisterLocale(createPtBRLocale())
}

// createEnUSLocale creates the English (United States) locale
func createEnUSLocale() *Locale {
	return &Locale{
		Code: "en-US",
		Name: "English (United States)",
		MonthNames: []string{
			"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December",
		},
		MonthAbbr: []string{
			"Jan", "Feb", "Mar", "Apr", "May", "Jun",
			"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		},
		WeekdayNames: []string{
			"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
		},
		WeekdayAbbr: []string{
			"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat",
		},
		AMPMNames: []string{"AM", "PM"},
		Ordinals:  createEnglishOrdinals(),
		TimeUnits: map[string]TimeUnitNames{
			"second":   {Singular: "second", Plural: "seconds"},
			"minute":   {Singular: "minute", Plural: "minutes"},
			"hour":     {Singular: "hour", Plural: "hours"},
			"day":      {Singular: "day", Plural: "days"},
			"week":     {Singular: "week", Plural: "weeks"},
			"month":    {Singular: "month", Plural: "months"},
			"year":     {Singular: "year", Plural: "years"},
			"moments":  {Singular: "a few seconds ago", Plural: "in a few seconds"},
			"patterns": {Singular: "%d %s ago", Plural: "in %d %s"},
		},
		DateFormats: map[string]string{
			"short":  "1/2/2006",
			"medium": "Jan 2, 2006",
			"long":   "January 2, 2006",
			"full":   "Monday, January 2, 2006",
		},
	}
}

// createEsESLocale creates the Spanish (Spain) locale
func createEsESLocale() *Locale {
	return &Locale{
		Code: "es-ES",
		Name: "Español (España)",
		MonthNames: []string{
			"enero", "febrero", "marzo", "abril", "mayo", "junio",
			"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
		},
		MonthAbbr: []string{
			"ene", "feb", "mar", "abr", "may", "jun",
			"jul", "ago", "sep", "oct", "nov", "dic",
		},
		WeekdayNames: []string{
			"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado",
		},
		WeekdayAbbr: []string{
			"dom", "lun", "mar", "mié", "jue", "vie", "sáb",
		},
		AMPMNames: []string{"AM", "PM"},
		Ordinals:  createSpanishOrdinals(),
		TimeUnits: map[string]TimeUnitNames{
			"second":   {Singular: "segundo", Plural: "segundos"},
			"minute":   {Singular: "minuto", Plural: "minutos"},
			"hour":     {Singular: "hora", Plural: "horas"},
			"day":      {Singular: "día", Plural: "días"},
			"week":     {Singular: "semana", Plural: "semanas"},
			"month":    {Singular: "mes", Plural: "meses"},
			"year":     {Singular: "año", Plural: "años"},
			"moments":  {Singular: "hace unos momentos", Plural: "en unos momentos"},
			"patterns": {Singular: "hace %d %s", Plural: "en %d %s"},
		},
		DateFormats: map[string]string{
			"short":  "2/1/2006",
			"medium": "2 ene 2006",
			"long":   "2 de enero de 2006",
			"full":   "lunes, 2 de enero de 2006",
		},
	}
}

// createFrFRLocale creates the French (France) locale
func createFrFRLocale() *Locale {
	return &Locale{
		Code: "fr-FR",
		Name: "Français (France)",
		MonthNames: []string{
			"janvier", "février", "mars", "avril", "mai", "juin",
			"juillet", "août", "septembre", "octobre", "novembre", "décembre",
		},
		MonthAbbr: []string{
			"janv", "févr", "mars", "avr", "mai", "juin",
			"juil", "août", "sept", "oct", "nov", "déc",
		},
		WeekdayNames: []string{
			"dimanche", "lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi",
		},
		WeekdayAbbr: []string{
			"dim", "lun", "mar", "mer", "jeu", "ven", "sam",
		},
		AMPMNames: []string{"AM", "PM"},
		Ordinals:  createFrenchOrdinals(),
		TimeUnits: map[string]TimeUnitNames{
			"second":   {Singular: "seconde", Plural: "secondes"},
			"minute":   {Singular: "minute", Plural: "minutes"},
			"hour":     {Singular: "heure", Plural: "heures"},
			"day":      {Singular: "jour", Plural: "jours"},
			"week":     {Singular: "semaine", Plural: "semaines"},
			"month":    {Singular: "mois", Plural: "mois"},
			"year":     {Singular: "an", Plural: "ans"},
			"moments":  {Singular: "il y a quelques instants", Plural: "dans quelques instants"},
			"patterns": {Singular: "il y a %d %s", Plural: "dans %d %s"},
		},
		DateFormats: map[string]string{
			"short":  "02/01/2006",
			"medium": "2 janv 2006",
			"long":   "2 janvier 2006",
			"full":   "lundi 2 janvier 2006",
		},
	}
}

// createDeDELocale creates the German (Germany) locale
func createDeDELocale() *Locale {
	return &Locale{
		Code: "de-DE",
		Name: "Deutsch (Deutschland)",
		MonthNames: []string{
			"Januar", "Februar", "März", "April", "Mai", "Juni",
			"Juli", "August", "September", "Oktober", "November", "Dezember",
		},
		MonthAbbr: []string{
			"Jan", "Feb", "Mär", "Apr", "Mai", "Jun",
			"Jul", "Aug", "Sep", "Okt", "Nov", "Dez",
		},
		WeekdayNames: []string{
			"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag",
		},
		WeekdayAbbr: []string{
			"So", "Mo", "Di", "Mi", "Do", "Fr", "Sa",
		},
		AMPMNames: []string{"AM", "PM"},
		Ordinals:  createGermanOrdinals(),
		TimeUnits: map[string]TimeUnitNames{
			"second":   {Singular: "Sekunde", Plural: "Sekunden"},
			"minute":   {Singular: "Minute", Plural: "Minuten"},
			"hour":     {Singular: "Stunde", Plural: "Stunden"},
			"day":      {Singular: "Tag", Plural: "Tage"},
			"week":     {Singular: "Woche", Plural: "Wochen"},
			"month":    {Singular: "Monat", Plural: "Monate"},
			"year":     {Singular: "Jahr", Plural: "Jahre"},
			"moments":  {Singular: "vor wenigen Augenblicken", Plural: "in wenigen Augenblicken"},
			"patterns": {Singular: "vor %d %s", Plural: "in %d %s"},
		},
		DateFormats: map[string]string{
			"short":  "2.1.2006",
			"medium": "2. Jan 2006",
			"long":   "2. Januar 2006",
			"full":   "Montag, 2. Januar 2006",
		},
	}
}

// createZhHansLocale creates the Chinese (Simplified) locale
func createZhHansLocale() *Locale {
	return &Locale{
		Code: "zh-Hans",
		Name: "中文 (简体)",
		MonthNames: []string{
			"一月", "二月", "三月", "四月", "五月", "六月",
			"七月", "八月", "九月", "十月", "十一月", "十二月",
		},
		MonthAbbr: []string{
			"1月", "2月", "3月", "4月", "5月", "6月",
			"7月", "8月", "9月", "10月", "11月", "12月",
		},
		WeekdayNames: []string{
			"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六",
		},
		WeekdayAbbr: []string{
			"周日", "周一", "周二", "周三", "周四", "周五", "周六",
		},
		AMPMNames: []string{"上午", "下午"},
		Ordinals:  createChineseOrdinals(),
		TimeUnits: map[string]TimeUnitNames{
			"second":   {Singular: "秒", Plural: "秒"},
			"minute":   {Singular: "分钟", Plural: "分钟"},
			"hour":     {Singular: "小时", Plural: "小时"},
			"day":      {Singular: "天", Plural: "天"},
			"week":     {Singular: "周", Plural: "周"},
			"month":    {Singular: "个月", Plural: "个月"},
			"year":     {Singular: "年", Plural: "年"},
			"moments":  {Singular: "刚刚", Plural: "马上"},
			"patterns": {Singular: "%d%s前", Plural: "%d%s后"},
		},
		DateFormats: map[string]string{
			"short":  "2006/1/2",
			"medium": "2006年1月2日",
			"long":   "2006年1月2日",
			"full":   "2006年1月2日星期一",
		},
	}
}

// createPtBRLocale creates the Portuguese (Brazil) locale
func createPtBRLocale() *Locale {
	return &Locale{
		Code: "pt-BR",
		Name: "Português (Brasil)",
		MonthNames: []string{
			"janeiro", "fevereiro", "março", "abril", "maio", "junho",
			"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
		},
		MonthAbbr: []string{
			"jan", "fev", "mar", "abr", "mai", "jun",
			"jul", "ago", "set", "out", "nov", "dez",
		},
		WeekdayNames: []string{
			"domingo", "segunda-feira", "terça-feira", "quarta-feira", "quinta-feira", "sexta-feira", "sábado",
		},
		WeekdayAbbr: []string{
			"dom", "seg", "ter", "qua", "qui", "sex", "sáb",
		},
		AMPMNames: []string{"AM", "PM"},
		Ordinals:  createPortugueseOrdinals(),
		TimeUnits: map[string]TimeUnitNames{
			"second":   {Singular: "segundo", Plural: "segundos"},
			"minute":   {Singular: "minuto", Plural: "minutos"},
			"hour":     {Singular: "hora", Plural: "horas"},
			"day":      {Singular: "dia", Plural: "dias"},
			"week":     {Singular: "semana", Plural: "semanas"},
			"month":    {Singular: "mês", Plural: "meses"},
			"year":     {Singular: "ano", Plural: "anos"},
			"moments":  {Singular: "há poucos instantes", Plural: "em poucos instantes"},
			"patterns": {Singular: "há %d %s", Plural: "em %d %s"},
		},
		DateFormats: map[string]string{
			"short":  "02/01/2006",
			"medium": "2 de jan de 2006",
			"long":   "2 de janeiro de 2006",
			"full":   "segunda-feira, 2 de janeiro de 2006",
		},
	}
}

// Ordinal creation functions

func createEnglishOrdinals() map[int]string {
	ordinals := make(map[int]string)

	// Special cases for 11th, 12th, 13th
	for i := 11; i <= 13; i++ {
		ordinals[i] = "th"
	}

	// Generate ordinals for 1-31 (days of month)
	for i := 1; i <= 31; i++ {
		if _, exists := ordinals[i]; exists {
			continue // Skip already set special cases
		}

		switch i % 10 {
		case 1:
			ordinals[i] = "st"
		case 2:
			ordinals[i] = "nd"
		case 3:
			ordinals[i] = "rd"
		default:
			ordinals[i] = "th"
		}
	}

	return ordinals
}

func createSpanishOrdinals() map[int]string {
	ordinals := make(map[int]string)

	// Spanish ordinals are more complex, but for dates we typically use cardinal numbers
	// For simplicity, we'll use "º" for masculine ordinals
	for i := 1; i <= 31; i++ {
		ordinals[i] = "º"
	}

	return ordinals
}

func createFrenchOrdinals() map[int]string {
	ordinals := make(map[int]string)

	// French uses "er" for 1st, "e" for others
	ordinals[1] = "er"
	for i := 2; i <= 31; i++ {
		ordinals[i] = "e"
	}

	return ordinals
}

func createGermanOrdinals() map[int]string {
	ordinals := make(map[int]string)

	// German ordinals use "." after the number
	for i := 1; i <= 31; i++ {
		ordinals[i] = "."
	}

	return ordinals
}

func createChineseOrdinals() map[int]string {
	ordinals := make(map[int]string)

	// Chinese doesn't typically use ordinal suffixes for dates
	for i := 1; i <= 31; i++ {
		ordinals[i] = "日" // "day" character
	}

	return ordinals
}

func createPortugueseOrdinals() map[int]string {
	ordinals := make(map[int]string)

	// Portuguese ordinals
	ordinals[1] = "º"
	for i := 2; i <= 31; i++ {
		ordinals[i] = "º"
	}

	return ordinals
}
