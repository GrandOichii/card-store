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


const App = () => {

    return <BrowserRouter>
        <Header />
        <Routes>
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
            <Route 
                path='/register' 
                element={<Register />}
            />
            {/* TODO: register */}
            {/* TODO: collections */}
            {/* TODO: admin */}
        </Routes>
        <Footer />
    </BrowserRouter>
}

export default App