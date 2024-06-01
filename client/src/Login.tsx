import { FormEvent, useState } from "react"
import { Button, Container, Form } from "react-bootstrap"
import axios from "./api/axios";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router-dom";
import { AxiosError, isAxiosError } from "axios";
import { login } from "./auth/login";

const Login = () => {
    // TODO add input checks
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [_1, setCookie, _2] = useCookies();
    const [validate, setValidate] = useState(false);
    const [processing, setProcessing] = useState(false);
    const [failedMsg, setFailedMsg] = useState('');

    const navigate = useNavigate();

    const canSubmit = (): boolean => {
        return !processing
    }

    const onSubmit = async (e: FormEvent) => {
        e.preventDefault();
        if (!canSubmit()) return;

        setProcessing(true);

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
            setProcessing(false);

            if (!isAxiosError(ex)) {
                console.error(ex);
                return;
            }
            
            const err = ex as AxiosError;
            if (err.response!.status == 401) {
                const data: any = err.response?.data;
                
                setFailedMsg(data);
                return;
            }

            // TODO add other status code handlers
            console.error(err);
        }
        
    }

    return (
    <Container>
        <h3>Login</h3>
        <Form noValidate validated={validate} onSubmit={onSubmit}>
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
                Login
            </Button>
        </Form>
    </Container>
    )
}

export default Login