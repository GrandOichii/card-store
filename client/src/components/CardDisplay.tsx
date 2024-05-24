import { ComponentProps } from "react";
import { Card } from "react-bootstrap"
import { Link } from "react-router-dom";

interface CardDisplayProps extends ComponentProps<"div"> {
    card: CardData
}

const CardDisplay = (props: CardDisplayProps) => {
    const card = props.card;
    
    return <Card className={props.className}>
        <Card.Img src={card.imageUrl}/>
        <Card.Body>
            <Card.Title>
                <Card.Link as={Link} to={`/cards/${props.card.id}`} className="stretched-link">
                    {card.name}
                </Card.Link>
            </Card.Title>
        </Card.Body>
        <Card.Footer className="text-end fs-4">
            {card.price}
        </Card.Footer>
    </Card>
}

export default CardDisplay;