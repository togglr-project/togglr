alter table rules drop constraint rules_flag_variant_check;

alter table rules alter column flag_variant_id set not null;

alter table rules drop column action;

drop type rule_action;
