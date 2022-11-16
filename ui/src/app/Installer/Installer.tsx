/* eslint-disable @typescript-eslint/no-unused-vars */
import * as React from 'react';
import { RestartDeployment } from './RestartDeployment';
import {
  Button,
  Card,
  CardBody,
  CardTitle,
  Gallery,
  PageSection,
  PageSectionVariants,
  Text,
  TextContent,
  Title,
  List,
  ListItem,
  Icon,
  TextVariants,
  Tooltip,
} from '@patternfly/react-core';
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { useEffect, useState } from 'react';
import { ProgressBar } from './ProgressBar';

interface Repository {
  name: string;
  branches: string | null;
  prs: string | null;
  workspaces: string;
  lastCommit: string;
}
const options = {
  method : 'POST',
  headers : {
    'Content-Type' : 'application/json'
  }
}
const Installer: React.FunctionComponent = () =>
{
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  var [percent, setPercent] = useState(0)
  var dataLength = 0

  useEffect(() => {
    fetch('http://localhost:9090/step', {
    }).then(response => {
      if(response.ok) {
        return response.json()
      }
      throw response;
    })
    .then(data => {
      setData(data);
    })
    .catch(error => {
      console.error("Error fetching deployment details", error);
      setError(error);
    })
    .finally(() => {
      setLoading(false);
    })
  }, []);

  dataLength = data.length

  const ProgressChangeHandler = (data) => {
    percent = percent + 100/dataLength;
    return <Text className='timeTaken' component={TextVariants.h5}>({data['executions'][0]['duration']})</Text>
  };

  return (
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
        }}
      >
        <div>
        <Card isHoverable isCompact style={{width:"203%"}}>
          <CardTitle>
            <Title headingLevel="h2" size="md">
              Deployment Steps
            </Title>
          </CardTitle>
          <CardBody className="cardbody">
        <List isPlain isBordered >
        {data.map(data => (
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
        ))}
        </List>
          </CardBody>
        </Card>
        </div>
        <br></br>
        <ProgressBar data = {data} data1 = {percent}></ProgressBar>
      </Gallery>
    </PageSection>
    <PageSection>
    </PageSection>
  </>
)};

export { Installer };
  function jsonData(jsonData: any): BodyInit | null | undefined {
    throw new Error('Function not implemented.');
  }

