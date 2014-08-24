window.App = Ember.Application.create();

//App.Store = DS.Store.extend();

//App.ApplicationAdapter = DS.FixtureAdapter.extend();

App.ApplicationAdapter = DS.RESTAdapter.extend({
        namespace: 'api'
    });
    
App.TodoAdapter = DS.FixtureAdapter.extend();    

