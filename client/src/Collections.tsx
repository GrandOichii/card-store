import { FormEvent, useEffect, useState } from "react";
import axios from "./api/axios";
import { AxiosError, isAxiosError } from "axios";
import { useNavigate } from "react-router-dom";
import CollectionDisplay from "./components/CollectionDisplay";
import { Alert, Button, CardGroup, Col, Container } from "react-bootstrap";
import { Form } from "react-bootstrap";


const Collections = () => {
    const [collections, setCollections] = useState<CollectionData[]>([]);
    const navigate = useNavigate();
    const [validate, setValidate] = useState(false);
    const [newName, setNewName] = useState('');
    const [newDescription, setNewDescription] = useState('')

    const perRow = 4;
    const splitCollections = (): CollectionData[][] => {
        let result = []
        var a = [...collections];
        while(a.length) {
            result.push(a.splice(0, perRow))
        }
        return result;
    };

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
                navigate('/login', {replace: false});
                return;
            }
            console.log(err);
            // TODO handle other errors
            const data: any = err.response?.data;
            
            return;
        }        
    };
    useEffect(() => {
        getCollections();
    }, []);

    const onCreateCollection = async (e: FormEvent) => {
        e.preventDefault();
        const form = e.currentTarget as HTMLFormElement
        if (!form.checkValidity()) {
            e.stopPropagation();
            setValidate(true)
            return;
        }
        const data = {
            'name': newName,
            'description': newDescription,
        }
        console.log(data);
        
        try {
            await axios.post('/collection', data, {withCredentials: true});
            getCollections();
        } catch (e) {
            console.log(e);
            // TODO handle
        }
    }

    // TODO add validation
    return (
        <Container>
            <h3>Collections</h3>
            <Form noValidate validated={validate} onSubmit={onCreateCollection} className="border border-primary p-2 rounded mb-1">
                <Form.Group controlId="formName">
                    <Form.Label>Name: </Form.Label>
                    <Form.Control
                        type="text"
                        placeholder="Enter collection name"
                        onChange={e => setNewName(e.target.value)}
                        required
                        minLength={3}
                    />
                    <Form.Control.Feedback type="invalid">Collection name has to be at least 3 characters long!</Form.Control.Feedback>
                </Form.Group>
                <Form.Group controlId="formDescription">
                    <Form.Label>Description: </Form.Label>
                    <Form.Control
                        type="text"
                        placeholder="Enter collection description"
                        onChange={e => setNewDescription(e.target.value)}
                    />
                </Form.Group>
                <Button type="submit" variant="outline-primary" className="my-2">Create collection</Button>
            </Form>
            {collections.length > 0 
                ? splitCollections().map((row, i) => (
                    <CardGroup key={i}>
                        {/* TODO margins mess up the ordering */}
                        {row.map(c => (
                            <Col key={c.id} className={`col-${12/perRow}`}>
                            <CollectionDisplay key={c.id} collection={c} />
                            </Col>
                        ))}
                    </CardGroup>
                ))
                : <Alert>No collections created yet!</Alert>
            }
        </Container>
    );
}

export default Collections;