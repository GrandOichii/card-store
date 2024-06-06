import { FormEvent, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Alert, Button, CardGroup, Col, Container, Form } from "react-bootstrap";

import axios from "./api/axios";
import { AxiosError, isAxiosError } from "axios";
import CollectionDisplay from "./components/CollectionDisplay";

const Collections = () => {
    const [collections, setCollections] = useState<CollectionData[]>([]);
    const navigate = useNavigate();
    const [validate, setValidate] = useState(false);
    const [newName, setNewName] = useState('');
    const [newDescription, setNewDescription] = useState('')
    const [collectionCreationError, setCollectionCreationError] = useState('');
    const [loggedIn, setLoggedIn] = useState(false);

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
            setLoggedIn(true);
            setCollections(resp.data);
        } catch (e: any) {
            if (!isAxiosError(e)) {
                console.error(e);
                return;
            }
            
            const err = e as AxiosError;
            if (err.response!.status == 401) {
                navigate('/login', {replace: false});
                return;
            }

            // TODO handle other errors
            console.error(err);
            
            return;
        }
    };
    
    useEffect(() => {
        getCollections();
    }, []);

    const onCreateCollection = async (e: FormEvent) => {
        e.preventDefault();
        setCollectionCreationError('');
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
        
        try {
            await axios.post('/collection', data, {withCredentials: true});
            getCollections();

            setNewName('');
            setNewDescription('');
        } catch (ex) {
            if (isAxiosError(ex)) {
                setCollectionCreationError(ex.message);
                return;
            }
            console.error(ex);
        }
    }

    return (loggedIn &&
        <Container>
            <h3>Collections</h3>
            <Form noValidate validated={validate} onSubmit={onCreateCollection} className="border border-primary p-2 rounded mb-1">
                <Form.Group controlId="formName">
                    <Form.Label>Name: </Form.Label>
                    <Form.Control
                        type="text"
                        value={newName}
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
                        value={newDescription}
                        placeholder="Enter collection description"
                        onChange={e => setNewDescription(e.target.value)}
                    />
                </Form.Group>
                <Button type="submit" variant="outline-primary" className="my-2">Create collection</Button>
                {collectionCreationError.length > 0 &&
                    <Alert variant="danger">Failed to create collection!</Alert>
                }
                
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