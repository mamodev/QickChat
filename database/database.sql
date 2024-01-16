create table users (
    id char (36) primary key,
    username varchar(255) not null unique,
    email varchar(255) not null,
    profile_picture varchar(255) default 'default.png',
    password varchar(255) not null,

    created_at timestamp default current_timestamp
);

create table chat (
    id char (36) primary key,
    name varchar(255) not null,
    picture varchar(255) default 'default.png',
    created_at timestamp default current_timestamp
);


create table chat_user (
    chat_id char (36),
    user_id char (36),
    created_at timestamp default current_timestamp,

    primary key (chat_id, user_id),
    foreign key (chat_id) references chat(id),
    foreign key (user_id) references users(id)
);

create table message (
    id char (36) primary key,
    sender_id char (36) not null,
    chat_id char (36) not null,

    message varchar(255) not null,

    created_at timestamp default current_timestamp,

    foreign key (sender_id) references users(id)
);


create table friend (
    user1_id char (36),
    user2_id char (36),
    created_at timestamp default current_timestamp,
    
    primary key (user1_id, user2_id),
    foreign key (user1_id) references users(id),
    foreign key (user2_id) references users(id)
);


create table friend_request (
    id char (36) primary key,
    sender_id char (36) not null,
    receiver_id char (36) not null,

    created_at timestamp default current_timestamp,

    foreign key (sender_id) references users(id),
    foreign key (receiver_id) references users(id)
);

create table push_subscription (
    id char (36) primary key,
    user_id char (36) not null,
    endpoint varchar(255) not null,

    created_at timestamp default current_timestamp,
    foreign key (user_id) references users(id)
);