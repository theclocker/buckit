import React from 'react';

export default class Menu extends React.Component {
    render() {
        return (
            <aside className="menu main-menu">
                <p className="menu-label">
                    Apis
                </p>
                <ul className="menu-list">
                    <li><a>Real-time dashboard</a></li>
                    <li>
                        <a>Manage your APIs</a>
                        <ul>
                            <li><a>Add new API</a></li>
                            <li><a>Manage credentials</a></li>
                            <li><a>Mapped endpoints</a></li>
                            <li><a>Statistics</a></li>
                            <li><a>Queues</a></li>
                            <li><a>Logs</a></li>
                        </ul>
                    </li>
                </ul>
                <p className="menu-label">
                    User management
                </p>
                <ul className="menu-list">
                    <li>
                        <a>Manage your users</a>
                        <ul>
                            <li><a>Add new user</a></li>
                            <li><a>User groups</a></li>
                            <li><a>Users access</a></li>
                            <li><a>Access levels</a></li>
                        </ul>
                    </li>
                </ul>
                <p className="menu-label">
                    Cluster management
                </p>
                <ul className="menu-list">
                    <li>
                        <a>Manage your servers</a>
                        <ul>
                            <li><a>Add new server</a></li>
                            <li><a>Edit users permissions</a></li>
                            <li><a>Cluster permissions</a></li>
                        </ul>
                    </li>
                </ul>
            </aside>
        );
    }
}