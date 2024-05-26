import { ComponentProps } from "react";
import { Row } from "react-bootstrap";

interface CollectionSlotDisplayProps extends ComponentProps<"div"> {
    collectionSlot: CollectionSlotData
}

const CollectionSlotDisplay = (props: CollectionSlotDisplayProps) => {
    const collectionSlot = props.collectionSlot;

    return (
        <Row>
            {collectionSlot.card.name}
        </Row>
    );
}

export default CollectionSlotDisplay;