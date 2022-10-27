import * as React from 'react';
import {
  Bullseye,
  Button,
  Card,
  CardBody,
  CardFooter,
  CardTitle,
  Divider,
  Gallery,
  PageSection,
  PageSectionVariants,
  Split,
  SplitItem,
  ToggleGroup,
  ToggleGroupItem,
  ToggleGroupItemProps,
  Progress,
  Text,
  TextContent,
  Title,
  List,
  ListItem,
  Icon,
  TextVariants,
  Tooltip,
  Spinner,
} from '@patternfly/react-core';
// import './Alignment.css'
import CheckCircleIcon from '@patternfly/react-icons/dist/esm/icons/check-circle-icon';
import ExclamationCircleIcon from '@patternfly/react-icons/dist/esm/icons/exclamation-circle-icon';
import { ArrowRightIcon, ExternalLinkSquareAltIcon } from '@patternfly/react-icons';
import { useEffect, useState } from 'react';

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
const Dashboard: React.FunctionComponent = () => 
{
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  var [percentage, setPercentage] = useState(0);
  const service = "sampleone"
  const flag = true
  const deploymentList = [['networkingAAPDeploy','52 seconds', true] ,['dnsAAPDeploy', '2 minutes 35 seconds', true], ['databaseAAPDeploy', '45 seconds', true], ['storageAAPDeploy', '1 minute 4 seconds', false], ['kubernetesAAPDeploy', '49 seconds', true],['operatorsAAPDeploy', '3 minutes 5 seconds', true]]
  const lenghtofdep = deploymentList.length;
  var progress = 0
  var percent = 30

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

  function handleClick(id) {
    fetch('http://127.0.0.1:9090/execution/3/restart', {
      method: 'POST', 
      mode: 'cors', 
      body: JSON.stringify(jsonData)
    })
   }
  
  return (
  <>
    <PageSection variant={PageSectionVariants.light}>
      <TextContent>
        <Text component="h1">Ansible Automation Platform Installer</Text>
        {/* <Text component="p">Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</Text> */}
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
          {/* <MDBRow>
          <MDBCol md='8'>
            md="8"
          </MDBCol>
          <MDBCol md='4'>
            md="4"
          </MDBCol>
        </MDBRow> */}

        {/* 
            width: auto;
            height: 100px;
            overflow: scroll;
     */}

        <List isPlain isBordered >
        {data.map(data => (
          <ListItem className='service-list'>
          <Text style={{'marginRight':'auto'}}>{data['name']}</Text>
          {data['executions'].length ? <div className='deployment-info'>
            {data['executions'][data['executions'].length-1]['provisioningState'] === 'Succeeded' ? <Text className='timeTaken' component={TextVariants.h5}>({data['executions'][0]['duration']})</Text> : <></>}
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
            {data['executions'][data['executions'].length-1]['provisioningState'] === 'Failed' ? <Button variant="primary" onClick={handleClick}>Retry step</Button>: <></>}
          </div>: <></>}
          {/* {data['executions'].length && data['executions'][data['executions'].length-1]['provisioningState'] != 'Failed' ? <></> : <Spinner isSVG size="md" aria-label="loading.." />} */}
        </ListItem>
        ))}
          {/* <ListItem className='service-list'>
            <Text style={{'marginRight':'auto'}}>{service}</Text> 
            <div className='deployment-info'>
              <Text className='timeTaken' component={TextVariants.h5}>(51 seconds)</Text>
              {flag ? <Icon className='icon1' status="success">
                <CheckCircleIcon/>
              </Icon> : <></>}
            </div>
          </ListItem>
          <ListItem className='service-list'> 
            <Text style={{'marginRight':'auto'}}>Second</Text> 
            <div className='deployment-info'>
              <Icon className='icon1' status="warning">
                <ExclamationCircleIcon/>
              </Icon>
              <Button variant="primary">Retry step</Button>
            </div>
          </ListItem>
          <ListItem>Third</ListItem>
          <ListItem>First</ListItem>
          <ListItem>Second</ListItem>
          <ListItem>Third</ListItem> */}
        </List>
          </CardBody>
          {/* <CardFooter>
            <Bullseye>
              <Button component="a" href="/controller/" target="_blank" variant="primary">
                Go to Controller
              </Button>
            </Bullseye>
          </CardFooter> */}
        </Card>
        </div>
        <br></br>
        {/* <Card isHoverable isCompact>
          <CardTitle>
            <Title headingLevel="h2" size="md">
              Insights Analytics
            </Title>
          </CardTitle>
          <CardBody>
            <span>
              Gain insights into your deployments through visual dashboards and organization statistics, calculate your
              return on investment, and explore automation processes details.{' '}
              <Button isInline variant="link" icon={<ExternalLinkSquareAltIcon />} iconPosition="right">
                Learn more
              </Button>
            </span>
            <br />
          </CardBody>
          <CardFooter>
            <Bullseye>
              <Button component="a" href="/" target="_blank" variant="primary">
                Go to Insights Analytics
              </Button>
            </Bullseye>
          </CardFooter>
        </Card> */}
        <div>
        <Card isHoverable isCompact style={{width:"203%"}}>
          <CardTitle>
                <Progress value={percent} title="Overall progress" />
                <br></br>
                <Button className='cancleButton' variant="secondary" onClick={handleClick}>Cancel Deployment</Button>       
          </CardTitle>
          <CardBody>
            {/* <span>
              Find and use content that is supported by Red Hat and our partners to deliver reassurance for the most
              demanding environments.{' '}
              <Button isInline variant="link" icon={<ExternalLinkSquareAltIcon />} iconPosition="right">
                Learn more
              </Button>
            </span> */}
            <br />
          </CardBody>
          {/* <CardFooter>
            <Bullseye>
              <Button component="a" href="/hub/" target="_blank" variant="primary">
                Go to Automation Hub
              </Button>
            </Bullseye>
          </CardFooter> */}
        </Card>
        </div>
        {/* <Card isHoverable isCompact>
          <CardTitle>
            <Title headingLevel="h2" size="md">
              Automation Services Catalog
            </Title>
          </CardTitle>
          <CardBody>
            <span>
              Use Automation Services Catalog to collect and distribute automation content, govern your content by
              designing and attaching approval processes, and esure required sign-off is obtained by assigned
              organizational groups.{' '}
              <Button isInline variant="link" icon={<ExternalLinkSquareAltIcon />} iconPosition="right">
                Learn more
              </Button>
            </span>
            <br />
          </CardBody>
          <CardFooter>
            <Bullseye>
              <Button component="a" href="/" target="_blank" variant="primary">
                Go to Automation Services Catalog
              </Button>
            </Bullseye>
          </CardFooter>
        </Card> */}
      </Gallery>
    </PageSection>
    <PageSection>
      {/* <Card>
        <CardTitle>
          <Title headingLevel="h3" size="md">
            Additional Resources
          </Title>
        </CardTitle>
        <CardBody>
          <Split hasGutter>
            <SplitItem>
              <Title headingLevel="h4" size="md">
                Resource 1
              </Title>
              <br />
              <p>Mustache tilde tumeric everyday carry vegan 3 wolf moon palo santo.</p>
              <Button isInline variant="link" icon={<ArrowRightIcon />} iconPosition="right">
                Learn more
              </Button>
            </SplitItem>
            <Divider isVertical />
            <SplitItem>
              <Title headingLevel="h4" size="md">
                Resource 2
              </Title>
              <br />
              <p>Kitsch sriracha jean shorts humblebrag DIY pop-up echo park. Ears back wide eyed kitty.</p>
              <Button isInline variant="link" icon={<ArrowRightIcon />} iconPosition="right">
                Learn more
              </Button>
            </SplitItem>
            <Divider isVertical />
            <SplitItem>
              <Title headingLevel="h4" size="md">
                Resource 3
              </Title>
              <br />
              <p>Sleep on keyboard toy mouse squeak roll over. Refuse to drink water except out of someone’s glass.</p>
              <Button isInline variant="link" icon={<ArrowRightIcon />} iconPosition="right">
                Learn more
              </Button>
            </SplitItem>
          </Split>
        </CardBody>
      </Card> */}
    </PageSection>
  </>
)};

export { Dashboard };
  function jsonData(jsonData: any): BodyInit | null | undefined {
    throw new Error('Function not implemented.');
  }

