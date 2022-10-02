DROP TABLE if EXISTS users;

CREATE TABLE users (
                       id SERIAL NOT NULL,
                       created_at timestamp,
                       username VARCHAR(64),
                       password_hash varchar(128),
                       email VARCHAR(64),
                       stuff BOOLEAN
);

DROP TABLE if EXISTS clients;

CREATE TABLE clients(
                        id SERIAL NOT NULL,
                        contact_login varchar(64),
                        end_date DATE NOT NULL ,
                        ordered_proxy_dhcp INTEGER,
                        payment float
);

DROP TABLE if EXISTS proxy;

CREATE TABLE proxy(
                      id serial not null,
                      dhcp serial not null,
                      tel_number integer
);