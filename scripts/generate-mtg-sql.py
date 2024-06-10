import argparse
import requests
import json

LANGUAGE = ''
CARDS_URL_FORMAT = 'https://api.scryfall.com/cards/named?fuzzy={}'
RESULT_FILE_FORMAT = '{}.sql'
RESULT = ''
INSERT_CARD_KEY_FORMAT = 'INSERT INTO card_keys (id, eng_name) VALUES (\'{}\', \'{}\') ON CONFLICT DO NOTHING;'
INSERT_CARD_FORMAT = 'INSERT INTO cards (card_key_id, name, text, image_url, price, in_stock_amount, poster_id, card_type_id, language_id, expansion_id) VALUES (\'{}\', \'{}\', \'{}\', \'{}\', 0, 0, 1, \'MTG\', \'{}\', \'{}\') ON CONFLICT DO NOTHING;'
INSERT_EXPANSION_FORMAT = 'INSERT INTO expansions (id, short_name, full_name) VALUES (\'{}\', \'{}\', \'{}\') ON CONFLICT DO NOTHING;'
SET_SEARCH_URL_FORMAT = 'https://api.scryfall.com/sets/{}'

def create_cards_request_url(card_name: str):
    return CARDS_URL_FORMAT.format(card_name.replace(' ', '+'))

def format_card_key(card: dict):
    return 'mtg_{}'.format(
        card['name'].lower().replace('\'', '').replace(' ', '_')
    )

def format_card_expansion(card: dict):
    return 'mtg_{}'.format(card['set'])

parser = argparse.ArgumentParser(
    prog='PopulateMtgCards',
    description='Fetches the cards from scryfall.com and uses them to populate the db')

parser.add_argument('-c', '--card')
parser.add_argument('-s', '--set')
parser.add_argument('-l', '--language')
parser.add_argument('-1', '--single', action='store_true')
parser.add_argument('-i', '--image', default='large')
args = parser.parse_args()

def append_set(card):    
    global RESULT
    RESULT += '\n' + INSERT_EXPANSION_FORMAT.format('mtg_' + card['set'], card['set'].upper(), card['set_name'])

def append_card(card):
    card_key = format_card_key(card)
    card_name = card['name']
    card_text = ''
    if card['lang'] != 'en':
        card_name = card['printed_name']
    else:
        card_text = card['oracle_text']
    card_name = card_name.replace('\'', '\'\'')
    card_text = card_text.replace('\'', '\'\'').replace('\n', '\\n')
    if not args.image in card['image_uris']:
        raise Exception(f'failed to find {args.image} image')
    card_image = card['image_uris'][args.image]
    card_language = card['lang']
    
    global RESULT
    RESULT += INSERT_CARD_FORMAT.format(
        card_key,
        card_name,
        card_text,
        card_image,
        card_language.upper(),
        format_card_expansion(card)
    ) + '\n'

def append_card_key(card):
    card_key = format_card_key(card)
    global RESULT
    RESULT += '\n' + INSERT_CARD_KEY_FORMAT.format(card_key, card['name'].replace('\'', '\'\'')) + '\n'

def fetch_card(name: str):
    print('fetching using fuzzy search...')
    url = create_cards_request_url(name)
    req = requests.get(url)
    card = req.json()
    if card['object'] == 'error':
        print('Failed to find card:', card['details'])
        quit(1)
    print('fetched instance, appending card key creation...')
    append_card_key(card)
    print('appended, fetching all printings...')
    data = {'has_more': True}
    url = card['prints_search_uri'] + '&include_multilingual=true'
    cards = []
    page = 0
    while data['has_more']:
        req = requests.get(url)
        data = req.json()
        if data['object'] == 'error':
            print('Failed to fetch page', page, ':', data['details'])
            quit(1)
        url = data['next_page'] if 'next_page' in data else ''
        print('fetched page', page)
        page += 1
        cards += data['data']

    print('Fetched total of', len(cards), 'cards, lowering to language selecton...')
    cards = [c for c in cards if LANGUAGE is None or c['lang'] == LANGUAGE]
    print('lowered to', len(cards), 'cards total')
    count = 0
    for card in cards:
        try:
            print('appending set', card['set_name'])
            append_set(card)
            print('appending card', card['name'])
            append_card(card)
            count += 1
            if args.single:
                print('appended single!')
                return
            print('appended!')
        except Exception as e:
            print('failed to append card:', e)
    print(f'in total, appended {count}/{len(cards)} cards')
    
def fetch_set(set_name: str):
    print(f'fetching set set {set_name}...')
    resp = requests.get(SET_SEARCH_URL_FORMAT.format(set_name))
    data = resp.json()
    if data['object'] == 'error':
        print('Failed to fetch set:', data['details'])
        quit(1)
    print('fetched set data, searching set cards...')
    url = data['search_uri'] + '&include_multilingual=true'
    data = {'has_more': True}
    cards = []
    page = 0
    while data['has_more']:
        resp = requests.get(url)
        data = resp.json()
        if data['object'] == 'error':
            print('failed to fetch page', page, ':', data['details'])
            quit(1)
        cards += data['data']
        url = data['next_page'] if 'next_page' in data else ''
        page += 1
        print('fetched page', page)
    print(f'in total fetched {len(cards)}, narrowing down to selected language...')
    lang_map = {}
    for card in cards:
        if not card['lang'] in lang_map:
            lang_map[card['lang']] = 0
        lang_map[card['lang']] += 1
    print('language counts:')
    for key, value in lang_map.items():
        print(f'{key}: {value}')
    print('selected language:', LANGUAGE)
    cards = [c for c in cards if LANGUAGE is None or c['lang'] == LANGUAGE]
    print(f'narrowed down to {len(cards)} cards, generating sql...')
    count = 0
    card_set = cards[0]['set']
    print('appending set', card_set)
    append_set(cards[0])
    for card in cards:
        try:
            append_card_key(card)
            append_card(card)
            count += 1
            print('appended!')
        except Exception as e:
            print('failed to append card', card['id'], ':', e)
    print(f'in total, appended {count}/{len(cards)} cards')

def save_result(name):
    fname = RESULT_FILE_FORMAT.format(name)
    print('saving result to', fname)
    open(fname, 'w').write(RESULT)
    print('done!')

def main():
    global LANGUAGE
    LANGUAGE = args.language
    if args.card is not None:
        fetch_card(args.card)
        save_result(args.card)
        return
    if args.set is not None:
        fetch_set(args.set.lower())
        save_result(args.set)
        return
    print('Specify card or set to be fetched')

if __name__ == '__main__':
    main()