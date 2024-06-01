-- password: password
insert into users (id, username, password_hash, email, is_admin, verified) values (1, 'admin', '$2a$14$cWg.v20w8okniqXTCw4r8u2PzaD0qeQS7ydPsx8GSf9UPvPHl2dAG', 'admin@mail.com', true, true);
insert into carts (user_id) values (1);

insert into card_types (id, long_name, short_name) values ('MTG', 'Magic: the Gathering', 'magic');
insert into card_types (id, long_name, short_name) values ('YGO', 'Yu-Gi-Oh!', 'yugioh');

insert into languages (id, long_name) values ('EN', 'English');
insert into languages (id, long_name) values ('RU', 'Russian');
insert into languages (id, long_name) values ('JA', 'Japanese');

insert into foilings (id, label, descriptive_name) values ('mtg_foil', 'Foil', 'Standard MTG foiling');
INSERT INTO cards (card_key_id, name, text, image_url, price, in_stock_amount, poster_id, card_type_id, language_id, expansion_id, foiling_id) VALUES ('mtg_zoetic_cavern', 'Zoetic Cavern', '{T}: Add {C}.\nMorph {2} (You may cast this card face down as a 2/2 creature for {3}. Turn it face up any time for its morph cost.)', 'https://cards.scryfall.io/large/front/3/0/30f1373a-e0ee-4ba2-b5c0-14a6efca4df8.jpg?1699878537', 0, 0, 1, 'MTG', 'EN', 'mtg_plst', 'mtg_foil') ON CONFLICT DO NOTHING;
