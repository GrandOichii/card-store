import {
    Routes,
    Route,
    BrowserRouter
} from 'react-router-dom'

import Header from "./Header"
import About from './About'
import Footer from './Footer'
import Cards from './Cards'
import Login from './Login'
import Register from './Register'
import LargeCardDisplay from './LargeCardDisplay'
import Collections from './Collections'
import LargeCollectionDisplay from './LargeCollectionDisplay'
import { useEffect, useState } from 'react'
import axios from './api/axios'
import { ExpansionsContext, LanguagesContext } from './context'


const App = () => {
    const [languages, setLanguages] = useState([]);
    const [expansions, setExpansions] = useState([]);

    useEffect(() => {
        console.log('FETCH');
        [
            {
                'path': '/card/languages',
                'setter': setLanguages,
            },
            {
                'path': '/card/expansions',
                'setter': setExpansions,
            },
        ].forEach(e => {
            
            (async () => {
                // TODO catch errors
                const resp = await axios.get(e.path)
                e.setter(resp.data);
            })();
        });
    }, []);

    return(
        <LanguagesContext.Provider value={languages}>
            <ExpansionsContext.Provider value={expansions}>
                <BrowserRouter>
                    <Header />
                    <Routes>
                        <Route 
                            path='/about' 
                            element={<About />}
                        />
                        <Route 
                            path='/cards/:type/all' 
                            element={<Cards />}
                        />
                        <Route 
                            path='/cards/:id' 
                            element={<LargeCardDisplay />}
                        />
                        <Route 
                            path='/login' 
                            element={<Login />}
                        />
                        <Route 
                            path='/register' 
                            element={<Register />}
                        />
                        <Route 
                            path='/collections' 
                            element={<Collections />}
                        />
                        <Route 
                            path='/collections/:id' 
                            element={<LargeCollectionDisplay />}
                        />
                        {/* TODO: admin */}
                    </Routes>
                    <Footer />
                </BrowserRouter>
            </ExpansionsContext.Provider>
        </LanguagesContext.Provider>
    ) 
}

export default App