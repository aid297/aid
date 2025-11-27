export const payments = { WORKDAY: 1, WEEKEND: 1, HOLIDAY: 1, "HOLIDAY-3": 3, IGNORE: 0 };

export class Payment {
    _data = 1;

    constructor(dateType) {
        this._data = payments[dateType] || 1;
    }

    static new = dateType => new Payment(dateType);

    get value() {
        return this._data;
    }

    set value(value) {
        this._data = value;
    }
}
