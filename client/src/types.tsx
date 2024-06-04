type CardType = {
    name: string
    longName: string
};

type CardTypeData = {
    longName: string
    shortName: string
};

type LanguageData = {
    id: string
    longName: string
};

type ExpansionData = {
    id: string
	shortName: string
	fullName: string
};

type FoilingData = {
    id: string
    label: string
    descriptiveName: string
};

type CardData = {
    id: string
    name: string
    text: string
    imageUrl: string
    price: number
    type: CardTypeData
    language: LanguageData
    key: string
    expansion: string
    expansionName: string
    inStockAmount: number
    foiling: FoilingData
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

type PrivateUserData = {
    username: string
    verified: boolean
    isAdmin: boolean
}