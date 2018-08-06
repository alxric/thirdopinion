import React from "react";                                                                                                          1
import axios from "axios";

function VoteButtons(props) {
	return (
		<div className="voteWrapper">
			<div className="voteHeader">
				Who do you agree with?
			</div>
			<div className="voteButtonDiv leftFloat">
				<button className="btn btn-default person1" type="submit" id={'voteButton_' + props.Index + '_' + props.ArgumentID + '_1'} onClick={props.Obj.handleClick}>
					Person 1
				</button>
			</div>
			<div className="voteButtonDiv leftFloat">
				<button className="btn btn-default person1" type="submit" id={'voteButton_' + props.Index + '_' + props.ArgumentID + '_2'} onClick={props.Obj.handleClick}>
					Person 2
				</button>
			</div>
			<div className="clear"></div>
		</div>
	);
}

function VoteDisplay(props) {
	return (
		<div className="voteWrapper">
			<div className="voteHeader">
				Votes so far
			</div>
			<div className="voteDisplay">
				<span className="person1 voteSpan">
					{props.P1Votes + ' (' + props.P1Percent + '%)'}
				</span>
				<span className="person2 voteSpan">
					{props.P2Votes + ' (' + props.P2Percent + '%)'}
				</span>
			</div>
		</div>
	);
}

export default class View extends React.Component {
	constructor(props) {
		super(props);
		this.state = {argument_id: '', args: [], error: '', votes: null};
		this.state.votes = JSON.parse(window.votes);
		this.handleClick = this.handleClick.bind(this);
		this.generateArguments();
	}

	generateFilter() {
		var argument_id = '';
		var baseURI = window.location.href;
		var s = baseURI.split("/view/");
		if(s.length == 2) {
			argument_id = s[1]
		}
		this.state.argument_id = argument_id
	}

	generateArguments() {
		this.generateFilter();
		var endpoint = '/api/arguments';
		if(this.state.argument_id != '') {
			endpoint += '?id=' + this.state.argument_id;
		}
		axios.get(endpoint, {
			argument: {
				id: parseInt(this.state.argument_id),
			},
		})
		.then((response) => {
			switch(response.data.error) {
				case "":
					this.setState({args: response.data["arguments"]});
					break;
				default:
					this.setState({error: response.data["error"]});
					break;
			}
		})
		.catch(function (error) {
			console.log(error);
		});
	}

	handleClick(event) {
		if(!this.props.Loggedin) {
			window.location.href='/login';
		}
		var idVals = event.target.id.split("_");
		if(idVals.length != 4) {
			console.log("Invalid button ID");
			return
		}
		var index = idVals[1];
		var argID = parseInt(idVals[2]);
		var person = parseInt(idVals[3]);
		var targs = this.state.args;

		const client = axios.create({
			headers: {
				'X-Requested-With': 'XMLHttpRequest',
				'X-CSRF-TOKEN' : this.props.SessionKey,
			},
		});

		client.post('/api/vote', {
			vote: {
				argument_id: argID,
				person: person
			}
		})
		.then((response) => {
			switch(response.data["error"]) {
				case "":
					switch(person) {
						case 1:
							targs[index]['votes'].person_1++;
							break;
						case 2:
							targs[index]['votes'].person_2++;
							break;
					}
					this.setState({args: targs});
					votes = this.state.votes;
					votes[argID] = person;
					this.setState({votes: votes});
					break;
				default:
					console.log(response.data["error"]);
					break;
			}
		})
		.catch(function (error) {
			console.log(error);
		});
	}

    render() {
		if(this.state.args.length == 0) {
			return (
				<div>No arguments found</div>
			);
		}
		return (
			<div id="arguments">
                {this.state.args.map((argument, index) => {
					var totalVotes = argument['votes'].person_1 + argument['votes'].person_2;
					if(totalVotes == 0) {
						var p1percent = 0;
						var p2percent = 0;
					} else {
						var p1percent =  Math.ceil((argument['votes'].person_1 / totalVotes) * 100);
						var p2percent = 100-p1percent;
					}

					const hasVoted = this.props.Loggedin && argument['id'] in this.state.votes;
					let voteArea;
					if(hasVoted) {
						voteArea = <VoteDisplay P1Votes={argument['votes'].person_1} P2Votes={argument['votes'].person_2} P1Percent={p1percent} P2Percent={p2percent} />
					} else {
						voteArea = <VoteButtons Index={index} ArgumentID={argument['id']} Obj={this} />
					}
                    return (
						<div className="argumentWrapper">
							<div className="argumentHeader">
								<a className="argumentLink" href={'/view/' + argument['id']}>
									<h3>{argument['title']}</h3>
								</a>
							</div>
							<div className="opinionContainer">
								{argument['opinions'].map((opinion, oindex) => {
									return (
										<div className="opinionWrapper">
											<div className={'person' + opinion['person']}>
												{'<Person ' + opinion['person'] + '> ' + opinion['text'] }
											</div>
										</div>
									);
								})}
							</div>
							{voteArea}
						</div>
					);
                  })}
			</div>
		);
    }
}
