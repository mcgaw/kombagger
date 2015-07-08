import Ember from 'ember';

export default Ember.Controller.extend({
	leaderboard : [],
	fixAvatarPics : function() {
		this.leaderboard.forEach(function(rider) {
			var url = rider.koms[0].athlete_profile;
			console.log(url);
			if (url.match('avatar')) {
				rider.athlete_profile = '/assets/images/avatar_large.png';
			} else {
				rider.athlete_profile = url;
			}
		})
	}
});

