-- password: password
insert into users (id, username, password_hash, email, is_admin, verified) values (1, 'admin', '$2a$14$cWg.v20w8okniqXTCw4r8u2PzaD0qeQS7ydPsx8GSf9UPvPHl2dAG', 'admin@mail.com', true, true);
insert into carts (user_id) values (1);

insert into card_types (id, long_name, short_name) values ('MTG', 'Magic: the Gathering', 'magic');
insert into card_types (id, long_name, short_name) values ('YGO', 'Yu-Gi-Oh!', 'yugioh');

insert into languages (id, long_name) values ('ENG', 'English');
insert into languages (id, long_name) values ('RUS', 'Russian');

insert into card_keys (id, eng_name) values ('mtg_toxic_iguanar', 'Toxic Iguanar');
insert into cards (card_key_id, name, text, image_url, price, poster_id, card_type_id, language_id) values ('mtg_toxic_iguanar', 'Ядовитый Игуанар', 'Существо — Ящер\nЯдовитый Игуанар имеет Смертельное касание, пока вы контролируете зеленый перманент.', 'https://cards.scryfall.io/large/front/e/9/e99a47ef-cff6-474f-a697-c84b0d29f8c0.jpg', 9.03, 1, 'MTG', 'RUS');
insert into card_keys (id, eng_name) values ('mtg_song_of_the_wordsoul', 'Song of the Worldsoul');
insert into cards (card_key_id, name, text, image_url, price, poster_id, card_type_id, language_id) values ('mtg_song_of_the_wordsoul', 'Song of the Worldsoul', 'Enchantment\nWhenever you cast a spell, populate.', 'https://cards.scryfall.io/large/front/b/b/bb73ec0d-f582-4f74-9b5c-180fe3aedcf6.jpg', 577.90, 1, 'MTG', 'ENG');
insert into card_keys (id, eng_name) values ('mtg_iron_myr', 'Iron Myr');
insert into cards (card_key_id, name, text, image_url, price, poster_id, card_type_id, language_id) values ('mtg_iron_myr', 'Iron Myr', 'Artifact Creature — Myr\n{T}: Add {R} yo your mana pool.', 'https://cards.scryfall.io/large/front/5/b/5bd0a588-b695-4060-b5d5-c6a74710ff0f.jpg', 22.58, 1, 'MTG', 'ENG');
insert into cards (card_key_id, name, text, image_url, price, poster_id, card_type_id, language_id) values ('mtg_iron_myr', 'Железный Миэр', 'Артефактное Существо — Миэр\n{T}: добавьте {R} в ваше хранилище маны.', 'https://cards.scryfall.io/large/front/f/b/fbb03b29-bf90-4404-a88b-83531def35d6.jpg', 22.58, 1, 'MTG', 'RUS');
insert into card_keys (id, eng_name) values ('ygo_moisture_creature', 'Moisture Creature');
insert into cards (card_key_id, name, text, image_url, price, poster_id, card_type_id, language_id) values ('ygo_moisture_creature', 'Moisture Creature', '[Fairy/Effect]\nIf you Tribute Summon this monster by Tributing 3 monsters on the field, destroy all Spell and Trap Cards on your opponent''s side of the field.', 'https://images.ygoprodeck.com/images/cards/75285069.jpg', 10, 1, 'YGO', 'ENG');