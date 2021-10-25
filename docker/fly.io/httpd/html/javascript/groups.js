import { update, deleted } from './tabular.js'
import { DB } from './db.js'

export function refreshed () {
  const groups = [...DB.groups.values()].sort((p, q) => p.index - q.index)

  realize(groups)

  groups.forEach(o => {
    const row = updateFromDB(o.OID, o)
    if (row) {
      if (o.status === 'new') {
        row.classList.add('new')
      } else {
        row.classList.remove('new')
      }
    }
  })

  DB.refreshed('groups')
}

function updateFromDB (oid, record) {
  const row = document.querySelector("div#groups tr[data-oid='" + oid + "']")

  if (record.status === 'deleted' || !row) {
    return
  }

  const name = row.querySelector(`[data-oid="${oid}.1"]`)
  const doors = [...DB.doors.values()].filter(o => o.status && o.status !== '<new>' && o.status !== 'deleted')

  row.dataset.status = record.status

  update(name, record.name)

  doors.forEach(o => {
    const td = row.querySelector(`td[data-door="${o.OID}"]`)

    if (td) {
      const e = td.querySelector('.field')
      const d = record.doors.get(`${e.dataset.oid}`)

      update(e, d && d.allowed)
    }
  })

  return row
}

function realize (groups) {
  const table = document.querySelector('#groups table')
  const thead = table.tHead
  const tbody = table.tBodies[0]

  const doors = new Map([...DB.doors.values()]
    .filter(o => o.status && o.status !== '<new>' && o.status !== 'deleted')
    .sort((p, q) => p.index - q.index)
    .map(o => [o.OID, o]))

  // ... columns

  const columns = table.querySelectorAll('th.door')
  const cols = new Map([...columns].map(c => [c.dataset.door, c]))
  const missing = [...doors.values()].filter(o => o.OID === '' || !cols.has(o.OID))
  const surplus = [...cols].filter(([k]) => !doors.has(k))

  missing.forEach(o => {
    const th = thead.rows[0].lastElementChild
    const padding = thead.rows[0].appendChild(document.createElement('th'))

    padding.classList.add('colheader')
    padding.classList.add('padding')

    th.classList.replace('padding', 'door')
    th.dataset.door = o.OID
    th.innerHTML = o.name
  })

  surplus.forEach(([, v]) => {
    v.remove()
  })

  // ... rows

  groups.forEach(o => {
    let row = tbody.querySelector("tr[data-oid='" + o.OID + "']")

    if (o.status === 'deleted') {
      deleted('groups', row)
      return
    }

    if (!row) {
      row = add(o.OID, o)
    }

    const columns = row.querySelectorAll('td.door')
    const cols = new Map([...columns].map(c => [c.dataset.door, c]))
    const missing = [...doors.values()].filter(o => o.OID === '' || !cols.has(o.OID))
    const surplus = [...cols].filter(([k]) => !doors.has(k))

    missing.forEach(o => {
      const door = o.OID.match(/^0\.2\.([1-9][0-9]*)$/)[1]
      const template = document.querySelector('#door')

      const uuid = row.id
      const oid = row.dataset.oid + '.2.' + door
      const ix = row.cells.length - 1
      const cell = row.insertCell(ix)

      cell.classList.add('door')
      cell.dataset.door = o.OID
      cell.innerHTML = template.innerHTML

      const flag = cell.querySelector('.flag')
      const field = cell.querySelector('.field')

      flag.classList.add(`d${door}`)
      field.classList.add(`d${door}`)

      flag.id = 'F' + oid

      field.id = uuid + '-' + `d${door}`
      field.dataset.oid = oid
      field.dataset.record = uuid
      field.dataset.original = ''
      field.dataset.value = ''
      field.checked = false
    })

    surplus.forEach(([, v]) => {
      v.remove()
    })
  })
}

function add (oid, record) {
  const uuid = 'R' + oid.replaceAll(/[^0-9]/g, '')
  const tbody = document.getElementById('groups').querySelector('table tbody')

  if (tbody) {
    const template = document.querySelector('#group')
    const row = tbody.insertRow()

    row.id = uuid
    row.classList.add('group')
    row.classList.add('new')
    row.dataset.oid = oid
    row.dataset.status = 'unknown'
    row.innerHTML = template.innerHTML

    const commit = row.querySelector('td span.commit')
    commit.id = uuid + '_commit'
    commit.dataset.record = uuid
    commit.dataset.enabled = 'false'

    const rollback = row.querySelector('td span.rollback')
    rollback.id = uuid + '_rollback'
    rollback.dataset.record = uuid
    rollback.dataset.enabled = 'false'

    const fields = [
      { suffix: 'name', oid: `${oid}.1`, selector: 'td input.name', flag: 'td img.name' }
    ]

    fields.forEach(f => {
      const field = row.querySelector(f.selector)
      const flag = row.querySelector(f.flag)

      if (field) {
        field.id = uuid + '-' + f.suffix
        field.value = ''
        field.dataset.oid = f.oid
        field.dataset.record = uuid
        field.dataset.original = ''
        field.dataset.value = ''

        flag.id = 'F' + f.oid
      } else {
        console.error(f)
      }
    })

    return row
  }
}
