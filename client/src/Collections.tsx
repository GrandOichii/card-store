import { useEffect, useState } from "react";
import axios from "./api/axios";
import { AxiosError, isAxiosError } from "axios";
import { useNavigate } from "react-router-dom";
import CollectionDisplay from "./components/CollectionDisplay";
import { Alert, Col, Container, Row } from "react-bootstrap";


const Collections = () => {
    const [collections, setCollections] = useState<CollectionData[]>([]);
    const navigate = useNavigate();
    const perRow = 4;
    const splitCollections = (): CollectionData[][] => {
        let result = []
        var a = [...collections];
        while(a.length) {
            result.push(a.splice(0, perRow))
        }
        return result;
    };

    useEffect(() => {
        const getCollections = async () => {
            try {
                const resp = await axios.get('/collection/all', {withCredentials: true});
                setCollections(resp.data);
            } catch (e: any) {
                if (!isAxiosError(e)) {
                    console.log(e);
                    return;
                }
                
                const err = e as AxiosError;
                if (err.response!.status == 401) {
                    navigate('/login');
                    return;
                }
                console.log(err);
                const data: any = err.response?.data;
                
                return;
            }        
        };
        getCollections();
    }, []);

    return <div>
        <Container>
            {collections.length > 0 
                ? splitCollections().map((row, i) => (
                    <Row key={i} className='mb-2'>
                        {row.map(c => (
                            <Col key={c.id} className={`col-${12/perRow}`}>
                                <CollectionDisplay collection={c} />
                            </Col>
                        ))}
                    </Row>
                ))
                : <Alert>No collections created yet!</Alert>
            }
        </Container>
    </div>
}

export default Collections;