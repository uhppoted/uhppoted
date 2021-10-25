/* global messages */

import { postAsForm } from './uhppoted.js'

export function login (event) {
  event.preventDefault()

  const credentials = {
    uid: document.getElementById('uid').value,
    pwd: document.getElementById('pwd').value
  }

  postAsForm('/authenticate', credentials)
    .then(response => {
      switch (response.status) {
        case 200:
          if (response.redirected) {
            return response.url
          } else {
            return '/index.html'
          }

        case 401:
          throw new Error(messages.unauthorized)

        default:
          throw new Error(response.text())
      }
    })
    .then(url => {
      window.location = url
    })
    .catch(function (err) {
      warning(`Error logging in (${err.message.toLowerCase()})`)
    })
}

export function showHidePassword () {
  const pwd = document.getElementById('pwd')
  const eye = document.getElementById('eye')

  if (pwd.type === 'password') {
    pwd.type = 'text'
    eye.src = 'images/eye-slash-solid.svg'
  } else {
    pwd.type = 'password'
    eye.src = 'images/eye-solid.svg'
  }
}

function warning (msg) {
  const message = document.getElementById('message')
  const text = document.getElementById('warning')

  if (text != null) {
    text.innerText = msg
    message.style.visibility = 'visible'
  } else {
    alert(msg)
  }
}
