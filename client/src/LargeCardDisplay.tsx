import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "./api/axios";

// TODO display current price
// TODO add button to add to a collection
const LargeCardDisplay = () => {
    const { id } = useParams();    
    const [card, setCard] = useState<CardData>();
    const [failed, setFailed] = useState(false);

     useEffect(() => {
        const getCard = async () => {
            try {
                const resp = await axios.get(`/card/${id}`);
                setCard(resp.data);
            } catch (e) {
                setFailed(true);
            }
        }
        
        getCard();
    }, []);
    return <div className="container">
        <h1>{card?.name}</h1>
        <div className="row">
            <img className="col-lg-4" src={card?.imageUrl}></img>
            <span style={{whiteSpace: 'pre-wrap'}} className="col">{card?.text.replace('\\n', '\n')}</span>
        </div>
    </div>
}

export default LargeCardDisplay