import argparse
import requests
import json

LANGUAGE = ''
CARDS_REQUEST_FORMAT = 'https://api.scryfall.com/cards/named?fuzzy={}'
RESULT_FILE_FORMAT = '{}.sql'
RESULT = ''
INSERT_CARD_KEY_FORMAT = 'INSERT INTO card_keys (id, eng_name) VALUES (\'{}\', \'{}\');'
INSERT_CARD_FORMAT = 'INSERT INTO cards (card_key_id, name, text, image_url, price, poster_id, card_type_id, language_id) VALUES (\'{}\', \'{}\', \'{}\', \'{}\', 0, 1, \'MTG\', \'{}\');'


def create_cards_request_url(card_name: str):
    return CARDS_REQUEST_FORMAT.format(card_name.replace(' ', '+'))

def format_card_key(card: dict):
    return 'mtg_{}'.format(
        card['name'].lower().replace('\'', '').replace(' ', '_')
    )

parser = argparse.ArgumentParser(
    prog='PopulateMtgCards',
    description='Fetches the cards from scryfall.com and uses them to populate the db')

parser.add_argument('-c', '--card')
parser.add_argument('-s', '--set')
parser.add_argument('-l', '--language')
parser.add_argument('-1', '--single', action='store_true')
args = parser.parse_args()

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
    if not 'large' in card['image_uris']:
        # TODO add option
        raise Exception('failed to find large image')
    card_image = card['image_uris']['large']
    card_language = card['lang']
    
    global RESULT
    RESULT += '\n' + INSERT_CARD_FORMAT.format(
        card_key,
        card_name,
        card_text,
        card_image,
        card_language.upper()
    )

def append_card_key(card):
    card_key = format_card_key(card)
    global RESULT
    RESULT += '\n' + INSERT_CARD_KEY_FORMAT.format(card_key, card['name']) + '\n'

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
    req = requests.get(card['prints_search_uri'] + '&include_multilingual=true')
    data = req.json()

    print('Fetched total of', data['total_cards'], 'cards, lowering to language selecton...')
    cards = [c for c in data['data'] if LANGUAGE is None or c['lang'] == LANGUAGE]
    print('lowered to', len(cards), 'cards total')
    count = 0
    for card in cards:
        try:
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
    print('not implemented')
    pass

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
        fetch_set(args.set)
        save_result(args.set)
        return
    print('Specify card or set to be fetched')

if __name__ == '__main__':
    print(args.single)
    main()