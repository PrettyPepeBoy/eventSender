CREATE TABLE IF NOT EXISTS mails
(
    id       integer primary key,
    email    TEXT unique not null,
    password TEXT unique not null
);

CREATE TABLE IF NOT EXISTS users
(
    id    integer primary key,
    email TEXT unique not null,
    constraint email_fk FOREIGN KEY (email) references mails (email)
)