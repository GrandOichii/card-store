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
                        <NavDropdown.Item href="/cards/MTG">Magic: the Gathering</NavDropdown.Item>
                        <NavDropdown.Item href="/cards/YGO">
                            Yu-Gi-Oh!
                        </NavDropdown.Item>
                    </NavDropdown>
                    <Nav.Link as={Link} to={loggedIn() ? "/collections" : "/login"}>Collection</Nav.Link>
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