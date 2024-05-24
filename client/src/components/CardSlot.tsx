import { ComponentProps } from "react";
import { Row } from "react-bootstrap";

interface CardSlotDisplayProps extends ComponentProps<"div"> {
    cardSlot: CardSlotData
}

const CardSlotDisplay = (props: CardSlotDisplayProps) => {
    const cardSlot = props.cardSlot;

    return (
        <Row>
            {cardSlot.card.name}
        </Row>
    );
}

export default CardSlotDisplay;