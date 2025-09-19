package featuresprocessor

import (
	"sort"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/rom8726/etoggle/internal/domain"
)

const cronMaxIters = 11_000

func (s *Service) BuildFeatureTimeline(
	feature domain.FeatureExtended,
	from time.Time,
	to time.Time,
) ([]domain.TimelineEvent, error) {
	featurePrepared := MakeFeaturePrepared(feature)

	events := []domain.TimelineEvent{
		{
			Time:    from,
			Enabled: IsFeatureActiveNow(featurePrepared, from),
		},
		{
			Time:    to,
			Enabled: IsFeatureActiveNow(featurePrepared, to),
		},
	}

	// 1. Базовое состояние — если нет расписаний, то считаем от created_at
	if len(feature.Schedules) == 0 {
		if feature.Enabled {
			events = append(events, domain.TimelineEvent{
				Time:    feature.CreatedAt,
				Enabled: true,
			})
		}

		return events, nil
	}

	// 2. Для каждого расписания собираем переключения
	for _, sched := range feature.Schedules {
		// Считаем начальное состояние
		if sched.StartsAt != nil && sched.StartsAt.After(from) && sched.StartsAt.Before(to) {
			events = append(events, domain.TimelineEvent{
				Time:    *sched.StartsAt,
				Enabled: sched.Action == domain.FeatureScheduleActionEnable,
			})
		}

		if sched.EndsAt != nil && sched.EndsAt.After(from) && sched.EndsAt.Before(to) {
			events = append(events, domain.TimelineEvent{
				Time:    *sched.EndsAt,
				Enabled: sched.Action != domain.FeatureScheduleActionEnable,
			})
		}

		// Если есть cron — разворачиваем его на весь интервал
		if sched.CronExpr != nil && *sched.CronExpr != "" {
			parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
			schedule, err := parser.Parse(*sched.CronExpr)
			if err != nil {
				continue
			}

			// Учитываем таймзону расписания
			loc := from.Location()
			if sched.Timezone != "" {
				if tz, err := time.LoadLocation(sched.Timezone); err == nil {
					loc = tz
				}
			}

			// Чтобы включать событие ровно в момент 'from' (from — включительно),
			// стартуем с на 1нс раньше.
			cursor := from.In(loc).Add(-time.Nanosecond)

			cnt := 0
			for {
				next := schedule.Next(cursor)
				// to — исключительно: отбрасываем next >= to
				if !next.Before(to.In(loc)) {
					break
				}

				// Добавляем событие включения/выключения
				events = append(events, domain.TimelineEvent{
					Time:    next,
					Enabled: sched.Action == domain.FeatureScheduleActionEnable,
				})

				// Если задана продолжительность для cron, добавляем обратное событие
				if sched.CronDuration != nil && *sched.CronDuration > 0 {
					endTime := next.Add(*sched.CronDuration)
					// Проверяем, что конец не выходит за границы запрашиваемого интервала
					if endTime.Before(to.In(loc)) {
						events = append(events, domain.TimelineEvent{
							Time:    endTime,
							Enabled: sched.Action != domain.FeatureScheduleActionEnable,
						})
					}
				}

				cursor = next

				cnt++
				if cnt > cronMaxIters {
					break
				}
			}
		}
	}

	// 3. Отсортировать события по времени
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	return events, nil
}
