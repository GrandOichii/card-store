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

type CollectionData = {
    id : number,
    name: string,
    description: string,
    // TODO add card slots
};

type LoginData = {
    username: string,
    password: string
};