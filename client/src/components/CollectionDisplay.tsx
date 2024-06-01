import { ComponentProps } from "react";
import { Card } from "react-bootstrap";
import { Link } from "react-router-dom";

interface CollectionDisplayProps extends ComponentProps<"div"> {
    collection: CollectionData
}

const CollectionDisplay = (props: CollectionDisplayProps) => {
    const collection = props.collection;
    return (
        <Card className="h-100">
            <Card.Body>
                <Card.Link as={Link} to={`/collections/${collection.id}`} className="stretched-link">
                    {collection.name}
                </Card.Link>
                <Card.Text>
                    {collection.description}
                </Card.Text>
            </Card.Body>
        </Card>
    )
}

export default CollectionDisplay;