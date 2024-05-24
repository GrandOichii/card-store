import { useEffect, useState } from "react";
import { Container } from "react-bootstrap";
import { useParams } from "react-router-dom";
import axios from "./api/axios";



const LargeCollectionDisplay = () => {
    const [collection, setCollection] = useState<CollectionData>();
    const { id } = useParams();

    const getCard = async () => {
        try {
            const resp = await axios.get(`collection/${id}`, {withCredentials: true});
            setCollection(resp.data);
        } catch (e) {
            // TODO handle error
            console.log(e);
        }
    };
    useEffect(() => {
        getCard()
    }, []);

    return (
        <Container>
            <h1>{collection?.name}</h1>
        </Container>
    )
}

export default LargeCollectionDisplay;