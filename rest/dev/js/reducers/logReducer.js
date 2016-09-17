
var defaultLog = {
    data: [],
    pending: false,
    err: null
};


export default function reducer(state=defaultLog , action){
    var newState = state
    switch (action.type) {
        case "FETCH_LOG_PENDING": {
            newState = Object.assign({}, state,
                {
                    pending: true
                }
            );
            break;
        }
        case "FETCH_LOG_FULFILLED": {
            newState = Object.assign({}, state,
                {
                    data: action.payload,
                    pending: false
                }
            );
            break;
        }
        case "FETCH_LOG_REJECTED": {
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
