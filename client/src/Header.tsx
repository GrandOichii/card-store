import { Container, Nav, NavDropdown, Navbar } from "react-bootstrap";
import { useCookies } from "react-cookie";
import {
    Link
} from 'react-router-dom'

const Header = () => {
    const [cookies, _1, _2] = useCookies()
    const loggedIn = () => cookies['loggedIn'];    

    return <Navbar expand="lg" className="bg-body-tertiary">
        <Container>
            <Navbar.Brand as={Link} to="/">Card Store</Navbar.Brand>
            <Navbar.Toggle aria-controls="navbar-elements" />
            <Navbar.Collapse id="navbar-elements">
                <Nav className="me-auto">
                    <Nav.Link as={Link} to="/">Home</Nav.Link>
                    <NavDropdown title="Cards" id="basic-nav-dropdown">
                        {/* FIXME forcing these to be as={Link} results in not being able to navigate between them */}
                        <NavDropdown.Item href="/cards/MTG/all">Magic: the Gathering</NavDropdown.Item>
                        <NavDropdown.Item href="/cards/YGO/all">Yu-Gi-Oh!</NavDropdown.Item>
                    </NavDropdown>
                    {/* TODO move this to the Collections component */}
                    <Nav.Link as={Link} to={loggedIn() ? "/collections" : "/login"}>Collections</Nav.Link>
                    <Nav.Link as={Link} to="/about">About</Nav.Link>
                </Nav>
                {!loggedIn() && 
                    <Nav>
                        <Nav.Link as={Link} to='/login'>Login</Nav.Link>
                        <Nav.Link as={Link} to='/register'>Register</Nav.Link>
                    </Nav>
                }
            </Navbar.Collapse>
        </Container>
    </Navbar>
    
}

export default Header;