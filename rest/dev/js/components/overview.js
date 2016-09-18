import React from 'react';
import { connect } from "react-redux";
import { fetchLog, sortClick } from 'js/actions/logActions';
import Log from 'js/components/log';



class OverviewComponent extends React.Component {
    constructor(props){
        super(props)

        this.sortClick = this.sortClick.bind(this)
    }

    componentWillMount() {
        this.props.dispatch(fetchLog());
    }

    sortClick(field, asc) {
        this.props.dispatch({
            type: "CLICK_LOG_SORT",
            payload: {
                order: field,
                asc: asc
            }
        });
    }


    render () {
        const {logData} = this.props
        return(
            <Log
                data={logData.data}
                pending={logData.pending}
                order={logData.order}
                asc={logData.asc}
                sortClick={this.sortClick}
                err={logData.err}
                />
        )
    }
}

const mapStateToProps = (store) => {
    return {
        logData: store.log
    };
}


const Overview = connect(
    mapStateToProps
)(OverviewComponent)


export default Overview
