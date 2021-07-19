let apiHost
if (process.env.NODE_ENV === "development") {
  apiHost = 'http://localhost:8000'
}else {
  apiHost = window.location.origin
}

const apiPrefix = apiHost+'/api'

export {apiHost,apiPrefix}
