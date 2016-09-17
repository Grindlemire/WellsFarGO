import React from 'react';
import { connect } from "react-redux";
import { fetchLog } from 'js/actions/logActions';
import Log from 'js/components/log';



class OverviewComponent extends React.Component {
    constructor(props){
        super(props)
    }

    componentWillMount() {
        this.props.dispatch(fetchLog());
    }


    render () {
        const {logData} = this.props

        return(
            <Log
                data={logData.data}
                pending={logData.pending}
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
