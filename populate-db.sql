-- password: password
insert into users (id, username, password_hash, email, is_admin, verified) values (1, 'admin', '$2a$14$cWg.v20w8okniqXTCw4r8u2PzaD0qeQS7ydPsx8GSf9UPvPHl2dAG', 'admin@mail.com', true, true);
insert into carts (user_id) values (1);

insert into card_types (id, long_name, short_name) values ('MTG', 'Magic: the Gathering', 'magic');
insert into card_types (id, long_name, short_name) values ('YGO', 'Yu-Gi-Oh!', 'yugioh');

insert into languages (id, long_name) values ('EN', 'English');
insert into languages (id, long_name) values ('RU', 'Russian');
insert into languages (id, long_name) values ('JA', 'Japanese');

insert into foilings (id, label, descriptive_name) values ('mtg_foil', 'Foil', 'Standard MTG foiling');