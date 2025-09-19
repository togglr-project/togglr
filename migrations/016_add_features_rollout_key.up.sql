alter table features add column rollout_key varchar(50);

create index if not exists idx_features_rollout_key on features(rollout_key) where rollout_key is not null;
