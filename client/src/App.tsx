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


const App = () => {

    return <BrowserRouter>
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
            {/* TODO: register */}
            {/* TODO: admin */}
        </Routes>
        <Footer />
    </BrowserRouter>
}

export default App