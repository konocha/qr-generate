create table qr(
    qr_id int primary key auto_increment,
    qr_value varchar(100) not null, 
    user_id int not null,
    foreign key (user_id) references users(user_id) on delete cascade
);