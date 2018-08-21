import React from 'react';
import Menu from "Components/navigation/menu";

export default class Topnav extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            show: false
        };
        this.toggleShow = this.toggleShow.bind(this);
    }

    toggleShow() {
        this.setState((prevState) => ({
            show: !prevState.show
        }))
    }

    render() {
        return (
            <div>
                <nav className="navbar has-shadow" role="navigation" aria-label="main navigation">
                    <div className="navbar-brand">
                        <a className="navbar-item">
                            <img src="https://bulma.io/images/bulma-logo.png"
                                 alt="Bulma: a modern CSS framework based on Flexbox" width="112" height="28"/>
                        </a>
                        <a role="button" onClick={this.toggleShow} className="navbar-burger" aria-label="menu" aria-expanded="false">
                            <span aria-hidden="true"/>
                            <span aria-hidden="true"/>
                            <span aria-hidden="true"/>
                        </a>
                    </div>
                </nav>
                {this.state.show &&
                    <div className="navbar-menu is-active">
                        <Menu/>
                    </div>
                }
            </div>
        );
    }
}