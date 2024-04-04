import React from 'react';
import { ExpandableSection, Flex, FlexItem, List, ListItem, StackItem } from '@patternfly/react-core';
import { EngineConfiguration } from '../apis/types';
import './EngineConfigurationInfo.css'

interface IEngineConfiguration {
	engineConfig:EngineConfiguration
}

const timeOnlyRegex = /(\d{2}):(\d{2}):(\d{2})/	// regexp used for picking up time values from ISO date & time string

const getFormattedTime = (secondsToFormat:number):string => {
	// this only handles times smaller than 24 hours, should be enough
	var date = new Date(0);
	date.setSeconds(secondsToFormat);
	const timeString = date.toISOString()
	const matchesArray = timeOnlyRegex.exec(timeString)
	if (matchesArray != null && matchesArray.length === 4) {
  	return (matchesArray.slice(1).map((element, i) => {		// returned array has full match on index 0, hence starting at 1
      const e = parseInt(element,10)
      return e === 0 ? "" : (i===0) ? `${e} hour${e>1?"s":""} `: (i===1)? `${e} minute${e>1?"s":""} ` : `${e} second${e>1?"s":""}`
    }).join(""))
	}
	return ""
}

export const EngineConfigurationInfo = ({engineConfig}:IEngineConfiguration ) => {
	return (
		<StackItem>
			<ExpandableSection toggleText="Deployment engine configuration">
				<List isPlain isBordered className='engine-configuration pf-u-box-shadow-md'>
					<ListItem>
						<Flex>
							<FlexItem align={{default: 'alignLeft'}}>Maximum time to restart failed deployment step:</FlexItem>
							<FlexItem align={{default: 'alignRight'}}>{getFormattedTime(engineConfig.stepRestartTimeout)}</FlexItem>
						</Flex>
					</ListItem>
					<ListItem>
						<Flex>
							<FlexItem align={{default: 'alignLeft'}}>Maximum number of restarts for each failed deployment step:</FlexItem>
							<FlexItem align={{default: 'alignRight'}}>{engineConfig.stepMaxRetries}</FlexItem>
						</Flex>
					</ListItem>
					<ListItem>
						<Flex>
							<FlexItem align={{default: 'alignLeft'}}>Maximum deployment step run time:</FlexItem>
							<FlexItem align={{default: 'alignRight'}}>{getFormattedTime(engineConfig.stepDeploymentTimeout)}</FlexItem>
						</Flex>
					</ListItem>
					<ListItem>
						<Flex>
							<FlexItem align={{default: 'alignLeft'}}>Maximum total deployment engine run time:</FlexItem>
							<FlexItem align={{default: 'alignRight'}}>{getFormattedTime(engineConfig.overallTimeout)}</FlexItem>
						</Flex>
					</ListItem>
				</List>
			</ExpandableSection>
		</StackItem>
	)
}
