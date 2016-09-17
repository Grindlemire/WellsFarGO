import React from 'react';

class Log extends React.Component {
    constructor(props){
        super(props)
    }

    render () {
        const {data, pending, err} = this.props

        var columns = (
                <tr>
                    <th>Date</th>
                    <th>Location</th>
                    <th>Amount</th>
                </tr>
            )

        var rows = data.map((log) => {
            return (
                <tr key={log.id} className={(log.amount<0) ? "greenRow" : "redRow"}>
                    <td className="cell">{log.date}</td>
                    <td className="cell">{log.location}</td>
                    <td className="cell">{log.amount}</td>
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
