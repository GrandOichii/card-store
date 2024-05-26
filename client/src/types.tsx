type CardType = {
    name: string,
    longName: string
};

type CardData = {
    id: string,
    imageUrl: string,
    name: string,
    price: number,
    text: string
};

type CollectionSlotData = {
    card: CardData,
    amount: number
};

type CollectionData = {
    id : number,
    name: string,
    description: string,
    cards: CollectionSlotData[]
};

type LoginData = {
    username: string,
    password: string
};