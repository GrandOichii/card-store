import { SyntheticEvent, useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';
import { Card, Col, Container, Form, Row } from 'react-bootstrap';

const Cards = () => {
    const { type } = useParams();    
    
    const [cards, setCards] = useState<CardData[]>([])
    const [failed, setFailed] = useState(false)
    const [cardSize, setCardSize] = useState(3)

    const splitCards = (): CardData[][] => {
        let result = []
        var a = [...cards];
        while(a.length) {
            result.push(a.splice(0, cardSize))
        }
        return result;
    };

    useEffect(() => {
        const getCards = async () => {
            try {
                const resp = await axios.get(`/card?type=${type}`);
                setCards(resp.data);
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

    return <div>
        {
            failed
            ? <div className="alert alert-danger" role='alert'>
                Failed to fetch cards!
            </div>
            : <div>
                <div className='row w-25'>
                    <Form.Label className='col-auto'>Card size: </Form.Label>
                    <Form.Select 
                        className='col'
                        onChange={onCardSizeChange}
                        defaultValue={cardSize}
                    >
                        <option value={4}>Small</option>
                        <option value={3}>Medium</option>
                    </Form.Select>
                </div>
                <Container>
                    {splitCards().map((row, i) => (
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