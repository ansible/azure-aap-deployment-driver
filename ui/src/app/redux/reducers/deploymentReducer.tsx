import { ActionTypes } from "../contants/action-types";

const initialState = {
    products: [],
}
export const deploymentReducer = (state = initialState, {type, payload}) => {
    switch(type) {
        case ActionTypes.DEPLOYMENT_STEPS:
            return {...state, products: payload};
        default:
            return state; 
    }
}