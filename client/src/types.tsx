type CardType = {
    name: string
    longName: string
};

type CardData = {
    id: string
    imageUrl: string
    name: string
    price: number
    text: string
};

type CollectionSlotData = {
    cardId: number
    amount: number
};

type CollectionData = {
    id : number
    name: string
    description: string
    cards: CollectionSlotData[]
};

type LoginData = {
    username: string
    password: string
};

type CardQueryResult = {
    cards: CardData[]
    totalCards: number
    perPage: number
};