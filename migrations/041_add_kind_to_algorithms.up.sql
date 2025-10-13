alter table algorithms
    add column kind varchar(50) not null default 'bandit';

comment on column algorithms.kind is
    'Algorithm family or category (e.g., bandit, rule-based, reinforcement, ml-model).';

update algorithms set kind = 'bandit' where slug in ('epsilon-greedy', 'thompson-sampling', 'ucb');

---

create or replace function check_feature_algorithms_kind_consistency()
    returns trigger as $$
declare
    existing_kind varchar(50);
    new_kind varchar(50);
begin
    select kind into new_kind from algorithms where id = new.algorithm_id;

    select distinct a.kind
    into existing_kind
    from feature_algorithms fa
             join algorithms a on a.id = fa.algorithm_id
    where fa.feature_id = new.feature_id
    limit 1;

    if existing_kind is not null and existing_kind <> new_kind then
        raise exception
            'Algorithm kind mismatch: existing feature algorithms are of kind %, but new algorithm is %',
            existing_kind, new_kind;
    end if;

    return new;
end;
$$ language plpgsql;

create trigger trg_check_feature_algorithms_kind
    before insert or update on feature_algorithms
    for each row execute function check_feature_algorithms_kind_consistency();
