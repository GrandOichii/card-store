import { SyntheticEvent, useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';
import { Dropdown, DropdownButton, Form } from 'react-bootstrap';

const Cards = () => {
    const { type } = useParams();    
    
    const [cards, setCards] = useState<CardData[]>([])
    const [failed, setFailed] = useState(false)
    const [cardSize, setCardSize] = useState(2)

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
                        className='w-25 col'
                        onChange={onCardSizeChange}
                        defaultValue={cardSize}
                    >
                        <option value={1}>Small</option>
                        <option value={2}>Medium</option>
                    </Form.Select>
                </div>
                <div className='d-flex'>
                    {cards.map(c => 
                        <div key={c.id} className={`col-lg-${cardSize}`}>
                            <CardDisplay card={c} />
                        </div>
                    )}
                </div>
            </div>
        }
    </div>
}

export default Cards;