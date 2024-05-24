import { useEffect, useState } from "react";
import { Alert, Button, Container, Row } from "react-bootstrap";
import { useParams } from "react-router-dom";
import axios from "./api/axios";
import CardSlotDisplay from "./components/CardSlot";

const LargeCollectionDisplay = () => {
    const [collection, setCollection] = useState<CollectionData>();
    const { id } = useParams();

    const getCard = async () => {
        try {
            const resp = await axios.get(`collection/${id}`, {withCredentials: true});
            setCollection(resp.data);
        } catch (e) {
            // TODO handle error
            console.log(e);
        }
    };
    useEffect(() => {
        getCard()
    }, []);

    const onImportFromClipboard = async () => {
        const text = await navigator.clipboard.readText();
        console.log(text);
        
    };

    return (
        <Container>
            {!!collection && (
                <div>
                    <h1>{collection?.name}</h1>
                    <h3>Cards</h3>
                    {collection.cards.length === 0 && (
                        <Alert variant="info">No cards added yet!</Alert>
                    )}
                    <Button onClick={onImportFromClipboard}>Add from clipboard</Button>
                    {collection?.cards.map(c => (
                        <CardSlotDisplay cardSlot={c} />
                    ))}
                </div>
            )}
        </Container>
    )
}

export default LargeCollectionDisplay;