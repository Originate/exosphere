class IndexController {

  constructor({send}) {
    this.send = send
  }

  index(req, res) {
    this.send('todo.list', {}, (todos) => {
      res.render('index', {todos})
    })
  }

}

module.exports = IndexController
