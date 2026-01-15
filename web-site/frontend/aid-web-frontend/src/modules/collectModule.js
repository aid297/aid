import { DateTitle } from './dateTitleModule';

export const CollectType = {
	name: 'NAME',
	collect: 'COLLECT',
	rowNumber: 'ROW-NUMBER',
};

/**
 * 统计数据单项
 */
export class CollectStatistic {
	_normal = false; // 正常
	_overtime = false; // 加班
	_annualLeave = false; // 年假
	_paternityLeave = false; // 陪产假
	_compensatoryLeave = false; // 调休
	_personalLeave = false; // 事假
	_sickLeave = false; // 病假
	_absenteeism = false; // 旷工
	_missingClockIn = false; // 上班缺卡
	_missingClockOut = false; // 下班缺卡
	_lateClockIn = false; // 上班迟到
	_earlyClockOut = false; // 下班早退
	_reset = false; // 休息

	constructor() {}

	static new = () => new CollectStatistic();

	get normal() {
		return this._normal;
	}
	get overtime() {
		return this._overtime;
	}
	get annualLeave() {
		return this._annualLeave;
	}
	get paternityLeave() {
		return this._paternityLeave;
	}
	get compensatoryLeave() {
		return this._compensatoryLeave;
	}
	get personalLeave() {
		return this._personalLeave;
	}
	get sickLeave() {
		return this._sickLeave;
	}
	get absenteeism() {
		return this._absenteeism;
	}
	get missingClockIn() {
		return this._missingClockIn;
	}
	get missingClockOut() {
		return this._missingClockOut;
	}
	get lateClockIn() {
		return this._lateClockIn;
	}
	get earlyClockOut() {
		return this._earlyClockOut;
	}
	get reset() {
		return this._reset;
	}

	set normal(value) {
		this._normal = value;
	}
	set overtime(value) {
		this._overtime = value;
	}
	set annualLeave(value) {
		this._annualLeave = value;
	}
	set paternityLeave(value) {
		this._paternityLeave = value;
	}
	set compensatoryLeave(value) {
		this._compensatoryLeave = value;
	}
	set personalLeave(value) {
		this._personalLeave = value;
	}
	set sickLeave(value) {
		this._sickLeave = value;
	}
	set absenteeism(value) {
		this._absenteeism = value;
	}
	set missingClockIn(value) {
		this._missingClockIn = value;
	}
	set missingClockOut(value) {
		this._missingClockOut = value;
	}
	set lateClockIn(value) {
		this._lateClockIn = value;
	}
	set earlyClockOut(value) {
		this._earlyClockOut = value;
	}
	set reset(value) {
		this._reset = value;
	}
}

/**
 * 模型：打卡数据
 */
export class Collect {
	_type = 'COLLECT';
	_dateTitle = DateTitle.make();
	_contentType = '';
	_value = '';
	_status = false;
	_desc = '';
	_collectStatistic = CollectStatistic.new();
	_print = false;

	constructor(content = null, contentType = null, dateTitle = null, print = false) {
		this._type = 'COLLECT';
		this._dateTitle = dateTitle || DateTitle.make();
		this._value = content || '';
		this._contentType = contentType || '';
		this._print = print;
		this._setValue(content || '');
	}

	static new = (content = null, contentType = null, dateTitle = null, print = false) => new Collect(content, contentType, dateTitle, print);

	static make = () => new Collect();

	get type() {
		return this._type;
	}
	get dateTitle() {
		return this._dateTitle;
	}
	get contentType() {
		return this._contentType;
	}
	get value() {
		return this._value;
	}
	get status() {
		return this._status;
	}
	get statistic() {
		return this._collectStatistic;
	}

	set type(value) {
		this._type = value;
	}
	set dateTitle(value) {
		this._dateTitle = value;
	}
	set contentType(value) {
		this._contentType = value;
	}
	set value(value) {
		this._setValue(value);
	}
	set status(value) {
		this._status = value;
	}
	set statistic(value) {
		this._collectStatistic = value;
	}

	_setValue(value) {
		this._value = value;
		if (this._contentType === CollectType.collect) {
			if (value === '') {
				this._status = false;
				this._desc = '无记录';
			}

			let matches = this._match(value);
			if (matches[0] === '未匹配') {
				this._status = false;
				this._desc = '未匹配';
			}

			if (matches.includes('正常')) {
				this._status = true;
				this._collectStatistic.normal = true;
			}
			if (matches.includes(',加班')) {
				this._status = true;
				this._collectStatistic.overtime = true;
			}
			if (matches.includes('年假')) {
				this._status = true;
				this._collectStatistic.annualLeave = true;
			}
			if (matches.includes('陪产假')) {
				this._status = true;
				this._collectStatistic.paternityLeave = true;
			}
			if (matches.includes('调休')) {
				this._status = true;
				this._collectStatistic.compensatoryLeave = true;
			}
			if (matches.includes('事假')) {
				this._status = true;
				this._collectStatistic.personalLeave = true;
			}
			if (matches.includes('病假')) {
				this._status = true;
				this._collectStatistic.sickLeave = true;
			}
			if (matches.includes('旷工')) {
				this._status = true;
				this._collectStatistic.absenteeism = true;
			}
			if (matches.includes('上班缺卡')) {
				this._status = true;
				this._collectStatistic.missingClockIn = true;
			}
			if (matches.includes('下班缺卡')) {
				this._status = true;
				this._collectStatistic.missingClockOut = true;
			}
			if (matches.includes('上班迟到')) {
				this._status = true;
				this._collectStatistic.lateClockIn = true;
			}
			if (matches.includes('下班早退')) {
				this._status = true;
				this._collectStatistic.earlyClockOut = true;
			}
			// if (matches.includes('休息')) {
			//     this._status = true
			//     this._collectStatistic.reset = true
			// }
			// if (matches.includes('休息并打卡')) {
			//     this._status = true
			//     this._collectStatistic.reset = true
			// }

			if (this._print === true) console.log('#汇总数据：', value, matches, this._collectStatistic);
		}
	}

	_match(value) {
		const pattern = /(正常|,加班|年假|陪产假|调休|事假|病假|旷工|上班缺卡|下班缺卡|上班迟到|下班早退|休息并打卡)/g;

		const matches = String(value).match(pattern);
		return matches ? matches : ['未匹配'];
	}
}
