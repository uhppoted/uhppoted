import { update, deleted } from './tabular.js'
import { DB } from './db.js'

export function refreshed () {
  const cards = [...DB.cards.values()].sort((p, q) => p.created.localeCompare(q.created))
  const pagesize = 1

  realize(cards)

  // renders a 'page size' chunk of cards
  const f = function (offset) {
    let ix = offset
    let count = 0
    while (count < pagesize && ix < cards.length) {
      const o = cards[ix]
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

  // sorts the table rows by 'created'
  const g = function () {
    const table = document.querySelector('#cards table')
    const tbody = table.tBodies[0]

    tbody.sort((p, q) => {
      const u = DB.cards.get(p.dataset.oid)
      const v = DB.cards.get(q.dataset.oid)

      return u.created.localeCompare(v.created)
    })
  }

  const chunk = offset => new Promise(resolve => {
    f(offset)
    resolve(true)
  })

  async function * render () {
    for (let ix = 0; ix < cards.length; ix += pagesize) {
      yield chunk(ix).then(() => ix)
    }
  }

  (async function loop () {
    for await (const _ of render()) {
      // empty
    }
  })()
    .then(() => g())
    .then(() => DB.refreshed('cards'))
    .catch(err => console.error(err))
}

function updateFromDB (oid, record) {
  const row = document.querySelector("div#cards tr[data-oid='" + oid + "']")

  if (record.status === 'deleted' || !row) {
    return
  }

  const name = row.querySelector(`[data-oid="${oid}.1"]`)
  const number = row.querySelector(`[data-oid="${oid}.2"]`)
  const from = row.querySelector(`[data-oid="${oid}.3"]`)
  const to = row.querySelector(`[data-oid="${oid}.4"]`)
  const groups = [...DB.groups.values()].filter(g => g.status && g.status !== '<new>' && g.status !== 'deleted')

  row.dataset.status = record.status

  update(name, record.name)
  update(number, record.number)
  update(from, record.from)
  update(to, record.to)

  groups.forEach(g => {
    const td = row.querySelector(`td[data-group="${g.OID}"]`)

    if (td) {
      const e = td.querySelector('.field')
      const g = record.groups.get(`${e.dataset.oid}`)

      update(e, g && g.member)
    }
  })

  return row
}

function realize (cards) {
  const table = document.querySelector('#cards table')
  const thead = table.tHead
  const tbody = table.tBodies[0]

  const groups = new Map([...DB.groups.values()]
    .filter(o => o.status && o.status !== '<new>' && o.status !== 'deleted')
    .sort((p, q) => p.created.localeCompare(q.created))
    .map(o => [o.OID, o]))

  // ... columns

  const columns = table.querySelectorAll('th.group')
  const cols = new Map([...columns].map(c => [c.dataset.group, c]))
  const missing = [...groups.values()].filter(o => o.OID === '' || !cols.has(o.OID))
  const surplus = [...cols].filter(([k]) => !groups.has(k))

  missing.forEach(o => {
    const th = thead.rows[0].lastElementChild
    const padding = thead.rows[0].appendChild(document.createElement('th'))

    padding.classList.add('colheader')
    padding.classList.add('padding')

    th.classList.replace('padding', 'group')
    th.dataset.group = o.OID
    th.innerHTML = o.name
  })

  surplus.forEach(([, v]) => {
    v.remove()
  })

  // ... rows

  cards.forEach(o => {
    let row = tbody.querySelector("tr[data-oid='" + o.OID + "']")

    if (o.status === 'deleted') {
      deleted('cards', row)
      return
    }

    if (!row) {
      row = add(o.OID, o)
    }

    const columns = row.querySelectorAll('td.group')
    const cols = new Map([...columns].map(c => [c.dataset.group, c]))
    const missing = [...groups.values()].filter(o => o.OID === '' || !cols.has(o.OID))
    const surplus = [...cols].filter(([k]) => !groups.has(k))

    missing.forEach(o => {
      const group = o.OID.match(/^0\.4\.([1-9][0-9]*)$/)[1]
      const template = document.querySelector('#group')

      const uuid = row.id
      const oid = row.dataset.oid + '.5.' + group
      const ix = row.cells.length - 1
      const cell = row.insertCell(ix)

      cell.classList.add('group')
      cell.dataset.group = o.OID
      cell.innerHTML = template.innerHTML

      const flag = cell.querySelector('.flag')
      const field = cell.querySelector('.field')

      flag.classList.add(`g${group}`)
      field.classList.add(`g${group}`)

      flag.id = 'F' + oid

      field.id = uuid + '-' + `g${group}`
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
  const tbody = document.getElementById('cards').querySelector('table tbody')

  if (tbody) {
    const template = document.querySelector('#card')
    const row = tbody.insertRow()

    row.id = uuid
    row.classList.add('card')
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
      { suffix: 'name', oid: `${oid}.1`, selector: 'td input.name', flag: 'td img.name' },
      { suffix: 'number', oid: `${oid}.2`, selector: 'td input.number', flag: 'td img.number' },
      { suffix: 'from', oid: `${oid}.3`, selector: 'td input.from', flag: 'td img.from' },
      { suffix: 'to', oid: `${oid}.4`, selector: 'td input.to', flag: 'td img.to' }
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
