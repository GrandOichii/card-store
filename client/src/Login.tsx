import { FormEvent, useState } from "react"
import { Button, Form } from "react-bootstrap"
import axios from "./api/axios";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router-dom";
import login from './auth/login'

const Login = () => {
    // TODO block login button when processing request
    // TODO add input checks
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [_1, setCookie, _2] = useCookies();
    const [validate, setValidate] = useState(false)
    const navigate = useNavigate();

    const canSubmit = (): boolean => {
        // TODO
        return true
    }

    const onSubmit = async (e: FormEvent) => {
        e.preventDefault();

        const form = e.currentTarget as HTMLFormElement
        if (!form.checkValidity()) {
            e.stopPropagation();
            setValidate(true)
            return;
        }
        
        const data: LoginData = {
            'username': username,
            'password': password,
        };
        try {
            await login(axios, data, setCookie, navigate)
        } catch (ex) {
            // TODO add error handling
            console.log(ex);
        }
        
    }

    return <Form noValidate validated={validate} onSubmit={onSubmit}>
        <Form.Group controlId="formUsername">
            <Form.Label>Username</Form.Label>
            <Form.Control 
                type="text" 
                placeholder="Enter username" 
                onChange={e => setUsername(e.target.value)}
                required
                minLength={4}
                maxLength={20}
            />
            <Form.Control.Feedback type="invalid">Username length must be between 4 and 20 characters</Form.Control.Feedback>
        </Form.Group>
        <Form.Group className="mb-3" controlId="formPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control 
                type="password" 
                placeholder="Password" 
                onChange={e => setPassword(e.target.value)}
                required
                minLength={8}
                maxLength={20}
            />
            <Form.Control.Feedback type="invalid">Password length must be between 8 and 20 characters</Form.Control.Feedback>
        </Form.Group>
        <Button variant="primary" type="submit" disabled={!canSubmit()}>
            Login
        </Button>
    </Form>
}

export default Login