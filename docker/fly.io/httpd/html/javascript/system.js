import * as controllers from './controllers.js'
import * as LAN from './interfaces.js'

export function refreshed () {
  LAN.refreshed()
  controllers.refreshed()
}
