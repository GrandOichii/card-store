import { ComponentProps, useEffect, useState } from "react";
import axios from "../api/axios";
import { slotStockAmountToVariant, toDescriptiveString } from "../utility/card";

interface CollectionSlotDisplayProps extends ComponentProps<"div"> {
    collectionSlot: CollectionSlotData,
    collectionId: number
}

const CollectionSlotDisplay = (props: CollectionSlotDisplayProps) => {
    const {collectionSlot, collectionId} = props;

    const [card, setCard] = useState<CardData>();

    useEffect(() => {
        const fetch = async () => {
            // TODO handle errors
            const resp = await axios.get(`/card/${collectionSlot.cardId}`)
            setCard(resp.data);
        };
        fetch();
    }, []);

    return (
        <div className={!!card ? `rounded border border-${slotStockAmountToVariant(collectionSlot, card!)} py-2 ps-2 pe-3` : '...'}>
            {!!card && 
                <div className="d-flex">
                    <div className="w-100">
                        <a href={`/cards/${card.id}`}>{toDescriptiveString(card)}</a>
                    </div>
                    <div className="flex-shrink-1">
                        {collectionSlot.amount}
                    </div>
                </div>
            }
        </div>
    );
}

export default CollectionSlotDisplay;