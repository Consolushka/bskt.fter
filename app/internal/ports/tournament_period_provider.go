package ports

import "time"

// TournamentPeriodProvider — опциональная способность провайдера данных:
// сообщить плановые даты начала и конца сезона турнира.
//
// Реализуют только те провайдеры, у которых эти даты есть в API напрямую
// (например, API-Basketball отдаёт их в /leagues). Для остальных вызывающий код
// проверяет способность приведением типа и, если она не реализована, идёт по
// запасному пути (явные даты от пользователя или предупреждение).
type TournamentPeriodProvider interface {
	GetTournamentPeriod(season string) (start, end time.Time, err error)
}
