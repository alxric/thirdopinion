import React from "react";
import axios from "axios";

function HeaderLink(props) {
	return (
		<div class="leftFloat header-menu-item" onClick={props.Obj.handleClick}>
			{props.Text}
		</div>
	);
}


export default class Header extends React.Component {
	constructor(props) {
		super(props);
		this.handleClick = this.handleClick.bind(this);
	}

	handleClick(event) {
		switch(event.target.innerHTML) {
			case "Logout":
				const client = axios.create({
					headers: {
						'X-Requested-With': 'XMLHttpRequest',
						'X-CSRF-TOKEN' : this.props.SessionKey,
					},
				});

				client.post('/logout', {
				})
				.then((response) => {
					window.location.href='/';
				})
				.catch(function (error) {
					console.log(error);
				});
				break;
			case "Login":
				window.location.href='/login';
				break;
			case "Register":
				window.location.href='/register';
				break;
			case "Create":
				window.location.href='/create';
				break;
		}
	}


	render() {
		const isLoggedIn = this.props.Loggedin;
        return (
            <div id="header">
                <a href="/">
                    <img src="/static/img/logo.png" id="logoImg" />
                </a>
                Get a third opinion and solve your argument once and for all!
                <div className="header-right">
					{isLoggedIn ? (
						<div id="header-logged-in">
							<HeaderLink Obj={this} Text='Create' />
							<HeaderLink Obj={this} Text='Logout' />
						</div>
					) : (
						<div id="header-logged-out">
							<HeaderLink Obj={this} Text='Register' />
							<HeaderLink Obj={this} Text='Login' />
						</div>
					)}
                </div>
				<div class="clear"></div>
            </div>
        );
    }
}
