class HomeController {

  constructor({send}) {
    this.send = send
  }

  index(req, res) {
    res.render('index', {})
  };
}



module.exports = HomeController
