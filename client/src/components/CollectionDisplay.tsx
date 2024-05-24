import { ComponentProps } from "react";
import { Card } from "react-bootstrap";


interface CollectionDisplayProps extends ComponentProps<"div"> {
    collection: CollectionData
}

const CollectionDisplay = (props: CollectionDisplayProps) => {
    const collection = props.collection;
    return (
        <Card>
            <Card.Body>
                <Card.Title>
                    {collection.name}
                </Card.Title>
                <Card.Text>
                    {collection.description}
                </Card.Text>
            </Card.Body>
        </Card>
    )
}

export default CollectionDisplay;