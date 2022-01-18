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
    accumen int,
    location text,
    primary_molecule int,
    pfp_uri text,
    description varchar(500),
    twitter_uri text,
    discord_uri text,
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

create table if not exists custom_exams(
    exam_id serial unique,
    questions text[],
    answers jsonb,
    fk_element int references elements(element_id)
)

create table if not exists element_cert_keys(
    cert_keys text[] not null,
    cert_uri text not null,
    rubric_uri text,
    fk_element int references elements(element_id)
);

create table if not exists element_attempts(
    passed boolean not null,
    score smallint,
    public_key varchar(64) not null,
    fact varchar(64),
    fact_job_id varchar(64),
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