import { DateTitleService } from "src/services/dateTitleService"
import { ClockIn, ClockInType } from "src/modules/clockInModule"
import { Str } from "src/utils/str"

/**
 * 服务：打卡数据
 */
export class ClockInService {
  _original = []
  _dateTitles = DateTitleService.make()
  _clockInData = []

  constructor(original = null, dateTitles = null) {
    if (original !== null) this._original = original
    if (dateTitles !== null) this._dateTitles = dateTitles
  }

  static new = (original, dateTitles) => new ClockInService(original, dateTitles)

  static make = () => new ClockInService()

  get original() {
    return this._original
  }
  get dateTitles() {
    return this._dateTitles
  }
  get data() {
    return this._clockInData
  }

  set original(value) {
    this._original = value
  }
  set dateTitles(value) {
    this._dateTitles = value
  }

  parse() {
    this._original.forEach((row) => {
      let tmp = []
      const name = Str.new(row[1]).replace()
      tmp.push(ClockIn.new(row[0], ClockInType.rowNumber))
      tmp.push(ClockIn.new(name, ClockInType.name))
      this._dateTitles.data.forEach((item) => tmp.push(ClockIn.new(row[item.idx], ClockInType.clockIn, this._dateTitles, false)))
      this._clockInData.push(tmp)
    })

    return this
  }
}
