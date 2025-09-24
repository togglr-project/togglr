create or replace view v_project_health as
with feature_with_category as (
    select
        f.id as feature_id,
        f.project_id,
        c.id as category_id,
        c.name as category_name,
        c.slug as category_slug,
        f.enabled,
        exists (
            select 1
            from feature_tags ft
                     join tags t on ft.tag_id = t.id
                     join categories c2 on t.category_id = c2.id
            where ft.feature_id = f.id
              and c2.slug = 'auto-disable'
        ) as is_auto_disabled,
        exists (
            select 1
            from feature_tags ft
                     join tags t on ft.tag_id = t.id
                     join categories c2 on t.category_id = c2.id
            where ft.feature_id = f.id
              and c2.slug = 'guarded'
        ) as is_guarded
    from features f
             left join feature_tags ft on ft.feature_id = f.id
             left join tags t on ft.tag_id = t.id
             left join categories c on t.category_id = c.id
),
     pending as (
         select distinct pce.entity_id as feature_id
         from pending_change_entities pce
                  join pending_changes pc on pc.id = pce.pending_change_id
         where pc.status = 'pending'
           and pce.entity = 'feature'
     )
select
    fwc.project_id,
    count(distinct fwc.feature_id) as total_features,
    count(distinct case when fwc.enabled then fwc.feature_id end) as enabled_features,
    count(distinct case when not fwc.enabled then fwc.feature_id end) as disabled_features,
    count(distinct case when fwc.is_auto_disabled then fwc.feature_id end) as auto_disabled_features,
    count(distinct case when p.feature_id is not null then fwc.feature_id end) as pending_features
from feature_with_category fwc
         left join pending p on p.feature_id = fwc.feature_id
group by fwc.project_id;

---

create or replace view v_project_category_health as
with feature_with_category as (
    select
        f.id as feature_id,
        f.project_id,
        c.id as category_id,
        c.name as category_name,
        c.slug as category_slug,
        f.enabled,
        exists (
            select 1
            from feature_tags ft
                     join tags t on ft.tag_id = t.id
                     join categories c2 on t.category_id = c2.id
            where ft.feature_id = f.id
              and c2.slug = 'auto-disable'
        ) as is_auto_disabled,
        exists (
            select 1
            from feature_tags ft
                     join tags t on ft.tag_id = t.id
                     join categories c2 on t.category_id = c2.id
            where ft.feature_id = f.id
              and c2.slug = 'guarded'
        ) as is_guarded
    from features f
             left join feature_tags ft on ft.feature_id = f.id
             left join tags t on ft.tag_id = t.id
             left join categories c on t.category_id = c.id
),
     pending as (
         select distinct pce.entity_id as feature_id
         from pending_change_entities pce
                  join pending_changes pc on pc.id = pce.pending_change_id
         where pc.status = 'pending'
           and pce.entity = 'feature'
     )
select
    fwc.project_id,
    fwc.category_id,
    fwc.category_name,
    fwc.category_slug,
    count(distinct fwc.feature_id) as total_features,
    count(distinct case when fwc.enabled then fwc.feature_id end) as enabled_features,
    count(distinct case when not fwc.enabled then fwc.feature_id end) as disabled_features,
    count(distinct case when fwc.is_auto_disabled then fwc.feature_id end) as auto_disabled_features,
    count(distinct case when p.feature_id is not null then fwc.feature_id end) as pending_features
from feature_with_category fwc
         left join pending p on p.feature_id = fwc.feature_id
group by fwc.project_id, fwc.category_id, fwc.category_name, fwc.category_slug;
