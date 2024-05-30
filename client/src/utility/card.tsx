

export const toDescriptiveString = (c: CardData): string => {
    return `${c.name} [${c.expansion}] (${c.language.id})`
}