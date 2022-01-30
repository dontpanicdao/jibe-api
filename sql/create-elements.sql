create table if not exists elements(
    element_id serial unique,
    element_contract_id int unique,
    address varchar(64),
    name varchar(32) not null,
    n_protons int,
    provider varchar(32),
    molecule_address varchar(64),
    reward_erc20_address varchar(64),
    reward_amount_low text,
    reward_amount_high text,
    reward_symbol varchar(10),
    cert_uri text not null,
    rubric_uri text,
    rubric_hash_low text,
    rubric_hash_high text,
    up_votes int default 0,
    down_votes int default 0,
    num_fail int default 0,
    num_pass int default 0,
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
    fk_element int references elements(element_contract_id)
);

create table if not exists custom_exams(
    exam_id serial unique,
    questions text[],
    answers jsonb,
    fk_element int references elements(element_contract_id)
);

create table if not exists element_cert_keys(
    cert_keys text[] not null,
    fk_element int references elements(element_contract_id)
);

create table if not exists element_attempts(
    passed boolean default false,
    score smallint default 0,
    public_key varchar(64) not null,
    fact varchar(64),
    fact_low varchar(32),
    fact_high varchar(32),
    fact_job_id varchar(64),
    status text,
    l1_tx text,
    l2_tx text
    fk_element int references elements(element_contract_id),
    fk_user int references users(user_id)
);

create table if not exists proton_completions(
    passed boolean not null,
    score smallint,
    response_uri text,
    fk_proton int references protons(proton_id),
    fk_user int references users(user_id)
);

create table if not exists element_contract(
    id serial,
    address varchar(64),
    owner varchar(64),
    version smallint
);

create table if not exists webauthn_sessions(
    challenge text,
    display_name text,
    user_verification text,
    public_key varchar(64)
);

create table if not exists credentials(
    aaguid text,
    credential_id text,
    public_x text,
    public_y text,
    stark_key text,
    counter int,
    display_name text
);