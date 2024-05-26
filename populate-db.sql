-- password: password
insert into users (id, username, password_hash, email, is_admin, verified) values (1, 'admin', '$2a$14$cWg.v20w8okniqXTCw4r8u2PzaD0qeQS7ydPsx8GSf9UPvPHl2dAG', 'admin@mail.com', true, true);
insert into carts (user_id) values (1);

insert into card_types (id, long_name) values ('MTG', 'Magic: the Gathering');
insert into card_types (id, long_name) values ('YGO', 'Yu-Gi-Oh!');

insert into languages (id, long_name) values ('ENG', 'English');
insert into languages (id, long_name) values ('RUS', 'Russian');

insert into cards (name, text, image_url, price, poster_id, card_type_id, language_id) values ('Toxic Iguanar', 'Creature — Lizard\nToxic Iguanar has deathtouch as long as you control a green permanent.', 'https://cards.scryfall.io/large/front/2/8/28fd2dce-b91f-441f-a3ea-af87cc925713.jpg', 9.03, 1, 'MTG', 'RUS');
insert into cards (name, text, image_url, price, poster_id, card_type_id, language_id) values ('Song of the Worldsoul', 'Enchantment\nWhenever you cast a spell, populate.', 'https://cards.scryfall.io/large/front/b/b/bb73ec0d-f582-4f74-9b5c-180fe3aedcf6.jpg', 577.90, 1, 'MTG', 'ENG');
insert into cards (name, text, image_url, price, poster_id, card_type_id, language_id) values ('Iron Myr', 'Artifact Creature — Myr\n{T}: Add {R} yo your mana pool.', 'https://cards.scryfall.io/large/front/5/b/5bd0a588-b695-4060-b5d5-c6a74710ff0f.jpg', 22.58, 1, 'MTG', 'ENG');
insert into cards (name, text, image_url, price, poster_id, card_type_id, language_id) values ('Iron Myr', 'Artifact Creature — Myr\n{T}: Add {R} yo your mana pool.', 'https://cards.scryfall.io/large/front/5/b/5bd0a588-b695-4060-b5d5-c6a74710ff0f.jpg', 22.58, 1, 'MTG', 'RUS');
insert into cards (name, text, image_url, price, poster_id, card_type_id, language_id) values ('Moisture Creature', '[Fairy/Effect]\nIf you Tribute Summon this monster by Tributing 3 monsters on the field, destroy all Spell and Trap Cards on your opponent''s side of the field.', 'https://images.ygoprodeck.com/images/cards/75285069.jpg', 10, 1, 'YGO', 'ENG');