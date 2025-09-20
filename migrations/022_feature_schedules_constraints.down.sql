drop trigger if exists trg_check_feature_schedule_mode on feature_schedules;

drop function if exists check_feature_schedule_mode;

alter table feature_schedules
    drop constraint feature_schedules_no_overlap_guard;

drop index if exists feature_schedules_uniq_cron_guard;

drop extension if exists btree_gist;
