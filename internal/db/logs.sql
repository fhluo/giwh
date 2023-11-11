create table if not exists logs
(
    id         text primary key not null,
    uid        text,
    gacha_type text,
    name       text,
    item_id    text,
    item_type  text,
    rank_type  text,
    count      text,
    time       text,
    lang       text
);

-- 索引
create index if not exists logs_index on logs (uid, gacha_type, rank_type);

-- UID 列表
create view if not exists uid_list as
select distinct uid
from logs;

-- 编号
create view if not exists number as
select uid,
       gacha_type,
       time,
       name,
       lang,
       item_type,
       rank_type,
       id,
       row_number() over (partition by uid, case when gacha_type in (301, 400) then 301 else gacha_type end order by id) as number
from logs;

-- 上个 4/5 星的编号
create view if not exists last_number as
select *,
       lag(number)
           over (partition by uid, case when gacha_type in (301, 400) then 301 else gacha_type end, rank_type) as last_number
from number
where rank_type in (4, 5);

-- 4/5 星 pity
create view if not exists pity as
select uid,
       gacha_type,
       time,
       name,
       lang,
       item_type,
       rank_type,
       id,
       case when last_number is not null then number - last_number else number end as pity
from last_number;