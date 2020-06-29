create table users (
    id int primary key auto_increment,
    firstname varchar(100) not null,
    lastname varchar(100) not null,
    mail varchar(150) not null unique,
    password char(60) not null 
);

create table snippets (
    id char(16) primary key,
    title varchar(300) not null,
    content TEXT not null,
    create_date date not null,
    expiration_date date not null,
    is_public BOOLEAN,
    owner_id int not null,
    FOREIGN KEY (owner_id)  REFERENCES users (id)
);