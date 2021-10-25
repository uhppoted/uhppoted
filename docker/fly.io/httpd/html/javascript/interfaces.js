import { busy, unbusy, warning, postAsJSON } from './uhppoted.js'
import * as system from './system.js'
import { DB } from './db.js'

export function refreshed () {
  DB.interfaces.forEach(c => {
    updateFromDB(c.OID, c)
  })
}

function updateFromDB (oid, record) {
  const section = document.querySelector(`[data-oid="${oid}"]`)

  if (section) {
    const name = section.querySelector(`[data-oid="${oid}.1"]`)
    const bind = section.querySelector(`[data-oid="${oid}.2"]`)
    const broadcast = section.querySelector(`[data-oid="${oid}.3"]`)
    const listen = section.querySelector(`[data-oid="${oid}.4"]`)

    update(name, record.name)
    update(bind, record.bind)
    update(broadcast, record.broadcast)
    update(listen, record.listen)
  }
}

export function set (element, value, status) {
  const oid = element.dataset.oid
  const original = element.dataset.original
  const v = value.toString()
  const flag = document.getElementById(`F${oid}`)

  element.dataset.value = v

  if (v !== original) {
    mark('modified', element, flag)
  } else {
    unmark('modified', element, flag)
  }

  percolate(oid)
}

export function rollback (tag, element) {
  const section = document.getElementById(tag)
  const oid = section.dataset.oid

  const children = section.querySelectorAll(`[data-oid^="${oid}."]`)
  children.forEach(e => {
    const flag = document.getElementById(`F${e.dataset.oid}`)

    e.dataset.value = e.dataset.original
    e.value = e.dataset.original
    e.classList.remove('modified')

    if (flag) {
      flag.classList.remove('modified')
      flag.classList.remove('pending')
    }
  })

  section.classList.remove('modified')
}

export function commit (element) {
  const section = document.getElementById('interface')
  const oid = section.dataset.oid
  const list = []

  const children = section.querySelectorAll(`[data-oid^="${oid}."]`)
  children.forEach(e => {
    if (e.dataset.value !== e.dataset.original) {
      list.push(e)
    }
  })

  const records = []
  list.forEach(e => {
    const oid = e.dataset.oid
    const value = e.dataset.value
    records.push({ oid: oid, value: value })
  })

  const reset = function () {
    list.forEach(e => {
      const flag = document.getElementById(`F${e.dataset.oid}`)
      unmark('pending', e, flag)
      mark('modified', e, flag)
    })
  }

  const cleanup = function () {
    list.forEach(e => {
      const flag = document.getElementById(`F${e.dataset.oid}`)
      unmark('pending', e, flag)
    })
  }

  list.forEach(e => {
    const flag = document.getElementById(`F${e.dataset.oid}`)
    mark('pending', e, flag)
    unmark('modified', e, flag)
  })

  post('objects', records, reset, cleanup)
}

function modified (oid) {
  const container = document.querySelector(`[data-oid="${oid}"]`)
  let changed = false

  if (container) {
    const list = document.querySelectorAll(`[data-oid^="${oid}."]`)
    list.forEach(e => {
      changed = changed || e.classList.contains('modified')
    })

    if (changed) {
      container.classList.add('modified')
    } else {
      container.classList.remove('modified')
    }
  }
}

function update (element, value, status) {
  if (element) {
    const v = value.toString()
    const oid = element.dataset.oid
    const flag = document.getElementById(`F${oid}`)
    const previous = element.dataset.original

    element.dataset.original = v

    // check for conflicts with concurrently edited fields
    if (element.classList.contains('modified')) {
      if (previous !== v && element.dataset.value !== v) {
        mark('conflict', element, flag)
      } else if (element.dataset.value !== v) {
        unmark('conflict', element, flag)
      } else {
        unmark('conflict', element, flag)
        unmark('modified', element, flag)
      }

      percolate(oid)
      return
    }

    // check for conflicts with concurrently submitted fields
    if (element.classList.contains('pending')) {
      if (previous !== v && element.dataset.value !== v) {
        mark('conflict', element, flag)
      } else {
        unmark('conflict', element, flag)
      }

      return
    }

    // update fields not pending or modified
    if (element !== document.activeElement) {
      element.value = v
    }

    set(element, value)
  }
}

function mark (clazz, ...elements) {
  elements.forEach(e => {
    if (e) {
      e.classList.add(clazz)
    }
  })
}

function unmark (clazz, ...elements) {
  elements.forEach(e => {
    if (e) {
      e.classList.remove(clazz)
    }
  })
}

function percolate (oid) {
  let oidx = oid
  while (oidx) {
    const match = /(.*?)(?:[.][0-9]+)$/.exec(oidx)
    oidx = match ? match[1] : null
    if (oidx) {
      modified(oidx)
    }
  }
}

function post (tag, records, reset, cleanup) {
  busy()

  postAsJSON('/system', { [tag]: records })
    .then(response => {
      if (response.redirected) {
        window.location = response.url
      } else {
        switch (response.status) {
          case 200:
            response.json().then(object => {
              if (object && object.system && object.system.objects) {
                DB.updated('objects', object.system.objects)
              }

              system.refreshed()
            })
            break

          default:
            reset()
            response.text().then(message => { warning(message) })
        }
      }
    })
    .catch(function (err) {
      reset()
      warning(`Error committing record (ERR:${err.message.toLowerCase()})`)
    })
    .finally(() => {
      cleanup()
      unbusy()
    })
}
