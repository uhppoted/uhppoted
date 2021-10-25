import { deleted } from './tabular.js'
import { DB } from './db.js'

export function refreshed () {
  const entries = [...DB.logs().values()].sort((p, q) => q.timestamp.localeCompare(p.timestamp))
  const pagesize = 5

  realize(entries)

  // renders a 'page size' of log entries
  const f = function (offset) {
    let ix = offset
    let count = 0
    while (count < pagesize && ix < entries.length) {
      const o = entries[ix]
      const row = updateFromDB(o.OID, o)
      if (row) {
        if (o.status === 'new') {
          row.classList.add('new')
        } else {
          row.classList.remove('new')
        }
      }

      count++
      ix++
    }
  }

  // sorts the table rows by 'timestamp'
  const g = function () {
    const table = document.querySelector('#logs table')
    const tbody = table.tBodies[0]

    tbody.sort((p, q) => {
      const u = DB.logs().get(p.dataset.oid)
      const v = DB.logs().get(q.dataset.oid)

      return v.timestamp.localeCompare(u.timestamp)
    })
  }

  // hides/shows the 'more' button
  const h = function () {
    const table = document.querySelector('#logs table')
    const tfoot = table.tFoot
    const last = DB.lastLog()

    if (last && DB.logs().has(last)) {
      tfoot.classList.add('hidden')
    } else {
      tfoot.classList.remove('hidden')
    }
  }

  // initialises the rows asynchronously in small'ish chunks
  const chunk = offset => new Promise(resolve => {
    f(offset)
    resolve(true)
  })

  async function * render () {
    for (let ix = 0; ix < entries.length; ix += pagesize) {
      yield chunk(ix).then(() => ix)
    }
  }

  (async function loop () {
    for await (const _ of render()) {
      // empty
    }
  })()
    .then(() => g())
    .then(() => h())
    .then(() => DB.refreshed('logs'))
    .catch(err => console.error(err))
}

function realize (entries) {
  const table = document.querySelector('#logs table')
  const tbody = table.tBodies[0]

  entries.forEach(o => {
    let row = tbody.querySelector("tr[data-oid='" + o.OID + "']")

    if (o.status === 'deleted') {
      deleted('logs', row)
      return
    }

    if (!row) {
      row = add(o.OID, o)
    }
  })
}

function add (oid) {
  const uuid = 'R' + oid.replaceAll(/[^0-9]/g, '')
  const tbody = document.getElementById('logs').querySelector('table tbody')

  if (tbody) {
    const template = document.querySelector('#entry')
    const row = tbody.insertRow()

    row.id = uuid
    row.classList.add('entry')
    row.dataset.oid = oid
    row.dataset.status = 'unknown'
    row.innerHTML = template.innerHTML

    const commit = row.querySelector('td span.commit')
    if (commit) {
      commit.id = uuid + '_commit'
      commit.dataset.record = uuid
      commit.dataset.enabled = 'false'
    }

    const rollback = row.querySelector('td span.rollback')
    if (rollback) {
      rollback.id = uuid + '_rollback'
      rollback.dataset.record = uuid
      rollback.dataset.enabled = 'false'
    }

    const fields = [
      { suffix: 'timestamp', oid: `${oid}.1`, selector: 'td input.timestamp' },
      { suffix: 'uid', oid: `${oid}.2`, selector: 'td input.uid' },
      { suffix: 'module', oid: `${oid}.3`, selector: 'td input.module' },
      { suffix: 'module-id', oid: `${oid}.4`, selector: 'td input.module-id' },
      { suffix: 'module-name', oid: `${oid}.5`, selector: 'td input.module-name' },
      { suffix: 'module-field', oid: `${oid}.6`, selector: 'td input.module-field' },
      { suffix: 'details', oid: `${oid}.7`, selector: 'td input.details' }
    ]

    fields.forEach(f => {
      const field = row.querySelector(f.selector)
      const flag = row.querySelector(`td img.${f.suffix}`)

      if (field) {
        field.id = uuid + '-' + f.suffix
        field.value = ''
        field.dataset.oid = f.oid
        field.dataset.record = uuid
        field.dataset.original = ''
        field.dataset.value = ''

        if (flag) {
          flag.id = 'F' + f.oid
        }
      } else {
        console.error(f)
      }
    })

    return row
  }
}

function updateFromDB (oid, record) {
  const row = document.querySelector("div#logs tr[data-oid='" + oid + "']")

  if (record.status === 'deleted' || !row) {
    return
  }

  const timestamp = row.querySelector(`[data-oid="${oid}.1"]`)
  const uid = row.querySelector(`[data-oid="${oid}.2"]`)
  const module = row.querySelector(`[data-oid="${oid}.3"]`)
  const moduleID = row.querySelector(`[data-oid="${oid}.4"]`)
  const moduleName = row.querySelector(`[data-oid="${oid}.5"]`)
  const moduleField = row.querySelector(`[data-oid="${oid}.6"]`)
  const details = row.querySelector(`[data-oid="${oid}.7"]`)

  row.dataset.status = record.status

  update(timestamp, format(record.timestamp))
  update(uid, record.uid)
  update(module, record.module.type)
  update(moduleID, record.module.ID)
  update(moduleName, record.module.name.toLowerCase())
  update(moduleField, record.module.field.toLowerCase())
  update(details, record.module.details)

  return row
}

function update (element, value) {
  if (element && value !== undefined) {
    element.value = value.toString()
  }
}

function format (timestamp) {
  const dt = Date.parse(timestamp)
  const fmt = function (v) {
    return v < 10 ? '0' + v.toString() : v.toString()
  }

  if (!isNaN(dt)) {
    const date = new Date(dt)
    const year = date.getFullYear()
    const month = fmt(date.getMonth() + 1)
    const day = fmt(date.getDate())
    const hour = fmt(date.getHours())
    const minute = fmt(date.getMinutes())
    const second = fmt(date.getSeconds())

    return `${year}-${month}-${day} ${hour}:${minute}:${second}`
  }

  return ''
}
