create type rule_action as enum ('assign', 'include', 'exclude');

alter table rules add column action rule_action not null default 'assign';

alter table rules alter column flag_variant_id drop not null;

alter table rules
    add constraint rules_flag_variant_check
        check (
            (action = 'assign' and flag_variant_id is not null)
                or (action in ('include','exclude') and flag_variant_id is null)
            );
