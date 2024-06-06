import { ComponentProps } from "react";
import { Modal } from "react-bootstrap";
import { Button } from "react-bootstrap";


interface FailedToAddToCollectionModalProps extends ComponentProps<"div"> {
    onHide: () => void,
    show: boolean,
    cardName: string,
    collectionName: string
}

const FailedToAddToCollectionModal = (props: FailedToAddToCollectionModalProps) => {
    return (
        <Modal
            {...props}
            size="lg"
            aria-labelledby='contained-modal-title-vcenter'
            centered
        >
            <Modal.Header>
                <h4>Failed to add card</h4>
            </Modal.Header>
            <Modal.Body>
                {`Failed to add card ${props.cardName} to collection ${props.collectionName}!`}
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={props.onHide}>Close</Button>
            </Modal.Footer>
        </Modal>

    );
};

export default FailedToAddToCollectionModal;