// import { Str } from 'src/utils/str'
import moment from "moment";

export const ClockInType = {
    name: "NAME",
    clockIn: "CLOCK-IN",
    rowNumber: "ROW-NUMBER",
};

/**
 * 模型：打卡数据
 */
export class ClockIn {
    _type = "CLOCK-IN";
    _dateTitle = null;
    _contentType = "";
    _value = "";
    _status = false;
    _desc = "";
    _print = false;
    _s = "";

    constructor(content = null, contentType = null, dateTitle = null, print = false) {
        this._type = "CLOCK-IN";
        this._dateTitle = dateTitle;
        this._value = content || "";
        this._contentType = contentType || "";
        this._desc = "";
        this._print = print;
        this._setValue(content || "");
    }

    static new = (content = null, contentType = null, dateTitle = null, print = false) => new ClockIn(content, contentType, dateTitle, print);

    static make = () => new ClockIn();

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
    get desc() {
        return this._desc;
    }
    get print() {
        return this._print;
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
    set desc(value) {
        this._desc = value;
    }
    set print(value) {
        this._print = value;
    }

    _setValue(value = "") {
        this._value = value;
        if (this._contentType === ClockInType.clockIn) {
            let isAm = false;
            let isPm = false;

            if (value === "") {
                this._status = false;
                this._desc = "未打卡";
                return;
            }

            const matches = value.match(/(?:[01]?\d|2[0-3]):([0-5]\d)/g);
            if (matches === null || matches.length === 0) {
                if (["美国", "加拿大", "墨西哥", "巴西", "阿根廷", "智利", "英国", "法国", "德国", "意大利", "俄罗斯", "波兰", "埃及", "南非", "沙特阿拉伯", "土耳其", "印度", "印度尼西亚", "巴基斯坦", "孟加拉国", "泰国", "日本", "韩国", "澳大利亚", "新西兰"].includes(value)) {
                    this._status = true;
                    this._desc = `特殊打卡（${value}）`;
                    return;
                }

                this._status = false;
                this._desc = `打卡数据格式错误：${value}`;
                return;
            }

            matches.forEach(match => {
                const hour = moment(match, "HH:mm").format("HH");
                if (hour < 12) isAm = true;
                if (hour >= 12) isPm = true;
            });

            if (isAm && isPm) {
                this._status = true;
                this._desc = "正常打卡";
            } else if (isAm && !isPm) {
                this._desc = "缺少下班卡";
            } else if (!isAm && isPm) {
                this._desc = "缺少上班卡";
            } else {
                this._desc = "缺卡";
            }
        }

        if (this._print) console.log(`#打卡数据：${this._value} - 状态: ${this._status} - 描述: ${this._desc}`);
    }
}
