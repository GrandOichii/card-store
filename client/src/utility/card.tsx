

export const toDescriptiveString = (c: CardData): string => {
    return `${c.name} [${c.expansion}] (${c.language.id})`
}

export const formatPrice = (c: CardData): string => {
    return `${c.price}`;
}