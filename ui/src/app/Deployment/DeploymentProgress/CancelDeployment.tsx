import React from 'react';
import { Button, Checkbox, Modal, ModalVariant } from '@patternfly/react-core';
import MinusCircleIcon from '@patternfly/react-icons/dist/esm/icons/minus-circle-icon';
import { cancelDeployment } from '../../apis/deployment';

export const CancelDeployment = () => {

  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [isConfirmed, setIsConfirmed] = React.useState(false);

  const handleModalToggle = () => {
    setIsModalOpen(!isModalOpen);
    //reset confirmation on modal close
    setIsConfirmed(false)
  };

  const handleConfirmClick = () => {
    setIsConfirmed(!isConfirmed)
  }

  async function handleCancelClick() {
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
    <>
      <Button className='cancelButton' variant="secondary" onClick={() => handleModalToggle()}>Cancel Deployment</Button>
      <Modal
        variant={ModalVariant.small}
        title="Cancel Deployment"
        titleIconVariant={'warning'}
        isOpen={isModalOpen}
        onClose={handleModalToggle}
        actions={[
          <Button key="confirm" variant="danger" onClick={handleCancelClick} icon={<MinusCircleIcon />} isDisabled={!isConfirmed}>
            Cancel deployment
          </Button>,
          <Button key="cancel" variant="link" onClick={handleModalToggle}>
            Cancel
          </Button>
        ]}>
        Are you sure you want to cancel your deployment?
        <br/>
        <br/>
        <br/>
        <Checkbox label="Yes, I confirm that I want to cancel this deployment." id="cancel-confirm" onClick={handleConfirmClick}/>
        <br/>
      </Modal>
    </>
  )
}
