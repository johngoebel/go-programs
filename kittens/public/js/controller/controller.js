// Add a controller to dispatch create commands
App.CreateController = Ember.Controller.extend( {
    name: null,
    actions: {
        save:function() {
            var kitten =this.store.createRecord('kitten');
            kitten.set('name', this.get('name'));
            kitten.save().then(function() {
                this.transitionToRoute('index');
                this.set('name', '');
            }.bind(this));
        }
    }
});

App.EditController = Ember.ObjectController.extend({
    actions: {
        save: function() {
            var kitten = this.get('model');
            kitten.save().then(function() {
                this.transitionToRoute('index');
            }.bind(this));
        }
    }
});