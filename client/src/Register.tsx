import { FormEvent, useState } from "react"
import { Button, Container } from "react-bootstrap"
import axios from "./api/axios";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";
import { AxiosError, isAxiosError } from "axios";
import { Form } from "react-bootstrap";
import { login } from "./auth/login";

const Register = () => {
    const [username, setUsername] = useState<string>('');
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [failedMsg, setFailedMsg] = useState('')
    const [validate, setValidate] = useState(false)
    const [_1, setCookie, _2] = useCookies();
    const [processing, setProcessing] = useState(false);
    const navigate = useNavigate()

    const canSubmit = (): boolean => {
        return !processing;
    }

    const onSubmit = async (e: FormEvent) => {
        e.preventDefault();
        if (!canSubmit()) return;
        setProcessing(true);

        const form = e.currentTarget as HTMLFormElement
        if (!form.checkValidity()) {
            e.stopPropagation();
            setValidate(true);
            setProcessing(false);
            return;
        }

        const registerData = {
            'username': username,
            'email': email,
            'password': password,
        };
        const loginData: LoginData = {
            'username': username,
            'password': password,
        };
        try {
            await axios.post('/auth/register', registerData, {
                withCredentials: true,
            });
            await login(axios, loginData, setCookie, navigate);
        } catch (e: any) {
            setProcessing(false);
            if (!isAxiosError(e)) {
                console.error(e);
                return;
            }
            
            const err = e as AxiosError;
            if (err.response!.status == 400) {
                const data: any = err.response?.data;
                setFailedMsg(`Failed to register: ${data}`);
                return;
            }

            // TODO add other status code handlers
            console.error(err);
        }
    }

    return (
    <Container>
        <h3>Register a new account</h3>
        <Form noValidate validated={validate} onSubmit={onSubmit}>
            <Form.Group controlId="formEmail">
                <Form.Label>Email</Form.Label>
                <Form.Control 
                    type="email" 
                    placeholder="Enter email" 
                    onChange={e => setEmail(e.target.value)}
                    required
                />
                <Form.Control.Feedback type="invalid">Invalid email</Form.Control.Feedback>
            </Form.Group>
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
            {failedMsg &&
                <div className="alert alert-danger" role='alert'>
                    {failedMsg}
                </div>
            }
            <Button variant="primary" type="submit" disabled={!canSubmit()}>
                Register
            </Button>
        </Form>
    </Container>
    )
}

export default Register