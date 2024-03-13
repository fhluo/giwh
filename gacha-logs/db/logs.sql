create table if not exists log (
    id text primary key not null,
    uid text,
    gacha_type text,
    name text,
    item_id text,
    item_type text,
    rank_type text,
    count text,
    time text,
    lang text
);

-- 索引
create index if not exists logs_index on log (uid, gacha_type, rank_type);

-- UID 列表
create view if not exists uid as
select distinct uid
from log;

-- 编号
create view if not exists number as
select id as log_id,
    row_number() over (
        partition by uid,
        case
            when gacha_type in (301, 400) then 301
            else gacha_type
        end
        order by id
    ) as number
from log;

-- 4/5 星 pity
create view if not exists pity as
with last_number as (
    select id as log_id,
        lag (number) over (
            partition by uid,
            case
                when gacha_type in (301, 400) then 301
                else gacha_type
            end,
            rank_type
        ) as last_number
    from log
        join number on log.id = number.log_id
    where rank_type in (4, 5)
)
select id as log_id,
    case
        when last_number is not null then number - last_number
        else number
    end as pity
from log
    join number on id = number.log_id
    join last_number on id = last_number.log_id;


create view if not exists star5 as
select
uid,
    name,
    pity,
    date(time),
    item_type,
    gacha_type
from log
    join pity on id = pity.log_id
where rank_type = 5;