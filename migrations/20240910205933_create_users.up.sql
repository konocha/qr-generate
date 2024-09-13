create table users(
    user_id int primary key auto_increment,
    email varchar(100) unique not null,
    encrypted_password varchar(100) not null
);