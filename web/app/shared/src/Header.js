import React from "react";                                                                                                          1

function UserGreeting(props) {
  return <h1>Welcome back!</h1>;
}

function GuestGreeting(props) {
  return <h1>Please sign up.</h1>;
}

function LoginButton(props) {
  return (
    <button onClick={props.onClick}>
      Login
    </button>
  );
}

function LogoutButton(props) {
  return (
    <button onClick={props.onClick}>
      Logout
    </button>
  );
}

export default class Header extends React.Component {
	constructor(props) {
		super(props);
		 this.handleLoginClick = this.handleLoginClick.bind(this);
		 this.handleLogoutClick = this.handleLogoutClick.bind(this);
	}

	handleLoginClick() {
		this.setState({isLoggedIn: true});
	}

	handleLogoutClick() {
		this.setState({isLoggedIn: false});
	}

	render() {
		const isLoggedIn = this.props.Loggedin;
		let button;
		if (isLoggedIn) {
			button = <LogoutButton onClick={this.handleLogoutClick} />;
		} else {
			button = <LoginButton onClick={this.handleLoginClick} />
		}

        return (
            <div id="header">
                <a href="/">
                    <img src="/static/img/logo.png" id="logoImg" />
                </a>
                Get a third opinion and solve your argument once and for all!
                <div className="header-right">
                    Some buttons here {this.props.Email}
					{button}
					<div id="logout_div">
						LOGOUT
					</div>
                </div>
            </div>
        );
    }
}
