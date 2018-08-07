import React from "react";                                                                                                          1
import axios from "axios";

export default class Login extends React.Component {
	constructor(props) {
		super(props);
		this.state = {email: '', password: '', loginError: ''};

		this.handleChange = this.handleChange.bind(this);
		this.handleClick = this.handleClick.bind(this);
	}

	handleChange(event) {
		this.setState({loginError: ''});
		switch(event.target.id) {
			case "loginEmailInput":
				this.setState({email: event.target.value});
				break;
			case "loginPasswordInput":
				this.setState({password: event.target.value});
				break;
		}
	}

	handleClick(event) {
		var self = this;
		event.preventDefault();
		axios.post('/login', {
			register: {
				email: this.state.email,
				password: this.state.password
			}
		})
		.then(function (response) {
			switch(response.data.error) {
				case "":
					window.location = '/';
					break;
				default:
					self.setState({loginError: response.data.error});
					break;
			}
		})
		.catch(function (error) {
			console.log(error);
		});
	}

    render() {
		const isLoggedIn = this.props.Loggedin;
		if (isLoggedIn) {
			return (
				<div className="container-fluid">
					<h3>You are already logged in!</h3>
				</div>
			);
		} else {
			return (
				<div className="container-fluid">
					<div id="loginContainer">
						<div id="loginHeader">
							<h3>Login to your account</h3>
						</div>
						<form className="form-inline">
							<div className="loginRow">
								<div className="form-group loginDiv">
									<label className="sr-only" for="loginEmailInput">
										Email input
									</label>
									<input type="text" id="loginEmailInput" className="form-control loginInput" placeholder="Enter email address" onChange={this.handleChange} />
								</div>
							</div>
							<div className="loginRow">
								<div className="form-group loginDiv">
									<label className="sr-only" for="loginPasswordInput">
										Passowrd input
									</label>
									<input type="password" id="loginPasswordInput" className="form-control loginInput" placeholder="Enter password" onChange={this.handleChange} />
								</div>
								<div id="loginError" className="error loginError">
									{this.state.loginError}
								</div>
							</div>
							<div id="loginSubmitDiv" className="rightFloat">
								<button type="submit" className="btn btn-default" onClick={this.handleClick}>
									Login
								</button>
							</div>
							<div className="clear"></div>
						</form>
					</div>
				</div>
			);
		}
    }
}
