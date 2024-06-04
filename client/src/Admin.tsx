import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "./api/axios";
import { isAxiosError } from "axios";
import { Container } from "react-bootstrap";

// ! don't know how safe this is

const Admin = () => {
    const navigate = useNavigate();
    const [isAdmin, setIsAdmin] = useState(false);

    useEffect(() => {
        (async () => {
            try {
                const resp = await axios.get('/user', {withCredentials: true});
                const data = resp.data as PrivateUserData;
                if (!data.isAdmin) {
                    navigate('/');
                    return;
                }
                setIsAdmin(true);
            } catch (ex) {
                if (!isAxiosError(ex)) {
                    console.error(ex);
                    return;
                }

                const status = ex.response?.status;
                if (status == 401) {
                    navigate('/login');
                    return;
                }

                console.error(ex);
            }
            // TODO handle error
        })();
    }, []);
    return (<Container>
        {isAdmin && (
            <h2>Admin board</h2>
        )}
    </Container>);
};

export default Admin;