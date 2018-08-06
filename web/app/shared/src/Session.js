import React from "react";                                                                                                          1
import axios from "axios";

export const validateSession = (obj)=>{
	sessionKey = window.sessionKey;
	var result = '{"loggedin": false, "email":""}';
	const client = axios.create({
		headers: {
			'X-Requested-With': 'XMLHttpRequest',
			'X-CSRF-TOKEN' : String(sessionKey),
		},
	});
	client.post('/api/session/validate', {
		timeout: 2000
	})
	.then((response) => {
		switch(response.status) {
			case 200:
				switch(response.data["error"]) {
					case "":
						obj.setState({loggedin: true, email: response.data["user"]["email"], sessionKey: sessionKey});
						break;
					default:
						// console.log(response.data["error"]);
				}
				break;
			default:
				console.log(response);
		}
	})
	.catch(function (error) {
	});
}
