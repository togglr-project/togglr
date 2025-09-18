package featuresprocessor

import (
	"sort"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/rom8726/etoggle/internal/domain"
)

const cronMaxIters = 11_000

func BuildFeatureTimeline(
	feature domain.FeatureExtended,
	from time.Time,
	to time.Time,
) ([]domain.TimelineEvent, error) {
	var events []domain.TimelineEvent

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

			cursor := from
			cnt := 0
			for cursor.Before(to) {
				next := schedule.Next(cursor)
				if next.After(to) {
					break
				}

				events = append(events, domain.TimelineEvent{
					Time:    next,
					Enabled: sched.Action == domain.FeatureScheduleActionEnable,
				})

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
