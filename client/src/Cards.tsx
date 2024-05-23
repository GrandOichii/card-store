import { useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';

const Cards = () => {
    const { type } = useParams();    
    
    const [cards, setCards] = useState<CardData[]>([])
    const [failed, setFailed] = useState(false)

    useEffect(() => {
        const getCards = async () => {
            // TODO catch error
            try {
                const resp = await axios.get(`/card?type=${type}`);
                setCards(resp.data);
            } catch (e) {
                setFailed(true);
            }
        }
        
        getCards();
    }, []);

    return <div className='d-flex'>
        {
            failed
            ? <div className="alert alert-danger" role='alert'>
                Failed to fetch cards!
            </div>
            : cards.map(c => <CardDisplay key={c.id} card={c} />)
        }
    </div>
}

export default Cards;