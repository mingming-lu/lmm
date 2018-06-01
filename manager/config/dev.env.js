'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  API_URL_BASE: '"http://api.lmm.local"',
  IMAGE_URL_BASE: '"http://image.lmm.local"'
})
