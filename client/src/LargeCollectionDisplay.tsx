import { FormEvent, useEffect, useState } from "react";
import { Alert, Button, Container, Form, Row } from "react-bootstrap";
import { useParams } from "react-router-dom";
import axios from "./api/axios";
import CollectionSlotDisplay from "./components/CollectionSlotDisplay";
import { AxiosResponse, isAxiosError } from "axios";
import { toDescriptiveString } from "./utility/card";

const LargeCollectionDisplay = () => {
    const [collection, setCollection] = useState<CollectionData>();
    const { id } = useParams();
    const [newCardId, setNewCardId] = useState('');
    const [fetchedCard, setFetchedCard] = useState<CardData | null>(null);
    
    const canAddCardById = (): boolean => {
        return fetchedCard != null;
    }

    const getCollection = async () => {
        try {
            const resp = await axios.get(`collection/${id}`, {withCredentials: true});
            setCollection(resp.data);
        } catch (e) {
            // TODO handle error
            console.log(e);
        }
    };
    
    useEffect(() => {
        getCollection();
    }, []);

    useEffect(() => {
        fetchCard(newCardId)
            .then((card: CardData | null) => setFetchedCard(card))
    }, [newCardId]);

    const onImportFromClipboard = async () => {
        // TODO
        const text = await navigator.clipboard.readText();
        console.log(text);
    };

    const onCardIdSubmit = async (e: FormEvent) => {
        e.preventDefault();
        if (fetchedCard == null) {
            return;            
        }

        // TODO handle errors
        const cardId = Number(newCardId);
        const resp = await axios.post(`/collection/${id}`, {
            'cardId': cardId,
            'amount': 1
        }, {withCredentials: true});
        setCollection(resp.data);
    }

    const fetchCard = async (cardID: string): Promise<CardData | null> => {
        try {   
            const resp: AxiosResponse<CardData> = await axios.get(`card/${cardID}`);
            
            return resp.data;
        } catch (ex) {
            if (!isAxiosError(ex)) {
                console.error(ex);
            }
            return null;
        }
    }

    return (
        <Container>
            {!!collection && (
                <div>
                    <h1>{collection?.name}</h1>
                    <h3>Cards</h3>
                    

                    <Form onSubmit={onCardIdSubmit} className="my-2">
                        <div className="d-flex">
                            <Form.Control
                                className="col me-1"
                                type="text"
                                placeholder="Enter card ID"
                                onChange={e => setNewCardId(e.target.value)}
                                />
                            <Button
                                className="col-auto"
                                type="submit"
                                disabled={!canAddCardById()}
                            >Add</Button>
                        </div>
                        <Form.Text>
                            {!!fetchedCard && toDescriptiveString(fetchedCard)}
                        </Form.Text>
                    </Form>
                    {/* <Button onClick={onImportFromClipboard}>Add from clipboard</Button> */}

                    {collection.cards.length === 0 && (
                        <Alert variant="info" className="my-2">No cards added yet!</Alert>
                    )}

                    {collection?.cards.map(c => (
                        <div className="my-1">
                            <CollectionSlotDisplay collectionSlot={c} />

                        </div>
                    ))}
                </div>
            )}
        </Container>
    )
}

export default LargeCollectionDisplay;