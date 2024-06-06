import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "./api/axios";
import { isAxiosError } from "axios";
import { toDescriptiveString } from "./utility/card";
import { useCookies } from "react-cookie";
import { Alert, Button, Container } from "react-bootstrap";
import { isLoggedIn } from "./auth/login";
import FailedToAddToCollectionModal from "./components/FailedToAddToCollectionModal";

// TODO display current price
// TODO display foiling and card variant
const LargeCardDisplay = () => {
    const { id } = useParams();    
    const [card, setCard] = useState<CardData>();
    const [notFound, setNotFound] = useState(false);
    const [collections, setCollections] = useState<CollectionData[]>();
    const [showFailModal, setShowFailModal] = useState(false);
    const [lastCollectionName, setLastCollectionName] = useState('');

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
        } catch (ex) {
            if (isAxiosError(ex)) {
                console.log(ex.code);
                
                if (ex.response?.status === 404) {
                    setNotFound(true);
                    return;
                }
            };
            console.error(ex);
        }
    }

    const addCardTo = async (collection: CollectionData) => {
        const collectionId = collection.id;
        try {
            const resp = await axios.post(`/collection/${collectionId}`, {
                'cardId': card?.id,
                'amount': 1
            }, {withCredentials: true});
            // TODO add cool green notification at the bottom right of the screen that the card has been added
            getCard();
        } catch (ex) {
            if (isAxiosError(ex)) {
                // TODO handle error
                // return;
            }
            setLastCollectionName(collection.name);
            setShowFailModal(true);
            console.error(ex);
        }
    };


    return <div className="container">
        {notFound &&
            <Alert variant="danger">{`Card with id ${id} not found!`}</Alert>
        }
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
                                        onClick={() => addCardTo(c)}>
                                        {/* onClick={() => setShowFailModal(true)}> */}
                                    {c.name}</Button>
                                ))}
                            </Container>
                        )}
                    </>
                )}
            </>
        )}
        <FailedToAddToCollectionModal
            onHide={() => setShowFailModal(false)}
            show={showFailModal}
            cardName={toDescriptiveString(card)}
            collectionName={lastCollectionName}
        />
    </div>
}

export default LargeCardDisplay