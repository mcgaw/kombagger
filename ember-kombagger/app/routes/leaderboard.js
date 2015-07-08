import Ember from 'ember';

export default Ember.Route.extend({
	setupController: function setupController(controller, model) {

		/* jshint unused:false */
		var req = Ember.$.ajax('http://localhost:9090/leaderboard');

		req.done(function (msg) {
			Ember.set(controller, 'leaderboard', msg);
			controller.fixAvatarPics();
		});

		req.fail(function (resp, status, err) {
			controller.transitionToRoute('/error');
		});
	}

});