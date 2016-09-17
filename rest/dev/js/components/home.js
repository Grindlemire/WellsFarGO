import React from 'react';
import 'styles/index';



class Home extends React.Component {
    constructor(props){
        super(props)
    }

    render () {
        return(
            <div>
                {this.props.children}
            </div>

        )
    }
}

export default Home
