import { Container, Nav, Navbar } from "react-bootstrap";
import {
    Link
} from 'react-router-dom'


const Header = () => {

    return <Navbar expand="lg" className="bg-body-tertiary">
        <Container>
            <Navbar.Brand as={Link} to="/">Card Store</Navbar.Brand>
            <Navbar.Toggle aria-controls="navbar-elements" />
            <Navbar.Collapse id="navbar-elements">
            <Nav className="me-auto">
                <Nav.Link as={Link} to="/">Home</Nav.Link>
                <Nav.Link as={Link} to="/cards">Cards</Nav.Link>
                <Nav.Link as={Link} to="/about">About</Nav.Link>
            </Nav>
            </Navbar.Collapse>
        </Container>
    </Navbar>
    
}

export default Header;