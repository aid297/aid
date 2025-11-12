import moment from "moment"

export class EveryDaysModule {
  _data = []

  constructor() {}

  static new = () => new EveryDaysModule()

  get data() {
    return this._data
  }

  push(data) {
    this._data.push(data)
  }

  find(date = "") {
    let ret = []
    this._data.forEach((items) => {
      const s = items[2]
      const matches = s.match(/\b\d{2}-\d{2}-\d{2}\b/g)
      if (matches.length === 0) return
      if (matches[0] === moment(date).format("YY-MM-DD")) {
        ret = items
        return
      }
    })

    if (ret.length !== 0) {
      return ret[4]?.startsWith("次日")
    }

    return false
  }
}
