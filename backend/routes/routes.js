const express = require('express')

// using express router
const router = express.Router()

const emailController = require('../controllers/emailControllers')
const homeController = require('../controllers/homeControllers')

let appRoutes = (app) => {
    router.get('/', homeController.getHome)
    router.post('/send-email', emailController.sendMail)

    return app.use('/', router)
};

module.exports = appRoutes