import React from 'react';
import Menu from "Components/navigation/menu";

export default class Root extends React.Component {
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