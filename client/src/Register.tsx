import { FormEvent, useState } from "react"
import { Button, Form } from "react-bootstrap"
import axios from "./api/axios";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";

const Register = () => {
    // TODO block login button when processing request
    // TODO add input checks
    const [username, setUsername] = useState<string>('');
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [_1, setCookie, _2] = useCookies();
    const navigate = useNavigate()

    const canSubmit = (): boolean => {
        // TODO
        return true
    }

    const onSubmit = async (e: FormEvent) => {
        e.preventDefault();
        const registerData = {
            'username': username,
            'email': email,
            'password': password,
        };
        const loginData = {
            'username': username,
            'password': password,
        };
        try {
            await axios.post('/auth/register', registerData, {
                withCredentials: true,
            });
            // TODO duplicated code
            await axios.post('/auth/login', loginData, {
                withCredentials: true,
            });
            setCookie('loggedIn', true, {
                maxAge: 3600
            })
            // TODO change to personal page
            navigate("/about")
        } catch (ex) {
            // TODO handle error
            console.log(ex);
        }
    }

    return <Form onSubmit={onSubmit}>
        <Form.Group controlId="formEmail">
            <Form.Label>Email</Form.Label>
            <Form.Control 
                type="text" 
                placeholder="Enter email" 
                onChange={e => setEmail(e.target.value)}
            />
        </Form.Group>
        <Form.Group controlId="formUsername">
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
            Register
        </Button>
    </Form>
}

export default Register