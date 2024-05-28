import { ComponentProps, useEffect, useState } from "react";
import { Row } from "react-bootstrap";
import axios from "../api/axios";

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
        <Row className="rounded border border-primary p-2">
            {!!card && 
                <div className="d-flex">
                    <div className="w-100">
                        {card.name}
                    </div>
                    <div className="flex-shrink-1">
                        {collectionSlot.amount}
                    </div>                
                </div>
            }
        </Row>
    );
}

export default CollectionSlotDisplay;