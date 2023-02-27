import React from 'react';
import { Bullseye, Button, Modal, ModalVariant, StackItem } from '@patternfly/react-core';
import { cancelDeployment } from '../../apis/deployment';

export const CancelDeployment = () => {

  const [isModalOpen, setIsModalOpen] = React.useState(false);

  const handleModalToggle = () => {
    setIsModalOpen(!isModalOpen);
  };

  async function handleClick() {
    try {
      const cancelled = await cancelDeployment()
      // TODO add visual confirmation that deployment was cancelled
      if(cancelled){
        setIsModalOpen(!isModalOpen);
      }
      // TODO add something in case cancel was not successful
    } catch (error) {
      console.log(error)
    }
  }

  return (
    <StackItem>
      <Bullseye>
        <Button className='cancelButton' variant="secondary" onClick={() => handleModalToggle()}>Cancel Deployment</Button>
        <Modal
          variant={ModalVariant.small}
          title="Cancel Deployment"
          isOpen={isModalOpen}
          onClose={handleModalToggle}
          actions={[
            <Button key="confirm" variant="primary" onClick={handleClick}>
              Confirm
            </Button>,
            <Button key="cancel" variant="link" onClick={handleModalToggle}>
              No
            </Button>
          ]}>
          Are you sure you want to cancel your deployment? If so, click the 'Confirm' Button or press the 'No' Button to return to your Deployment.
        </Modal>
      </Bullseye>
    </StackItem>
  );
};
