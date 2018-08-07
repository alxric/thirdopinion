import React from "react";                                                                                                          1
import axios from "axios";

export default class View extends React.Component {
	constructor(props) {
		super(props);
		this.state = {error: '', title: '', titleClass: 'form-control', initialClass: 'form-control', currentOpinion: '', opinions: [], nextPerson: 1};
		this.handleChange = this.handleChange.bind(this);
		this.selectChange = this.selectChange.bind(this);
		this.handleClick = this.handleClick.bind(this);
		this.addOpinion = this.addOpinion.bind(this);
		this.deleteOpinion = this.deleteOpinion.bind(this);
	}

	selectChange(event) {
		var index = event.target.id.split("peopleSelect")[1];
		var opinions = this.state.opinions;
		var opinion = this.state.opinions[parseInt(index)];
		opinion["person"] = 3 - opinion["person"];
		opinions[index] = opinion;
		this.setState({opinions: opinions});
	}

	handleChange(event) {
		this.setState({error: ''});
		switch(event.target.id) {
			case "argumentTitleInput":
				this.setState({title: event.target.value, titleClass: 'form-control'});
				break;
			case "initialInput":
				this.setState({currentOpinion: event.target.value, initialClass: 'form-control'});
				break;
		}
	}

	handleClick(event) {
		event.preventDefault();
		const client = axios.create({
			headers: {
				'X-Requested-With': 'XMLHttpRequest',
				'X-CSRF-TOKEN' : this.props.SessionKey,
			},
		});

		client.post('/api/create', {
			argument: {
				title: this.state.title,
				opinions: this.state.opinions,
			}
		})
		.then((response) => {
			switch(response.data["error"]) {
				case "Title too short":
					this.setState({error: response.data["error"], titleClass: 'form-control errorInput'});
					break;
				case "You need at least two opinions for an argument":
					this.setState({error: response.data["error"], initialClass: 'form-control errorInput'});
					break;
				case "":
					window.location.href = '/view/'+response.data["argument_id"];
			}
			console.log(response);
		})
		.catch(function (error) {
			console.log(error);
		});
	}

	deleteOpinion(event) {
		var index = event.target.id.split("delete")[1];
		var opinions = this.state.opinions;
		opinions.splice(index, 1);
		this.setState({opinions: opinions});
	}

	addOpinion(event) {
		event.preventDefault();
		var opinions = this.state.opinions
		var opinion = '{"person":' + this.state.nextPerson + ',"text":"'+this.state.currentOpinion+'"}';
		if(this.state.currentOpinion == "") {
			return
		}
		var nextPerson = 3 - this.state.nextPerson;
		opinions.push(JSON.parse(opinion));
		this.setState({opinions: opinions, currentOpinion: '', nextPerson: nextPerson});
	}

	render() {
		const isLoggedIn = this.props.Loggedin;
		if (!isLoggedIn) {
			return (
				<div className="container-fluid">
					<h3>You need to <a href="/login">log in</a> to create an argument</h3>
				</div>
			);
		}
		var ph = "Enter the opinion of the first person";
		const formCorrect = (this.state.opinions.length >=2 && this.state.title.length >= 3);
		return (
			<div className="container-fluid">
				<div id="newArgumentDiv">
					<form className="form-inline">
						<div id="argumentTitleText">
							<h3>
								Title
							</h3>
						</div>
						<div id="argumentTitleDiv" className="form-group">
							<label className="sr-only" for="argumentTitleInput">
								Title input
							</label>
							<input type="text" className={this.state.titleClass} id="argumentTitleInput" placeholder="Enter argument title" onChange={this.handleChange} />
						</div>
						<div>
							<button type="submit" className="btn btn-default invisible" onClick={this.addOpinion}>
								Add
							</button>
						</div>
						<div id="argumentInputs">
							{this.state.opinions.map((opinion, index) => {
								return (
									<div>
										<div className="argumentInput">
											<div className="leftFloat peopleDiv">
												<label className="sr-only" for={'peopleSelect'+index}>
													Initial input
												</label>
												<select className="peopleSelect form-control" value={'Person '+opinion["person"]} id={'peopleSelect'+index} onChange={this.selectChange}>
													<option>
														Person 1
													</option>
													<option>
														Person 2
													</option>
												</select>
											</div>
											<div className="leftFloat argumentDiv">
												{opinion["text"]}
											</div>
											<div className="rightFloat deleteArgumentDiv">
												<span className="deleteSpan" id={'delete'+index} onClick={this.deleteOpinion}>
													DELETE
												</span>
											</div>
											<div className="clear"></div>
										</div>
										<div className="clear"></div>
									</div>
								);
						  	})}
						</div>
						<div id="inputDiv1" className="form-group">
							<label className="sr-only" for="initialInput">
								Initial input
							</label>
							<input type="text" id="initialInput" value={this.state.currentOpinion} className={this.state.initialClass} placeholder={ph} onChange={this.handleChange} />
						</div>
						<div>
							<div className="leftFloat">
								<span className="error" id="errorSpan">
									{this.state.error}
								</span>
							</div>
							<div className="rightFloat">
								{formCorrect ? (
									<button type="submit" className="btn btn-default" onClick={this.handleClick}>
										Create
									</button>
								) : (
									<button type="submit" className="btn btn-default" disabled>
										Create
									</button>
								)}
							</div>
							<div className="clear"></div>
						</div>
					</form>
				</div>
			</div>
		);
	}
}
