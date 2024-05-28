import { FormEvent, useEffect, useState } from "react";
import { Alert, Button, Container, Form, Row } from "react-bootstrap";
import { useParams } from "react-router-dom";
import axios from "./api/axios";
import CollectionSlotDisplay from "./components/CollectionSlotDisplay";

const LargeCollectionDisplay = () => {
    const [collection, setCollection] = useState<CollectionData>();
    const { id } = useParams();
    const [newCardId, setNewCardId] = useState('');

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

    const onCardIdSubmit = async (e: FormEvent) => {
        e.preventDefault();
        // TODO handle errors
        const cardId = Number(newCardId);
        const resp = await axios.post(`/collection/${id}`, {
            'cardId': cardId,
            'amount': 1
        }, {withCredentials: true});
        console.log(resp.data)
        setCollection(resp.data);
    }

    return (
        <Container>
            {!!collection && (
                <div>
                    <h1>{collection?.name}</h1>
                    <h3>Cards</h3>
                    {collection.cards.length === 0 && (
                        <Alert variant="info">No cards added yet!</Alert>
                    )}

                    <Form onSubmit={onCardIdSubmit}>
                        <Row>
                            <Form.Control
                                className="col"
                                type="text"
                                placeholder="Enter card ID"
                                onChange={e => setNewCardId(e.target.value)}
                            />
                            <Button
                                className="col-auto"
                                type="submit"
                            >Add</Button>
                        </Row>
                    </Form>
                    {/* <Button onClick={onImportFromClipboard}>Add from clipboard</Button> */}


                    {collection?.cards.map(c => (
                        <CollectionSlotDisplay collectionSlot={c} />
                    ))}
                </div>
            )}
        </Container>
    )
}

export default LargeCollectionDisplay;