import * as React from 'react';
import {
   Card, CardBody, CardTitle, Gallery, Icon, List, ListItem, PageSection, PageSectionVariants, TextContent, Text, Title, Tooltip, TextVariants,
  } from '@patternfly/react-core';
import { RestartDeployment } from './RestartDeployment';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { useSelector } from 'react-redux';
import { RootState } from '@app/redux/reducers';
import { ProgressBar } from './ProgressBar';

export const DeploymentSteps = () => {
{
    const deploymentSteps = useSelector((state: RootState) => state.allProducts.products);
    const error = useSelector((state: RootState) => state.allProducts.err);
    var [percent, setPercent] = React.useState(0)
    var dataLength = 0
    if(deploymentSteps) {
      dataLength = deploymentSteps.length
    }
    const ProgressChangeHandler = (data) => {
        percent = percent + 100/dataLength;
        return <Text className='timeTaken' component={TextVariants.h5}>({data['executions'][0]['duration']})</Text>
      };

    return(
        <>
        <PageSection variant={PageSectionVariants.light}>
          <TextContent>
            <Text component="h1">Ansible Automation Platform Installer</Text>
          </TextContent>
        </PageSection>
        <PageSection>
          <Gallery
            hasGutter
            maxWidths={{
              sm: '100%',
              lg: '49%',
            }}>
        <div>
        <Card isHoverable isCompact style={{width:"203%"}}>
          <CardTitle>
            <Title headingLevel="h2" size="md">
              Deployment Steps
            </Title>
          </CardTitle>
          <CardBody className="cardbody">
        <List isPlain isBordered >
        {error != null? <Text component="h1">Unable to reach servers</Text> : <>
        {deploymentSteps?.map(data => (
          <ListItem className='service-list'>
          <Text style={{'marginRight':'auto'}}>{data['name']}</Text>
          {data['executions'].length ? <div className='deployment-info'>
            {data['executions'][data['executions'].length-1]['provisioningState'] === 'Succeeded' ? ProgressChangeHandler(data): <></>}
            {data['executions'][data['executions'].length-1]['provisioningState'] === 'Succeeded' ? <Tooltip
      content={
        <div>Success</div>
      }
    ><Icon className='icon1' status="success">
              <CheckCircleIcon/>
            </Icon></Tooltip>: <Tooltip
      content={
        <div>{data['executions'][data['executions'].length-1]['message']}</div>
      }
    ><Icon className='icon1' status="warning">
                <ExclamationCircleIcon/>
              </Icon></Tooltip>}
            {data['executions'][data['executions'].length-1]['provisioningState'] === 'Failed' ? <Text className='attempt' component={TextVariants.h5}>({data['executions'].length} attempts)</Text> : <></>}
            <RestartDeployment data = {data}></RestartDeployment>          
          </div>: <></>}
        </ListItem>
        ))} </> }
        </List>
          </CardBody>
        </Card>
        </div>
        <br></br>
        <ProgressBar data = {deploymentSteps} data1 = {percent}></ProgressBar>
      </Gallery>
    </PageSection>
    <PageSection>
    </PageSection>
    </>
    )
}};
