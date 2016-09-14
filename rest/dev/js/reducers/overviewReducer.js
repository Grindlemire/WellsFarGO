
var defaultOverview = {
    data: "Before fetch",
    pending: false,
    err: null
};


export default function reducer(state=defaultOverview , action){
    var newState = state
    switch (action.type) {
        case "FETCH_OVERVIEW_PENDING": {
            newState = Object.assign({}, state,
                {
                    pending: true
                }
            );
            break;
        }
        case "FETCH_OVERVIEW_FULFILLED": {
            newState = Object.assign({}, state,
                {
                    data: action.payload,
                    pending: false
                }
            );
            break;
        }
        case "FETCH_OVERVIEW_REJECTED": {
            newState = Object.assign({}, state,
                {
                    err: action.payload,
                    pending: false
                }
            );
            break;
        }
    }
    return newState;
}
