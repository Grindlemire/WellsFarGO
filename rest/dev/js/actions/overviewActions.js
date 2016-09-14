import axios from "axios";
// import moment from 'moment';

export function fetchOverview(entity) {
    return function(dispatch) {
        dispatch({type: "FETCH_OVERVIEW_PENDING"});

        var req = {
            method: 'GET',
            url: "http://rest.learncode.academy/api/redux/friends"
        };

        axios(req)
            .then((response) => {
                dispatch({
                    type: "FETCH_OVERVIEW_FULFILLED",
                    payload: "Data Received"
                });
            })
            .catch((err) => {

                if(err.response) {
                    var error = err
                } else {
                    error = err.message
                }
                dispatch({
                    type: "FETCH_OVERVIEW_REJECTED",
                    payload: error
                });
            });

    };
}
