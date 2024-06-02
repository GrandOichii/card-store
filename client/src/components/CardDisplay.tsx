import { ComponentProps } from "react";
import { Card } from "react-bootstrap"
import { Link } from "react-router-dom";
import { toDescriptiveString, formatPrice } from "../utility/card";
import { StarFill } from "react-bootstrap-icons";

interface CardDisplayProps extends ComponentProps<"div"> {
    card: CardData
}

const CardDisplay = (props: CardDisplayProps) => {
    const card = props.card;
    
    return <Card className={props.className}>
        <Card.Img src={card.imageUrl}/>
        <Card.Body>
            <Card.Title className="d-flex align-items-center">
                {card.foiling.id.length > 0 && 
                    <StarFill className="me-2" color="gold" />
                }
                <Card.Link as={Link} to={`/cards/${props.card.id}`} className="">
                    {toDescriptiveString(card)}
                </Card.Link>
            </Card.Title>
        </Card.Body>
        <Card.Footer className="text-end fs-4">
            {formatPrice(card)}
        </Card.Footer>
    </Card>
}

export default CardDisplay;