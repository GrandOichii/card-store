// TODO! implement paging! can be confusing while testing

import { FormEvent, SyntheticEvent, useEffect, useState } from 'react';
import NumericInput from "react-numeric-input";
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';
import { Button, Col, Container, Form, Row } from 'react-bootstrap';

const Cards = () => {
    const { type } = useParams();    
    
    const [queryResult, setQueryResult] = useState<CardQueryResult>();
    const [failed, setFailed] = useState(false);
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
        const getCards = async () => {
            try {
                const resp = await axios.get(`/card?type=${type}`);
                setQueryResult(resp.data);
            } catch (e) {
                setFailed(true);
            }
        }
        
        getCards();
    }, []);

    const onCardSizeChange = (e: SyntheticEvent) => {
        const select = e.target as HTMLSelectElement;
        const v = parseInt(select.value);
        setCardSize(v);
    }

    const onPageSubmit = async (e: FormEvent) => {
        e.preventDefault();
        await fetchQuery(page);
    }

    const onQuerySubmit = async (e: FormEvent) => {
        e.preventDefault();
        setPage(1);
        await fetchQuery(1  );
    }
        
    const fetchQuery = async (p: number) => {
        const url = `/card?type=${type}&page=${p}&t=${keywords}`;
        const resp = await axios.get(url);
        setQueryResult(resp.data);
        // TODO catch error
    }

    return <div>
        {
            failed
            ? <div className="alert alert-danger" role='alert'>
                Failed to fetch cards!
            </div>
            : <div>
                <Container>
                    <Row className='align-items-center'>
                        <Form.Label className='col-auto'>Card size: </Form.Label>
                        <Form.Select 
                            className='col'
                            onChange={onCardSizeChange}
                            defaultValue={cardSize}
                        >
                            <option value={4}>Small</option>
                            <option value={3}>Medium</option>
                        </Form.Select>
                    </Row>
                    <Form onSubmit={onQuerySubmit}>
                        <Row>
                            <Form.Control
                                type="Search:"
                                className='col'
                                onChange={e => setKeywords(e.target.value)}
                            />
                            <Button 
                                type="submit" 
                                variant="primary"
                                className="col-auto"
                                disabled={keywords.length === 0}
                            >Search</Button>
                        </Row>
                    </Form>
                    <Form onSubmit={onPageSubmit}>
                        <Row>
                            <Form.Label>Page:</Form.Label>
                            <NumericInput
                                min={1}
                                value={page}
                                onChange={(
                                    value: number | null, 
                                    _stringValue: string, 
                                    _input: HTMLInputElement
                                ) => setPage(value!)}
                            />
                        </Row>
                    </Form>
                </Container>
                <Container>
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
            </div>
        }
    </div>
}

export default Cards;