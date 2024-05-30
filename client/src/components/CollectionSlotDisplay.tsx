import { ComponentProps, useEffect, useState } from "react";
import { Row } from "react-bootstrap";
import axios from "../api/axios";
import { toDescriptiveString } from "../utility/card";

interface CollectionSlotDisplayProps extends ComponentProps<"div"> {
    collectionSlot: CollectionSlotData
}

const CollectionSlotDisplay = (props: CollectionSlotDisplayProps) => {
    const collectionSlot = props.collectionSlot;    
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
        <div className="rounded border border-primary py-2 ps-2 pe-3">
            {!!card && 
                <div className="d-flex">
                    <div className="w-100">
                        {toDescriptiveString(card)}
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