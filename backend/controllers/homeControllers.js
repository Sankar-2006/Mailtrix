const path = require('path')

let getHome = (req, res) => {
    return res.sendFile(path.join(`${__dirname}/../$homePath`)) // TODO: $homepath need to updated by the real html file
}

module.exports = {
    getHome: getHome
}