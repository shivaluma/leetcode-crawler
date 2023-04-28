create table if not exists submissions (
    id serial primary key ,
    leetcode_username varchar(40) ,
    github_username varchar(40),
    problem_id varchar(40),
    submission_id varchar(40),
    sha      varchar(80) null ,
    updated_at timestamp not null
);


create unique index idx_submission_username_problem_id on submissions(leetcode_username, problem_id);
