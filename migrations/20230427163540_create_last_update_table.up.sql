create table last_updated_timestamp (
    leetcode_username varchar(40) primary key ,
    github_username varchar(40) unique not null ,
    last_updated_at timestamp not null
)