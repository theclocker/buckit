import React from 'react';
import ReactDOM from 'react-dom';
import Menu from "Components/navigation/menu";
import 'Src/style.scss';

class Index extends React.Component {
    render() {
        return (
            <div>
                <div className="columns">
                    <div className="column is-one-third is-hidden-touch">
                        <Menu/>
                    </div>
                    <div className="column">
                        123
                    </div>
                </div>
            </div>
        );
    }
}

ReactDOM.render(<Index/>, document.getElementById('app'));