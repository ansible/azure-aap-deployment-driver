import React from 'react';
import { Bullseye, Button, StackItem } from '@patternfly/react-core';
import { restartStep } from '../../apis/deployment';

interface IRestartDeploymentProps {
  stepExId: number
}
interface LoadingPropsType {
  spinnerAriaValueText: string;
  spinnerAriaLabelledBy?: string;
  spinnerAriaLabel?: string;
  isLoading: boolean;
}

export const RestartStep = ({ stepExId }: IRestartDeploymentProps) => {

  const [isPrimaryLoading, setIsPrimaryLoading] = React.useState<boolean>(false);
  const primaryLoadingProps = {} as LoadingPropsType;
  primaryLoadingProps.spinnerAriaValueText = 'Loading';
  primaryLoadingProps.spinnerAriaLabelledBy = 'primary-loading-button';
  primaryLoadingProps.isLoading = isPrimaryLoading;

  async function handleRestart(event) {
    event.preventDefault();
    try {
      setIsPrimaryLoading(!isPrimaryLoading)
      await restartStep(stepExId);
    } catch (error) {
      console.log(error)
    }
    setTimeout(()=>{
      setIsPrimaryLoading(false)
    }, 10000)
  }

  return (
    <StackItem>
      <Bullseye>
        <Button className='retryButton' id="primary-loading-button" variant="primary" onClick={handleRestart} {...primaryLoadingProps}>
          {isPrimaryLoading ? 'Restarting Step' : 'Restart Step'}</Button>
      </Bullseye>
    </StackItem>
  );
};
