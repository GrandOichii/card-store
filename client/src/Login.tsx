import { FormEvent, useState } from "react"
import { Button, Form } from "react-bootstrap"
import axios from "./api/axios";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router-dom";

const Login = () => {
    // TODO block login button when processing request
    // TODO add input checks
    // TODO add error handling
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [_1, setCookie, _2] = useCookies();
    const navigate = useNavigate();

    const canSubmit = (): boolean => {
        return true
    }

    const onSubmit = async (e: FormEvent) => {
        e.preventDefault();
        const data = {
            'username': username,
            'password': password,
        };
        try {
            await axios.post('/auth/login', data, {
                withCredentials: true,
            });
            setCookie('loggedIn', true, {
                maxAge: 3600
            })
            
            // TODO change to personal page
            navigate("/about")
        } catch (ex) {
            console.log(ex);
            
        }
        
    }

    return <Form onSubmit={onSubmit}>
        <Form.Group controlId="formEmail">
            <Form.Label>Username</Form.Label>
            <Form.Control 
                type="text" 
                placeholder="Enter username" 
                onChange={e => setUsername(e.target.value)}
            />
        </Form.Group>
        <Form.Group className="mb-3" controlId="formPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control 
                type="password" 
                placeholder="Password" 
                onChange={e => setPassword(e.target.value)}
            />
        </Form.Group>
        <Button variant="primary" type="submit" disabled={!canSubmit()}>
            Submit
        </Button>
    </Form>
}

export default Login