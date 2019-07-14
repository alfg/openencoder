-- auto-generated definition
create table jobs
(
  id             serial       not null,
  guid           varchar(128) not null
    constraint jobs_pk
    primary key,
  profile        varchar(128) not null,
  created_date   timestamp default CURRENT_DATE,
  status         varchar(64)
);

alter table jobs
  owner to postgres;

create unique index jobs_id_uindex
  on jobs (id);

create unique index jobs_guid_uindex
  on jobs (guid);

