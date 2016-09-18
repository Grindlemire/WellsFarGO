import axios from "axios";
import moment from 'moment';

export function fetchLog(entity) {
    return function(dispatch) {
        dispatch({type: "FETCH_LOG_PENDING"});

        var start = moment().subtract(1, 'years').format("MM/DD/YYYY");
        var end = moment().format("MM/DD/YYYY");

        var range = {
            start: start,
            end: end
        };

        var req = {
            method: 'POST',
            url: "transactions/range",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            data: JSON.stringify(range)
        };

        axios(req)
            .then((response) => {
                dispatch({
                    type: "FETCH_LOG_FULFILLED",
                    payload: response.data
                });
            })
            .catch((err) => {

                if(err.response) {
                    var error = err
                } else {
                    error = err.message
                }
                dispatch({
                    type: "FETCH_LOG_REJECTED",
                    payload: error
                });
            });

    };
}
