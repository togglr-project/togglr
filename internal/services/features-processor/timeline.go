package featuresprocessor

import (
	"sort"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/togglr-project/togglr/internal/domain"
)

const cronMaxIters = 11_000

func (*Service) BuildFeatureTimeline(
	feature domain.FeatureExtended,
	from time.Time,
	to time.Time,
) ([]domain.TimelineEvent, error) {
	// 1. Master Enable = OFF → фича полностью выключена
	if !feature.Enabled {
		return []domain.TimelineEvent{
			{Time: from, Enabled: false},
			{Time: to, Enabled: false},
		}, nil
	}

	// 2. Если нет расписаний, фича остается в ручном состоянии
	if len(feature.Schedules) == 0 {
		return []domain.TimelineEvent{
			{Time: from, Enabled: feature.Enabled},
			{Time: to, Enabled: feature.Enabled},
		}, nil
	}

	// 3. Определяем baseline состояние
	featurePrepared := MakeFeaturePrepared(feature)
	baseline := getScheduleBaseline(feature.Schedules)

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

	// 5. Для каждого расписания собираем переключения
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

		// Для cron расписаний: проверяем, нужно ли добавить событие перехода к baseline
		// в начале интервала, если cron активен, но мы находимся в промежутке между срабатываниями
		if sched.CronExpr != nil && *sched.CronExpr != "" && sched.CronDuration != nil && *sched.CronDuration > 0 {
			// Проверяем, активен ли cron в начале интервала
			if active, _ := IsScheduleActive(sched, featurePrepared.crons, from, feature.CreatedAt); active {
				// Находим время окончания текущего активного периода
				scheduleStart := feature.CreatedAt
				if sched.StartsAt != nil {
					scheduleStart = *sched.StartsAt
				}

				prevCron, found := findPrevCron(featurePrepared.crons[sched.ID], from, scheduleStart)
				if found {
					endOfActivePeriod := prevCron.Add(*sched.CronDuration)
					// Если конец активного периода находится в нашем интервале, добавляем событие
					if endOfActivePeriod.After(from) && endOfActivePeriod.Before(to) {
						events = append(events, domain.TimelineEvent{
							Time:    endOfActivePeriod,
							Enabled: baseline, // переход к baseline состоянию
						})
					}
				}
			}
		}

		// Если есть cron — разворачиваем его с учетом starts_at и ends_at
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

			// Определяем начало работы cron-расписания
			cronStart := feature.CreatedAt
			if sched.StartsAt != nil {
				cronStart = *sched.StartsAt
			}

			cronStart = cronStart.In(loc)

			// Определяем конец работы cron-расписания
			cronEnd := to.In(loc)
			if sched.EndsAt != nil {
				cronEnd = sched.EndsAt.In(loc)
			}

			// Cron работает только в пределах своего временного окна
			effectiveFrom := from.In(loc)
			effectiveTo := to.In(loc)

			// Ограничиваем интервал расписанием
			if cronStart.After(effectiveFrom) {
				effectiveFrom = cronStart
			}

			if cronEnd.Before(effectiveTo) {
				effectiveTo = cronEnd
			}

			// Если расписание не пересекается с запрашиваемым интервалом, пропускаем
			if effectiveFrom.After(effectiveTo) || effectiveFrom.Equal(effectiveTo) {
				continue
			}

			// Чтобы включать событие ровно в момент effectiveFrom (включительно),
			// стартуем с на 1нс раньше.
			cursor := effectiveFrom.Add(-time.Nanosecond)

			cnt := 0

			for {
				next := schedule.Next(cursor)
				// effectiveTo — исключительно: отбрасываем next >= effectiveTo
				if !next.Before(effectiveTo) {
					break
				}

				// Добавляем событие включения/выключения
				events = append(events, domain.TimelineEvent{
					Time:    next,
					Enabled: sched.Action == domain.FeatureScheduleActionEnable,
				})

				// Если задана продолжительность для cron, добавляем событие возврата к baseline
				if sched.CronDuration != nil && *sched.CronDuration > 0 {
					endTime := next.Add(*sched.CronDuration)
					// Проверяем, что конец не выходит за границы расписания и запрашиваемого интервала
					if endTime.Before(effectiveTo) && endTime.Before(cronEnd) {
						events = append(events, domain.TimelineEvent{
							Time:    endTime,
							Enabled: baseline, // возвращаемся к baseline состоянию
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
