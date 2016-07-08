class TodosController {

  constructor({send}) {
    this.send = send
  }

  create(req, res) {
    this.send('todo.create', req.body, () => {
      res.redirect('/')
    })
  }

}
module.exports = TodosController
