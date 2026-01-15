import { Payment } from './paymentModule';

export const DateType = {
	workday: 'WORKDAY',
	exWorkday: 'EX-WORKDAY',
	exHoliday: 'EX-HOLIDAY',
	weekend: 'WEEKEND',
	holiday: 'HOLIDAY',
	holiday3: 'HOLIDAY-3',
	ignore: 'IGNORE',
};

/**
 * 模型：日期表头
 */
export class DateTitle {
	_type = 'date';
	_dateType = '';
	_idx = 0;
	_label = '';
	_value = '';
	_payment = 1;

	constructor(idx = 0, dateType = '', label = '', value = '') {
		this._type = 'DATE';
		this._idx = idx;
		this._dateType = dateType;
		this._label = label;
		this._value = value;
		this._payment = Payment.new(dateType) || Payment.new(DateType.workday);
	}

	static new = (idx, dateType, label, value) => new DateTitle(idx, dateType, label, value);

	static make = () => new DateTitle();

	get type() {
		return this._type;
	}
	get dateType() {
		return this._dateType;
	}
	get idx() {
		return this._idx;
	}
	get label() {
		return this._label;
	}
	get value() {
		return this._value;
	}
	get payment() {
		return this._payment;
	}

	set type(value) {
		this._type = value;
	}
	set dateType(value) {
		this._dateType = value;
	}
	set idx(value) {
		this._idx = value;
	}
	set label(value) {
		this._label = value;
	}
	set value(value) {
		this._value = value;
	}
	set payment(value) {
		this._payment = value;
	}
}
