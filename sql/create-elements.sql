create table if not exists elements(
    element_id serial unique,
    element_contract_id int unique,
    address varchar(64),
    name varchar(32) not null,
    n_protons int,
    provider varchar(32),
    up_votes int,
    down_votes int,
    num_fail int,
    num_pass int,
    description varchar(500),
    tx_code varchar(25),
    transaction_hash varchar(64),
    content_hash varchar(64),
    dob bigint
);

create table if not exists users(
    user_id serial unique,
    address varchar(64) not null unique,
    username varchar(30) unique,
    pfp_uri text,
    description varchar(500),
    twitter_uri text,
    github_uri text,
    is_student boolean,
    is_teacher boolean
);

create table if not exists protons(
    proton_id serial unique,
    name varchar(32) not null,
    description varchar(500),
    base_uri text,
    fk_element int references elements(element_id)
);

create table if not exists facts(
    fact_id serial unique,
    fact varchar(64) not null,
    fact_hash varchar(64),
    fact_r varchar(64),
    fact_s varchar(64),
    fact_output text,
    fact_status text
);

create table if not exists element_keys(
    keys text[] not null,
    fk_element int references elements(element_id)
);

create table if not exists element_attempts(
    passed boolean not null,
    score smallint,
    fact_id int references facts(fact_id),
    element_id int references elements(element_id),
    fk_user int references users(user_id)
);

create table if not exists proton_completions(
    passed boolean not null,
    score smallint,
    response_uri text,
    fk_proton int references protons(proton_id),
    fk_user int references users(user_id)
);