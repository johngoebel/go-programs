// App.Router.map(function() {
//   this.route('create');
//   this.route('edit', {path: '/edit/:kitten_id'});
// });

 App.Router.map(function() {
    this.resource('create');
    this.route('edit', {path: '/edit/:kitten_id'});
});

App.IndexRoute = Ember.Route.extend({
      model: function() {
          return this.store.find('kitten');
      },

      actions: {
          deleteKitten: function(kitten) {
              kitten.deleteRecord();
              kitten.save();
          }
      }
});
