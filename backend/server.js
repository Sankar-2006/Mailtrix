const express = require('express')

const app = express()

const appRoutes = require('./routes/routes')

app.use(express.urlencoded({extended: true}))

appRoutes(app)

const port = 8080

app.listen(port, () => {
    console.log(`Server started move to localhost:8080/`)
})