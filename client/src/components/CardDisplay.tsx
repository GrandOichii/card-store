import { ComponentProps } from "react";
import { Card } from "react-bootstrap"

interface CardDisplayProps extends ComponentProps<"div"> {
    card: CardData
}

const CardDisplay = (props: CardDisplayProps) => {
    return <Card className="width: 18rem;">
        <Card.Img src={props.card.imageUrl}/>
        <Card.Body>
            <Card.Title>{props.card.name}</Card.Title>
        </Card.Body>
    </Card>
}

export default CardDisplay;

{/* <div class="card" style="width: 18rem;">
<img class="card-img-top" src="{{ .ImageUrl }}" alt="...">
<div class="card-body">
    <p class="card-title">{{ .Name }}</p>
</div>
</div> */}