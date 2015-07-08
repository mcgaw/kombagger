import sinon from 'sinon';

import {
  moduleFor,
  test
} from 'ember-qunit';

let fakeServer;

moduleFor('route:leaderboard', {
  beforeEach: function() {
  		fakeServer = sinon.fakeServer.create();
   },

  afterEach: function() {
  		fakeServer.restore();
  }
});

var sampleLeaderboardResponse = `[
    {
        "id": 1002872,
        "name": "Andy Schofield",
        "koms": [
            {
                "effort_id": 7103008581,
                "athlete_name": "Andy Schofield",
                "athlete_id": 1002872,
                "elapsed_time": 34,
                "moving_time": 34,
                "rank": 1,
                "athlete_profile": "https://dgalywyr863hv.cloudfront.net/pictures/athletes/1002872/2572875/2/large.jpg"
            }
        ]
    },
    {
        "id": 1149854,
        "name": "Gary Macdonald",
        "koms": [
            {
                "effort_id": 6232958474,
                "athlete_name": "Gary Macdonald",
                "athlete_id": 1149854,
                "elapsed_time": 246,
                "moving_time": 246,
                "rank": 1,
                "athlete_profile": "https://dgalywyr863hv.cloudfront.net/pictures/athletes/1149854/421167/3/large.jpg"
            }
        ]
    }]`;

test('leaderboard controller is populated with rider array', function(assert) {

	var expectedRiders = [{athleteName : 'rider1'}, {athleteName : 'rider2'}];
	var route = this.subject();

	fakeServer.respondWith("GET", /.*\/api\/leaderboard.*/,
            [200, { "Content-Type": "application/json" },
            sampleLeaderboardResponse]);


	// fakeServer.respondWith("GET", /.*\/api\/leaderboard.*/,
 //            [200, { "Content-Type": "application/json" },
 //            '[{ "athlete_name" : "rider1"}, { "athlete_name" : "rider2"}]']);


	var controller = sinon.spy();

	route.setupController(controller, null);

	fakeServer.respond();

	assert.equal(controller.leaderboard.length, 2, "incorrect number of riders");
	assert.deepEqual(controller.leaderboard[1], {id : 1149854, name : "Gary Macdonald",
		koms : [{effort_id : 6232958474, athlete_name : "Gary Macdonald", athlete_id : 1149854, elapsed_time : 246, moving_time : 246, rank : 1,
		athlete_profile : "https://dgalywyr863hv.cloudfront.net/pictures/athletes/1149854/421167/3/large.jpg"}] }, "rider data not populated correctly");
	

});