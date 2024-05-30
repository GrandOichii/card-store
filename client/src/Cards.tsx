// TODO! implement paging! can be confusing while testing

import { FormEvent, SyntheticEvent, useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';
import { Button, Col, Container, Form, Row } from 'react-bootstrap';

const Cards = () => {
    // TODO add alert when no cards were found
    const { type } = useParams();    
    
    const [queryResult, setQueryResult] = useState<CardQueryResult>();
    const [cardSize, setCardSize] = useState(3);
    const [keywords, setKeywords] = useState('');
    const [page, setPage] = useState(1);

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
    }, []);
    
    useEffect(() => {
        setQueryResult({
            cards: [],
            totalCards: 0,
            perPage: 0
        });
        
        fetchCards();
    }, [page]);

    const onCardSizeChange = (e: SyntheticEvent) => {
        const select = e.target as HTMLSelectElement;
        const v = parseInt(select.value);
        setCardSize(v);
    }

    // const onPageSubmit = async (e: FormEvent) => {
    //     e.preventDefault();
    //     await fetchCards();
    // }

    const onQuerySubmit = async (e: FormEvent) => {
        e.preventDefault();
        
        if (page == 1) {
            fetchCards();
            return;
        }
        setPage(1);
    }
        
    const fetchCards = async () => {
        let url = `/card?type=${type}&page=${page}`;
        if (keywords.length > 0) {
            url += `&t=${keywords}`
        }
        
        const resp = await axios.get(url);
        setQueryResult(resp.data);
        window.scrollTo(0, 0);

        // TODO catch error
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

    return <div>
        {/* failed
        ? <div className="alert alert-danger" role='alert'>
            Failed to fetch cards!
        </div> */}
        <div>
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
                    {!!queryResult && splitCards().map((row, i) => (
                        <Row key={i} className='mb-2'>
                            {row.map(c => (
                                <Col key={c.id} className={`col-${12/cardSize}`}>
                                    <CardDisplay card={c} />
                                </Col>
                            ))}
                        </Row>
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