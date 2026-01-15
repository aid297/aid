import { DateTitle } from './dateTitleModule';
import { ClockIn } from './clockInModule';
import { Collect } from './collectModule';

/**
 * 模型：统计
 */
export class Statistic {
	_dateTitle = DateTitle.make();
	_clockIn = ClockIn.make();
	_collect = Collect.make();

	constructor(dateTitle = null, clockIn = null, collect = null) {
		this._dateTitle = dateTitle || DateTitle.make();
		this._clockIn = clockIn || ClockIn.make();
		this._collect = collect || Collect.make();
	}

	static new = (dateTitle = null, clockIn = null, collect = null) => new Statistic(dateTitle, clockIn, collect);

	static make = () => new Statistic();

	get dateTitle() {
		return this._dateTitle;
	}
	get clockIn() {
		return this._clockIn;
	}
	get collect() {
		return this._collect;
	}

	set dateTitle(value) {
		this._dateTitle = value;
	}
	set clockIn(value) {
		this._clockIn = value;
	}
	set collect(value) {
		this._collect = value;
	}
}
