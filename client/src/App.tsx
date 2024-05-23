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


const App = () => {

    return <BrowserRouter>
        <Header />
        <Routes>
            {/* TODO: / */}
            <Route 
                path='/about' 
                element={<About />}
            />
            <Route 
                path='/cards/:type' 
                element={<Cards />}
            />
            
            <Route 
                path='/login' 
                element={<Login />}
            />

            {/* TODO: register */}
            {/* TODO: collections */}
            {/* TODO: admin */}
        </Routes>
        <Footer />
    </BrowserRouter>
}

export default App