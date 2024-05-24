import {AxiosInstance} from 'axios'
import { NavigateFunction } from 'react-router-dom';
import { CookieSetOptions } from 'universal-cookie';


const login = async (
        api: AxiosInstance, 
        loginData: LoginData, 
        setCookie: (name: string, value: any, options?: CookieSetOptions | undefined) => void,
        navigate: NavigateFunction
 ) => {
    await api.post('/auth/login', loginData, {
        withCredentials: true,
    });
    setCookie('loggedIn', true, {
        maxAge: 3600
    })
    navigate("/collections")
}

export default login;