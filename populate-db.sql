-- password: password
insert into users (id, username, password_hash, email, is_admin, verified) values (1, 'admin', '$2a$14$cWg.v20w8okniqXTCw4r8u2PzaD0qeQS7ydPsx8GSf9UPvPHl2dAG', 'admin@mail.com', true, true);

insert into cards (name, text, image_url, price, poster_id) values ('Toxic Iguanar', 'Creature — Lizard\nToxic Iguanar has deathtouch as long as you control a green permanent.', 'https://cards.scryfall.io/large/front/2/8/28fd2dce-b91f-441f-a3ea-af87cc925713.jpg?1562800112', 9.03, 1);