export class EventDefault {
	time: string
	type: string
	original: string
	constructor(time: string, type: string, original: string) {
		this.time = time
		this.type = type
		this.original = original
	}
}

export class EventError extends EventDefault { }
export class EventNotMatch extends EventError { constructor(time: string, original: string) { super(time, "不匹配", original) } }
export class EventErrLate extends EventError { constructor(time: string, original: string) { super(time, "迟到", original) } }
export class EventErrLeaveEarly extends EventError { constructor(time: string, original: string) { super(time, "早退", original) } }
export class EventErrWorkdayCardShortageWorking extends EventError { constructor(time: string, original: string) { super(time, "工作日缺卡(上班)", original) } }
export class EventErrWorkdayCardShortageClosed extends EventError { constructor(time: string, original: string) { super(time, "工作日缺卡(下班)", original) } }
export class EventErrOvertimeCardShortageWorking extends EventError {constructor(time: string, original: string) {super(time,"加班缺卡(上班)", original)}}
export class EventErrOvertimeCardShortageClosed extends EventError {constructor(time: string, original: string) {super(time,"加班缺卡(下班)", original)}}
export class EventErrAbsenteeism extends EventError {constructor(time: string, original: string) {super(time,"旷工", original)}}

export class EventOvertime extends EventDefault { }
export class EventOvertimeWeekend extends EventOvertime {constructor(time: string, original: string) {super(time,"周末加班", original)}}
export class EventOvertimeHoliday extends EventOvertime {constructor(time: string, original: string) {super(time,"节假日加班", original)}}