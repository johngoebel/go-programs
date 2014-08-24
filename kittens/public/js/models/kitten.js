App.Kitten = DS.Model.extend({
  name: DS.attr('string'),
  picture: DS.attr('string')
});


// ... additional lines truncated for brevity ...
App.Kitten.FIXTURES = [
 {
   id: 1,
   title: 'Learn Ember - wow.js',
   isCompleted: true
 },
 {
   id: 2,
   title: '...you dont say',
   isCompleted: false
 },
 {
   id: 3,
   title: 'Profit!',
   isCompleted: false
 }
];