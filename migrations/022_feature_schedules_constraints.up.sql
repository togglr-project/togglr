-- extension is needed for EXCLUDE CONSTRAINT
create extension if not exists btree_gist;

-- uniqueness of cron for each feature_id
create unique index if not exists feature_schedules_uniq_cron_guard
    on feature_schedules (feature_id)
    where cron_expr is not null;

-- prevent overlapping one-shot schedules for one feature_id
alter table feature_schedules
    add constraint feature_schedules_no_overlap_guard
        exclude using gist (
            feature_id with =,
            tstzrange(starts_at, ends_at, '[]') with &&
        )
        where (cron_expr is null);

-- trigger-function: prevents mixing cron and one-shot
create or replace function check_feature_schedule_mode()
    returns trigger as $$
begin
    if NEW.cron_expr is not null then
        -- if adding cron, check that there is no one-shot for this feature
        if exists(
            select 1 from feature_schedules
            where feature_id = NEW.feature_id
              and cron_expr is null
              and id <> NEW.id
        ) then
            raise exception 'Feature % already has one-shot schedules, cannot add cron', NEW.feature_id;
        end if;
    else
        -- if adding one-shot, check that there is no cron for this feature
        if exists(
            select 1 from feature_schedules
            where feature_id = NEW.feature_id
              and cron_expr is not null
              and id <> NEW.id
        ) then
            raise exception 'Feature % already has a cron schedule, cannot add one-shot', NEW.feature_id;
        end if;
    end if;

    return NEW;
end;
$$ language plpgsql;

-- the trigger
drop trigger if exists trg_check_feature_schedule_mode on feature_schedules;

create trigger trg_check_feature_schedule_mode
    before insert or update on feature_schedules
    for each row
execute function check_feature_schedule_mode();
