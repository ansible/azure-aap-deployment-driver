import { ActionTypes } from "../constants/action-types";

const initialState = {
    deploymentSteps: [],
    err: null
}
export const deploymentReducer = (state = initialState, {type, payload, error}) => {
    switch(type) {
        case ActionTypes.DEPLOYMENT_STEPS:
            return {...state, deploymentSteps: payload, err: error};
        default:
            return state; 
    }
}
