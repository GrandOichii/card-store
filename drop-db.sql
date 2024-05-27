drop table collection_slots;
drop table collections;
drop table cart_slots;
drop table carts;
drop table cards;
drop table languages;
drop table card_types;
drop table card_keys;
drop table users;

SELECT "cards"."id","cards"."created_at","cards"."updated_at","cards"."deleted_at","cards"."name","cards"."text","cards"."image_url","cards"."price","cards"."card_key_id","cards"."poster_id","cards"."card_type_id","cards"."language_id" 
FROM "cards" 
JOIN languages ON cards.language_id = languages.id 
JOIN card_types ON cards.card_type_id = card_types.id 
JOIN card_keys ON cards.card_key_id = card_keys.id 
WHERE 
(LOWER(name) like $1$ AND 
(LOWER(language_id) = $2$ 
OR LOWER(languages.long_name) = $3$ 
OR LOWER(name) like $4$ 
OR LOWER(card_type_id) = $5$ 
OR LOWER(card_types.short_name) = $6$ 
OR LOWER(card_keys.eng_name) = $7$) 
AND 
(LOWER(language_id) = $8$ 
OR LOWER(languages.long_name) = $9$ 
OR LOWER(name) like $10$ 
OR LOWER(card_type_id) = $11$ 
OR LOWER(card_types.short_name) = $12$ 
OR LOWER(card_keys.eng_name) = $13$)) 
AND "cards"."deleted_at" IS NULL LIMIT $14$