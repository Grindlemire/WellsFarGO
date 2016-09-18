import React from 'react';

class Log extends React.Component {
    constructor(props){
        super(props)
        this.curOrder = ""
        this.sign = 1;

    }

    componentDidMount() {

    }


    compare(a,b, field, asc) {
        if (a[field] < b[field]){
            if(asc) {
                return -1
            } else {
                return 1
            }
        }

        if (a[field] > b[field]){
            if (asc) {
                return 1
            } else {
                return -1
            }
        }
        return 0;


    }

    clickField(field) {
        if (field == this.props.order) {
            this.props.sortClick(field, !(this.props.asc));
        } else {
            this.props.sortClick(field, true);
        }

    }

    sort(data, field, asc) {
        data.sort((a, b) => {
            return this.compare(a, b, field, asc)
        });
        return data
    }

    render () {
        const {data, order, asc, pending, err} = this.props
        var sortedData = this.sort(data, order, asc)

        var columns = (
                <tr>
                    <th onClick={() => this.clickField("date")}>Date</th>
                    <th onClick={() => this.clickField("location")}>Location</th>
                    <th onClick={() => this.clickField("amount")}>Amount</th>
                </tr>
            )

        var rows = data.map((log) => {
            return (
                <tr key={log.id} className={(log.amount<0) ? "greenRow" : "redRow"}>
                    <td className="date">{log.date}</td>
                    <td className="location">{log.location}</td>
                    <td className="amount">{log.amount}</td>
                </tr>
            )
        })


        return(
            <div className="row">
                <div className="col s12">
                    <table>
                        <thead>
                            { columns }
                        </thead>
                        <tbody>
                            { rows }
                        </tbody>
                    </table>
                </div>
            </div>
        )
    }
}

export default Log
