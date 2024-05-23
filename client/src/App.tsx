import {
    Routes,
    Route,
    BrowserRouter
} from 'react-router-dom'

import Header from "./Header"
import About from './About'
import Footer from './Footer'
import Cards from './Cards'


const App = () => {

    return <BrowserRouter>
        <Header />
        <Routes>
            <Route 
                path='/about' 
                element={<About />}
            />
            <Route 
                path='/cards' 
                element={<Cards />}
            />
        </Routes>
        <Footer />
    </BrowserRouter>
}

export default App