import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "./api/axios";
import { isAxiosError } from "axios";
import { toDescriptiveString } from "./utility/card";
import { useCookies } from "react-cookie";
import { Button, Container } from "react-bootstrap";
import { isLoggedIn } from "./auth/login";

// TODO display current price
const LargeCardDisplay = () => {
    const { id } = useParams();    
    const [card, setCard] = useState<CardData>();
    const [failed, setFailed] = useState(false);
    const [collections, setCollections] = useState<CollectionData[]>();

    const [cookies, _1, _2] = useCookies();

    const getCollections = async () => {
        // TODO handle errors
        const resp = await axios.get('/collection/all', {withCredentials: true});
        setCollections(resp.data);
    };
    
    useEffect(() => {
        getCollections();
    }, []);

    useEffect(() => {    
        getCard();
    }, []);

    const getCard = async () => {
        try {
            const resp = await axios.get(`/card/${id}`);
            setCard(resp.data);
        } catch (e) {
            setFailed(true);
        }
    }

    const addCardTo = async (collectionId: number) => {
        const resp = await axios.post(`/collection/${collectionId}`, {
            'cardId': card?.id,
            'amount': 1
        }, {withCredentials: true});
        // TODO catch errors
        // TODO add cool green notification at the bottom right of the screen that the card has been added
        getCard();
    };


    return <div className="container">
        {!!card && (
            <>
                <h1>{toDescriptiveString(card)}</h1>
                <div className="row mb-3">
                    <img className="col-lg-4" src={card?.imageUrl}></img>
                    <span style={{whiteSpace: 'pre-wrap'}} className="col">{card?.text.replace('\\n', '\n')}</span>
                </div>
                {isLoggedIn(cookies) && (
                    <>
                        <h2>Add to collections</h2>
                        {!!collections && (
                            <Container className="d-flex flex-wrap">
                                {collections.map(c => (
                                    <Button 
                                        key={c.id} 
                                        variant="primary" 
                                        className="m-1" 
                                        onClick={() => addCardTo(c.id)}>
                                    {c.name}</Button>
                                ))}
                            </Container>
                        )}
                    </>
                )}
            </>
        )}
    </div>
}

export default LargeCardDisplay