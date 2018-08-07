import React from "react";                                                                                                          1
import axios from "axios";

function renderPasswordError() {
	return (
		<div>
			Password needs to meet the following requirements:
			<ul>
				<li>Be between 8 and 16 characters long</li>
				<li>Contain at least one uppercase character</li>
				<li>Contain atleast one lowercase character</li>
				<li>Contain at least one digit</li>
				<li>Contain at least one special character</li>
			</ul>
		</div>
	);
}

function RenderInput(props) {
	var errorMessage = props.Error;
	if(props.Name == "Password" && props.Error == true) {
		errorMessage = renderPasswordError();
	}
	return (
		<div>
			<div className="registerRow">
				<div id={'register'+props.Name+'Div'} className="registerDiv form-group">
					<label className="sr-only" for={'register'+props.Name+'Input'}>
						{props.Name} input
					</label>
					<input type={props.Type} id={'register'+props.Name+'Input'} className="form-control registerInput" placeholder={props.Placeholder} onChange={props.Obj.handleChange} />
				</div>
			</div>
			<div id={'register'+props.Name+'Error'} className="error registerError">
				{errorMessage}
			</div>
		</div>
	);
}

function occurrences(string, subString, allowOverlapping) {
	string += "";
	subString += "";
	if (subString.length <= 0) return (string.length + 1);
	var n = 0,
		pos = 0,
		step = allowOverlapping ? 1 : subString.length;
	while (true) {
		pos = string.indexOf(subString, pos);
		if (pos >= 0) {
			++n;
			pos += step;
		} else break;
	}
	return n;
}

export default class Register extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			email: '',
			password: '',
			confirmPW: '',
			registerError: '',
			emailError: '',
			passwordError: false,
			confirmError: '',
			registered: false
		};

		this.handleChange = this.handleChange.bind(this);
		this.handleClick = this.handleClick.bind(this);
	}

	handleClick(event) {
		event.preventDefault();
		axios.post('/api/register', {
			register: {
				email: this.state.email,
				password: this.state.password,
				confirm_password: this.state.confirmPW
			}
		})
		.then((response) => {
			switch(response.status) {
				case 200:
					switch(response.data["error"]) {
						case "":
							this.setState({registered: true});
							break;
						default:
							this.setState({registerError: response.data["error"]});
							break;
					}
					break;
				default:
					this.setState({registerError: 'Unexpected error. Please try again later'});
					break;
				}
		})
		.catch(function (error) {
			console.log(error);
		});
	}

	validateEmail(email) {
		if(email == "") {
			return
		}
		if(occurrences(email, "@") != 1) {
			this.setState({emailError: 'Invalid email address'});
			return
		}
		var vals = email.split("@");
		if(vals.length != 2 || vals[0] == "" || vals[1] == "") {
			this.setState({emailError: 'Invalid email address'});
			return
		}
	}

	validatePassword(password) {
		if(password == "") {
			return
		}
		var upperCase = false
		var lowerCase = false
		var digit = false
		var special = false
		var i = password.length;
		var specialChars = /[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;
		while (i--) {
			if(specialChars.test(password.charAt(i))) {
				special = true;
				continue;
			}
			if(password.charAt(i) == password.charAt(i).toUpperCase() && isNaN(parseInt(password.charAt(i), 10))) {
				upperCase = true;
			}
			if(password.charAt(i) == password.charAt(i).toLowerCase() && isNaN(parseInt(password.charAt(i), 10))) {
				lowerCase = true;
			}
			if(!isNaN(parseInt(password.charAt(i), 10))) {
				digit = true;
			}
		}
		if(!upperCase || !lowerCase || !digit || !special || password.length <8 || password.length > 16) {
			this.setState({passwordError: true});
		}
	}

	handleChange(event) {
		this.setState({registerError: ''});
		switch(event.target.id) {
			case "registerEmailInput":
				this.setState({email: event.target.value, emailError: ''});
				this.validateEmail(event.target.value)
				break;
			case "registerPasswordInput":
				this.setState({password: event.target.value, passwordError: false});
				this.validatePassword(event.target.value)
				break;
			case "registerConfirmInput":
				var confirmError = '';
				if(this.state.password != event.target.value) {
					confirmError = 'Passwords don\'t match!';
				}
				this.setState({confirmPW: event.target.value, confirmError: confirmError});

				break;
		}
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
			if(this.state.registered) {
				return (
					<div className="container-fluid">
						<h3>
							Account registered. You can now <a href="/login">log in!</a> (this needs to be updated with a confirmation email etc)
						</h3>
					</div>
				);
			} else {
				const formCorrect = (this.state.emailError == "" && this.state.passwordError == "" && this.state.confirmError == "" && this.state.email != "" && this.state.password != "" && this.state.confirmPW != "");
				return (
					<div className="container-fluid">
						<div id="registerContainer">
							<div id="registerHeader">
								<h3>
									Create new account
								</h3>
							</div>
							<form className="form-inline">
								<RenderInput Name="Email" Type="text" Placeholder="Enter email address" Error={this.state.emailError} Obj={this} />
								<RenderInput Name="Password" Type="password" Placeholder="Enter password" Error={this.state.passwordError} Obj={this} />
								<RenderInput Name="Confirm" Type="password" Placeholder="Confirm password" Error={this.state.confirmError} Obj={this} />
								<div className="leftFloat error" id="registerError">
									{this.state.registerError}
								</div>
								<div id="registerSubmitDiv" className="rightFloat">
									{formCorrect ? (
										<button type="submit" className="btn btn-default" onClick={this.handleClick}>
											Register
										</button>
									) : (
										<button type="submit" className="btn btn-default" disabled>
											Register
										</button>
									)}
								</div>
								<div className="clear"></div>
							</form>
						</div>
					</div>
				);
			}
		}
	}
}
