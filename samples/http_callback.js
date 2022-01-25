import http from "k6/http"

globalThis.requestP = (method, url, body, params) => {
  return new Promise((resolve, reject) => {
    params = params || {};
    params.callback = (resp, error) => {
      if (error != null) {
        reject(error)
      } else {
        resolve(resp)
      }
    }
    http.request(method, url, body, params)
  })
}

export default () => {
  http.request("GET", "https://httpbin.test.k6.io/delay/5", null, { callback: (resp) => { console.log(JSON.stringify(resp, null, "  ")) } })
  requestP("GET", "https://httpbin.test.k6.io/delay/4").then((resp) => { console.log(resp) })
  console.log("something")

}
