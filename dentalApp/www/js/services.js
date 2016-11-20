angular.module('starter.services', [])

.factory('Chats', function() {
  // Might use a resource here that returns a JSON array

  // Some fake testing data
  var chats = [{
    id: 0,
    name: '151662135',
    lastText: 'root canal treatment',
    face: 'img/tooth1.jpg'
  }, {
    id: 1,
    name: '5151515151',
    lastText: 'Lorem ipsum dolor sit amet, eu vivendum torquatos vim, te sed noster dolorem,' +
    ' aliquip noluisse mnesarchum sea id. Has mutat fugit nusquam in. Mel ad posse comprehensam,' +
    ' ut suas accumsan vim, mei id ullum virtute. At nec sint velit aeterno, fugit erant ius ex,' +
    ' te duo indoctum corrumpit. Eu sanctus periculis sed.Sit ne numquam ponderum instructior, ' +
    'cu eum rebum augue audire. Ad usu nonumy consetetur, eos ex sonet molestiae.' +
    ' Vide assum eruditi an quo. Nec erat postea ocurreret te, qui in atqui iisque repudiare.',
    face: 'img/tooth2.jpg'
  }, {
    id: 2,
    name: '5444754577',
    lastText: 'pulling',
    face: 'img/tooth3.jpg'
  }];

  return {
    all: function() {
      return chats;
    },
    remove: function(chat) {
      chats.splice(chats.indexOf(chat), 1);
    },
    get: function(chatId) {
      for (var i = 0; i < chats.length; i++) {
        if (chats[i].id === parseInt(chatId)) {
          return chats[i];
        }
      }
      return null;
    }
  };
});
