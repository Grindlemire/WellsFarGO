import React from 'react';
// import moment from 'moment';

import { connect } from "react-redux";

import { fetchOverview } from 'js/actions/overviewActions';



class OverviewComponent extends React.Component {
    constructor(props){
        super(props)
    }

    componentWillMount() {
        this.props.dispatch(fetchOverview());
    }


    render () {
        const {data, err, pending} = this.props

        if(err !== null) {
            var errElem = <div class="error">Error: {err}</div>
        }
        return(
            <div className="row">
                <div className="col s12">
                    <div>Hello World: {data}</div>
                    {pending && (() => <div>PENDING</div>)()}
                    {errElem}
                </div>
            </div>
        )
    }
}

const mapStateToProps = (store) => {
    return {
        data: store.overview.data,
        err: store.overview.err,
        pending: store.overview.pending
    };
}


const Overview = connect(
    mapStateToProps
)(OverviewComponent)


export default Overview
