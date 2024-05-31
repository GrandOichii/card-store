

export const toDescriptiveString = (c: CardData): string => {
    return `${c.name} [${c.expansion}] (${c.language.id})`
}

export const formatPrice = (c: CardData): string => {
    return `${c.price}`;
}

export const stockAmountToVariant = (c: CardData | undefined): string => {
    if (c == null) return 'dark';

    if (c.inStockAmount == 0) return 'danger';
    // TODO add 'warning' for small stock amounts
    return 'primary';
}