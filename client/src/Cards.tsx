import { useEffect, useState } from 'react';
import axios from './api/axios'
import CardDisplay from './components/CardDisplay';

const Cards = () => {
    const [cards, setCards] = useState<CardData[]>([])
    useEffect(() => {
        const getCards = async () => {
            // TODO catch error
            const resp = await axios.get('/card/all')
            setCards(resp.data)
        }
        getCards()
    }, [])
    return <div className='d-flex'>
        {cards.map(c => <CardDisplay card={c} />)}
    </div>
}

export default Cards;