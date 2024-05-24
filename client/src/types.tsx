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

type CardSlotData = {
    card: CardData,
    amount: number
};

type CollectionData = {
    id : number,
    name: string,
    description: string,
    cards: CardSlotData[]
};

type LoginData = {
    username: string,
    password: string
};