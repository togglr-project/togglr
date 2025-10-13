alter table algorithms
    add column kind varchar(50) not null default 'bandit';

comment on column algorithms.kind is
    'Algorithm family or category (e.g., bandit, rule-based, reinforcement, ml-model).';

update algorithms set kind = 'bandit' where slug in ('epsilon-greedy', 'thompson-sampling', 'ucb');
