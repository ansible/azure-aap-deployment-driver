import React from 'react';
import { Bullseye, Button, Modal, ModalVariant, StackItem } from '@patternfly/react-core';
import { cancelDeployment } from '../../apis/deployment';

export const CancelDeployment = ({setCancelled}) => {

  const [isModalOpen, setIsModalOpen] = React.useState(false);

  const handleModalToggle = () => {
    setIsModalOpen(!isModalOpen);
  };

  async function handleClick() {
    try {
      const cancelled = await cancelDeployment()
      // TODO add visual confirmation that deployment was cancelled
      console.log(`Deployment was cancelled: ${cancelled}`);
      setCancelled(true);
      document.getElementsByClassName("cancelButton")[0].remove();
      document.getElementsByClassName("retryButton")[0]?.remove();
      if(document.getElementsByClassName("infoText")){document.getElementsByClassName("infoText")[0].innerHTML = 
      "Your Ansible on Azure deployment is cancelled. You still need to delete the managed application from your Azure subscription."+
        "In the Azure Portal, navigate to 'Resource Groups', and then to the resource group where you deployed the instance of the managed application. "+
        "Select the managed application from the list of resources and then click 'Delete' to remove all resources associated with the managed application"}
      setIsModalOpen(!isModalOpen);
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
