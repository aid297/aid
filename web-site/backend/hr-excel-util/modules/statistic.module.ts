import moment from "moment"
import { EventDefault, EventErrAbsenteeism, EventErrLate, EventErrLeaveEarly, EventErrOvertimeCardShortageClosed, EventErrOvertimeCardShortageWorking, EventErrWorkdayCardShortageClosed, EventErrWorkdayCardShortageWorking, EventNotMatch, EventOvertimeHoliday, EventOvertimeWeekend } from "./event.module"

export class StatisticModule {
	private _originalTime?: moment.Moment = undefined
	private _originalTimeText?: string = ""
	private _isWeekend: boolean = false
	private _isHoliday: boolean = false
	private _isWorkday: boolean = true
	private _isAnnualVacation: boolean = false
	private _salaryMultiple: number = 1
	private _idx: number = -1
	private _originalMonthStatistic: Map<string, string> = new Map()
	private _originalClockIn: Map<string, string[]> = new Map()
	private _names: string[] = []
	private _analysisResult: Map<string, EventDefault[] | undefined> = new Map()

	constructor(time: string, text: string = "", idx: number = -1, workdays: string[] = [], annualVacations: string[] = []) {
		this._originalTime = moment(time)
		this._originalTimeText = text
		this._isWeekend = [0, 6].includes(this._originalTime.day())
		this._idx = idx

		this._isWorkday = workdays.includes(time)
		this._isAnnualVacation = annualVacations.includes(time)

		if (this._isWorkday) {
			this._salaryMultiple = 1
		} else if (this._isHoliday) {
			this._salaryMultiple = 3
		} else if (this._isWeekend) {
			this._salaryMultiple = 2
		} else if (this._isAnnualVacation) {
			this._salaryMultiple = 0
		}
	}

	get originalTime() { return this._originalTime }
	get text() { return this._originalTimeText }
	get isWeekend() { return this._isWeekend }
	get isHoliday() { return this._isHoliday }
	get isWorkday() { return this._isWorkday }
	get isAnnualVacation() { return this._isAnnualVacation }
	get salaryMultiple() { return this._salaryMultiple }
	get idx() { return this._idx }
	get originalMonthStatistic() { let m: any = {}; this._originalMonthStatistic.forEach((v, k) => { m[k] = v }); return m }
	get originalClockIn() { let m: any = {}; this._originalClockIn.forEach((v, k) => { m[k] = v }); return m }
	get names() { return this._names }
	get analysisResult() { let m: any = {}; this._analysisResult.forEach((v, k) => { m[k] = v }); return m }

	/**
	 * 转时间到字符串
	 * @returns string
	 */
	toString = (): string => this._originalTime ? this._originalTime.format("YYYY-MM-DD") : "无效时间"

	/**
	 * 转json字符串
	 * @returns string
	 */
	toJson = (): string => JSON.stringify({
		originalTime: this.originalTime,
		originalTimeString: this.toString(),
		text: this.text,
		isWeekend: this.isWeekend,
		isHoliday: this.isHoliday,
		isWorkday: this.isWorkday,
		isAnnualVacation: this.isAnnualVacation,
		salaryMultiple: this.salaryMultiple,
		idx: this._idx,
		originalMonthStatistic: this.originalMonthStatistic,
		originalClockIn: this.originalClockIn,
		analysisResult: this.analysisResult,
	})

	/**
	 * 转字典
	 * @returns {}
	 */
	toDict = (): {} => ({
		originalTime: this.originalTime,
		originalTimeString: this.toString(),
		text: this.text,
		isWeekend: this.isWeekend,
		isHoliday: this.isHoliday,
		isWorkday: this.isWorkday,
		isAnnualVacation: this.isAnnualVacation,
		salaryMultiple: this.salaryMultiple,
		idx: this._idx,
		originalMonthStatistic: this.originalMonthStatistic,
		originalClockIn: this.originalClockIn,
		analysisResult: this.analysisResult,
	})

	/**
	 * 设置原始月度汇总数据
	 * @param name 姓名
	 * @param original 原始数据
	 */
	setOriginalMonthStatistic = (name: string, original: string) => { this._originalMonthStatistic.set(name, original); this._names.push(name) }

	/**
	 * 设置每日打卡数据
	 * @param name 姓名
	 * @param originals 第一次和最后一次打卡时间
	 */
	setOriginalClockIn = (name: string, originals: string[]) => { this._originalClockIn.set(name, originals); this._names.push(name) }

	/**
	 * 分析数据
	 * @returns Map<string, EventDefault>
	 */
	analysis = () => {
		this._names = Array.from(new Set(this._names))

		this._names.forEach(name => {
			if (!this._analysisResult.has(name)) this._analysisResult.set(name, [])

			const monthStatistic = this._originalMonthStatistic.get(name) ?? ""
			const clockIn = this._originalClockIn.get(name) ?? []
			const rightClockInItems = ["补卡审批通过", "正常", "出差", "外勤", "请假"]
			const wrongClockInItems = ["缺卡", "迟到", "早退", "旷工"]

			if (!monthStatistic && clockIn.length === 0) return

			// 判断打卡是否正常
			let clockInStatus: string = ""
			let isRightClockIn: boolean =
				rightClockInItems.includes(clockIn[2] ?? "") && rightClockInItems.includes(clockIn[4] ?? "")
				||
				wrongClockInItems.includes(clockIn[2] ?? "") || wrongClockInItems.includes(clockIn[4] ?? "")

			// 判断打卡状态
			if (!isRightClockIn) {
				if (clockIn[2] === "缺卡") clockInStatus = "上班缺卡"
				if (clockIn[4] === "缺卡") clockInStatus = "下班缺卡"
				if (clockIn[2] === "迟到") clockInStatus = "上班迟到"
				if (clockIn[4] === "早退") clockInStatus = "下班早退"
				if (clockIn[2] === "旷工" || clockIn[4] === "旷工") clockInStatus = "旷工"
			}

			const matches = String(monthStatistic).match(/(正常|,加班|年假|陪产假|调休|事假|病假|旷工|上班缺卡|下班缺卡|上班迟到|下班早退|休息并打卡)/g) ?? ["未匹配"]

			if (matches.includes("未匹配")) this._analysisResult.get(name)?.push(new EventNotMatch(this.toString(), monthStatistic))
			if (matches.includes("正常")) this._analysisResult.get(name)?.push(new EventDefault(this.toString(), "正常", monthStatistic))
			if (matches.includes(",加班")) {
				if (!isRightClockIn && this._isWeekend) {
					switch (clockInStatus) {
						case "上班缺卡":
							this._analysisResult.get(name)?.push(new EventErrOvertimeCardShortageWorking(this.toString(), monthStatistic))
							break
						case "下班缺卡":
							this._analysisResult.get(name)?.push(new EventErrOvertimeCardShortageClosed(this.toString(), monthStatistic))
							break
					}
				} else {
					if (isRightClockIn && this._isWeekend) this._analysisResult.get(name)?.push(new EventOvertimeWeekend(this.toString(), monthStatistic))
					if (isRightClockIn && this._isHoliday) this._analysisResult.get(name)?.push(new EventOvertimeHoliday(this.toString(), monthStatistic))
				}
			}
			if (matches.includes("年假")) this._analysisResult.get(name)?.push(new EventDefault(this.toString(), "年假", monthStatistic))
			if (matches.includes("陪产假")) this._analysisResult.get(name)?.push(new EventDefault(this.toString(), "陪产假", monthStatistic))
			if (matches.includes("调休")) this._analysisResult.get(name)?.push(new EventDefault(this.toString(), "调休", monthStatistic))
			if (matches.includes("事假")) this._analysisResult.get(name)?.push(new EventDefault(this.toString(), "事假", monthStatistic))
			if (matches.includes("病假")) this._analysisResult.get(name)?.push(new EventDefault(this.toString(), "病假", monthStatistic))
			if (matches.includes("旷工")) this._analysisResult.get(name)?.push(new EventErrAbsenteeism(this.toString(), monthStatistic))
			if (matches.includes("上班缺卡")) this._analysisResult.get(name)?.push(new EventErrWorkdayCardShortageWorking(this.toString(), monthStatistic))
			if (matches.includes("下班缺卡")) this._analysisResult.get(name)?.push(new EventErrWorkdayCardShortageClosed(this.toString(), monthStatistic))
			if (matches.includes("上班迟到")) this._analysisResult.get(name)?.push(new EventErrLate(this.toString(), monthStatistic))
			if (matches.includes("下班早退")) this._analysisResult.get(name)?.push(new EventErrLeaveEarly(this.toString(), monthStatistic))
		})
	}

	// set monthStatistics(name: string, data: string) {
	// 	if (!data) return

	// 	const matches = String(data).match(/(正常|,加班|年假|陪产假|调休|事假|病假|旷工|上班缺卡|下班缺卡|上班迟到|下班早退|休息并打卡)/g) ?? ["未匹配"]
	// 	if (matches.includes("未匹配")) return

	// 	if (!this._logs.has(name)) this._logs.set(name, [])

	// 	if (matches.includes("正常")) this._logs.get(name)?.push(newLogDefault(this.toString(), data))
	// }

	// private _match(value: string): string[] {
	// 	return String(value).match(/(正常|,加班|年假|陪产假|调休|事假|病假|旷工|上班缺卡|下班缺卡|上班迟到|下班早退|休息并打卡)/g) ?? ["未匹配"]
	// }
}