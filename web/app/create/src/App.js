import React from "react";                                                                                                          1
import {validateSession} from '../../shared/src/Session.js'
import Header from "../../shared/src/Header.js";
import Create from "./Create.js";
import Footer from "../../shared/src/Footer.js";


export default class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = {sessionKey: '', loggedin: false, email: ''};
		validateSession(this);
	}

    render() {
        return (
			<div className="window">
				<Header Loggedin={this.state.loggedin} Email={this.state.email} SessionKey={this.state.sessionKey} />
				<Create Loggedin={this.state.loggedin} SessionKey={this.state.sessionKey} />
				<Footer />
			</div>
        );
    }
}
