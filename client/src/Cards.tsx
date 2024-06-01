import { FormEvent, SyntheticEvent, useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';
import { Alert, Button, Col, Container, Form, Image, Offcanvas, Row } from 'react-bootstrap';
import { toDescriptiveString } from './utility/card';
import { isAxiosError } from 'axios';

const Cards = () => {
    // TODO add adding to cart in offcanvas
    const { type } = useParams();    
    
    const [queryResult, setQueryResult] = useState<CardQueryResult>();
    const [cardSize, setCardSize] = useState(3);
    const [keywords, setKeywords] = useState('');
    const [page, setPage] = useState(1);
    const [selectedCard, setSelectedCard] = useState<CardData | null>();
    const [collections, setCollections] = useState<CollectionData[]>();
    const [queryFailed, setQueryFailed] = useState(false);

    const splitCards = (): CardData[][] => {
        let result = []
        var a = [...queryResult!.cards];
        while(a.length) {
            result.push(a.splice(0, cardSize))
        }
        return result;
    };

    useEffect(() => {
        fetchCards();
        getCollections();
    }, []);
    
    useEffect(() => {
        setQueryResult({
            cards: [],
            totalCards: 0,
            perPage: 0
        });
        
        fetchCards();
    }, [page]);

    const getCollections = async () => {
        try {
            const resp = await axios.get('/collection/all', {withCredentials: true});
            setCollections(resp.data);
        } catch (ex) {
            // TODO handle
            if (isAxiosError(ex)) {
                return;
            }
            console.error(ex);
        }
    };

    const onCardSizeChange = (e: SyntheticEvent) => {
        const select = e.target as HTMLSelectElement;
        const v = parseInt(select.value);
        setCardSize(v);
    }

    const onQuerySubmit = async (e: FormEvent) => {
        e.preventDefault();
        
        if (page == 1) {
            fetchCards();
            return;
        }
        setPage(1);
    }
        
    const fetchCards = async () => {
        setQueryFailed(false);
        let url = `/card?type=${type}&page=${page}`;
        if (keywords.length > 0) {
            url += `&t=${keywords}`
        }
        
        try {
            const resp = await axios.get(url);
            setQueryResult(resp.data);
        } catch (ex) {
            console.error(ex);
            
            setQueryFailed(true);
        }
        window.scrollTo(0, 0);
    }

    const maxPage = (): number => {
        return Math.ceil(queryResult!.totalCards / queryResult!.perPage)
    }

    const gotoPage = (p: number) => {
        if (p <= 0) p = 1;
        const maxP = maxPage();
        if (p > maxP) p = maxP;
        setPage(p);
    };

    const modPage = (pMod: number) => {
        gotoPage(page + pMod);
    }

    const [show, setShow] = useState(false);

    const handleClose = () => setShow(false);
    const handleShow = (c: CardData) => {
        setSelectedCard(c);
        setShow(true);
    } 

    const addSelectedCardTo = async (collectionId: number) => {
        const resp = await axios.post(`/collection/${collectionId}`, {
            'cardId': selectedCard?.id,
            'amount': 1
        }, {withCredentials: true});
        // TODO add some feedback
    };

    return <div>
        {/* failed
        ? <div className="alert alert-danger" role='alert'>
            Failed to fetch cards!
        </div> */}
        <div>
        {/* <Button variant="primary" onClick={handleShow}>
            Launch
        </Button> */}

        <Offcanvas show={show} onHide={handleClose} placement='end'>
                {!!selectedCard && (
                    <div>
                        <Offcanvas.Header closeButton>
                            <Offcanvas.Title>{toDescriptiveString(selectedCard)}</Offcanvas.Title>
                        </Offcanvas.Header>
                        <div className="mx-3">
                            <Image
                                className='img-fluid'
                                src={selectedCard.imageUrl}
                            />
                        </div>
                        <Offcanvas.Body>
                            {!!collections 
                                ? (
                                    <>
                                    <h2>Add card to collections</h2>
                                        <Container className="d-flex flex-wrap">
                                            {collections.map(c => (
                                                <Button 
                                                    key={c.id} 
                                                    variant="primary" 
                                                    className="m-1" 
                                                    onClick={() => addSelectedCardTo(c.id)}>
                                                {c.name}</Button>))
                                            }
                                        </Container>
                                    </>
                                )
                                :<Alert variant='info'>Log in to add cards to collections!</Alert>
                            }
                        </Offcanvas.Body>
                    </div>
                )}
        </Offcanvas>
            <Container>
                <div className='d-flex my-1 align-items-center'>
                    <div className='text-nowrap me-2'>Card size: </div>
                    <Form.Select 
                        className=''
                        onChange={onCardSizeChange}
                        defaultValue={cardSize}
                    >
                        <option value={4}>Small</option>
                        <option value={3}>Medium</option>
                    </Form.Select>
                </div>
                <Form onSubmit={onQuerySubmit} className='d-flex my-1'>
                    <Form.Control
                        placeholder='Enter keywords'
                        type="text"
                        className='me-1'
                        onChange={e => setKeywords(e.target.value)}
                    />
                    <Button 
                        type="submit" 
                        variant="primary"
                        className=""
                        disabled={keywords.length === 0}
                    >Search</Button>
                </Form>
                <Container className='my-3'>
                    {queryFailed && <Alert variant='danger'>Error while fetching cards!</Alert>}
                    {!!queryResult && (
                        queryResult.cards.length == 0
                        ? <Alert variant='warning'>No cards found!</Alert>
                        : splitCards().map((row, i) => (
                            <Row key={i} className='mb-2'>
                                {row.map(c => (
                                    <Col key={c.id} className={`col-${12/cardSize}`} onClick={() => handleShow(c)}>
                                        <CardDisplay card={c} />
                                    </Col>
                                ))}
                            </Row>
                        )
                    ))}
                </Container>
                {queryResult && (
                    <div className="d-flex justify-content-center my-3">
                        <button type="button" className="btn btn-outline-info mx-1" onClick={() => gotoPage(1)}>&lt;&lt;</button>
                        <button type="button" className="btn btn-outline-info mx-1" onClick={() => modPage(-1)}>&lt; Previous</button>
                        <button type="button" className="btn btn-outline-info mx-1" onClick={() => modPage(1)}>Next &gt;</button>
                        <button type="button" className="btn btn-outline-info mx-1" onClick={() => gotoPage(maxPage())}>&gt;&gt;</button>
                    </div>
                )}
            </Container>
        </div>
    </div>
}

export default Cards;