import { useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';
import { useParams } from 'react-router-dom';

const Cards = () => {
    const { type } = useParams();    
    
    const [cards, setCards] = useState<CardData[]>([])
    useEffect(() => {
        const getCards = async () => {
            // TODO catch error
            const resp = await axios.get(`/card?type=${type}`)
            setCards(resp.data)
        }
        getCards()
    }, [])
    return <div className='d-flex'>
        {cards.map(c => <CardDisplay key={c.id} card={c} />)}
    </div>
}

export default Cards;