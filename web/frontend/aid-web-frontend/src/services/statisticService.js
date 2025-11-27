import { DateType } from "src/modules/dateTitleModule";
import notify from "src/utils/notify";

/**
 * 服务：统计
 */
export class StatisticService {
    _dateTitles = [];
    _clockInData = {};
    _collectData = {};
    _finalStatistic = {};

    constructor(dateTitles = null, clockInData = null, collectData = null, everydayData = null) {
        this._dateTitles = dateTitles || [];
        this._clockInData = clockInData || {};
        this._collectData = collectData || {};
        this._everydayData = everydayData || {};
    }

    static new = (dateTitles = null, clockInData = null, collectData = null, everydayData = null) => new StatisticService(dateTitles, clockInData, collectData, everydayData);

    static make = () => new StatisticService();

    get dateTitles() {
        return this._dateTitles;
    }
    get clockInData() {
        return this._clockInData;
    }
    get collectData() {
        return this._collectData;
    }
    get data() {
        return this._finalStatistic;
    }

    set dateTitles(value) {
        this._dateTitles = value;
    }
    set clockInData(value) {
        this._clockInData = value;
    }
    set collectData(value) {
        this._collectData = value;
    }

    _check() {
        if (this._dateTitles.length === 0) {
            notify.error("日期表头错误");
            return false;
        }

        if (!this._clockInData) {
            notify.error("打卡数据错误");
            console.error("打卡数据错误", this._clockInData);
            return false;
        }

        if (!this._collectData) {
            notify.error("汇总数据错误");
            console.error("汇总数据错误", this._collectData);
            return false;
        }
        return true;
    }

    parse() {
        if (!this._check()) return this;

        this._dateTitles.forEach(dateTitle => {
            Object.entries(this._collectData).forEach(([name, collectDatum]) => {
                if (!this._finalStatistic[name])
                    this._finalStatistic[name] = {
                        weekendOvertime: 0, // 周末加班
                        holidayOvertime: 0, // 假日加班
                        holiday3Overtime: 0, // 三薪日加班
                        annualLeave: 0, // 年假
                        paternityLeave: 0, // 陪产假
                        compensatoryLeave: 0, // 调休
                        personalLeave: 0, // 事假
                        sickLeave: 0, // 病假
                        absenteeism: 0, // 旷工
                        missingClockIn: 0, // 上班缺卡
                        missingClockOut: 0, // 下班缺卡
                        lateClockIn: 0, // 上班迟到
                        earlyClockOut: 0, // 下班早退
                        overtimeClockOut: 0, // 加班缺卡
                        reset: 0, // 休息
                        log: [], // 日志
                        dateTitle: dateTitle, // 日期类型
                        collectData: [],
                        clockInData: [],
                    };

                let collectDatumItem = collectDatum[dateTitle.idx];
                let clockInDatumItem = this._clockInData[name][dateTitle.idx];

                // // todo: 调试用
                // if (name === "李伟") {
                //   console.log(dateTitle.value)
                //   console.log("日期", dateTitle)
                //   console.log("统计", collectDatumItem)
                //   console.log("打卡", clockInDatumItem)
                // }

                this._finalStatistic[name].collectData.push(collectDatumItem);
                this._finalStatistic[name].clockInData.push(clockInDatumItem);

                // 处理数据
                // 国假加班
                if (dateTitle.dateType === DateType.holiday && collectDatumItem.statistic.overtime) {
                    if (clockInDatumItem.status) {
                        this._finalStatistic[name].holidayOvertime++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】国假加班 → ${collectDatumItem.value}`);
                    } else {
                        if (!this._everydayData[name].find(dateTitle.value)) {
                            this._finalStatistic[name].overtimeClockOut++;
                            this._finalStatistic[name].log.push(`【${dateTitle.value}】加班缺卡 → ${clockInDatumItem.desc}`);
                            return;
                        }

                        this._finalStatistic[name].holidayOvertime++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】国假加班 → ${collectDatumItem.value}`);
                    }
                }

                // 周末加班
                if (dateTitle.dateType === DateType.weekend && collectDatumItem.statistic.overtime) {
                    if (clockInDatumItem.status) {
                        this._finalStatistic[name].weekendOvertime++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】周末加班 → ${collectDatumItem.value}`);
                    } else {
                        if (!this._everydayData[name].find(dateTitle.value)) {
                            this._finalStatistic[name].overtimeClockOut++;
                            this._finalStatistic[name].log.push(`【${dateTitle.value}】加班缺卡 → ${clockInDatumItem.desc}`);
                            return;
                        }

                        this._finalStatistic[name].weekendOvertime++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】周末加班 → ${collectDatumItem.value}`);
                    }
                }

                // 三薪日加班
                if (dateTitle.dateType === DateType.holiday3 && collectDatumItem.statistic.overtime) {
                    if (clockInDatumItem.status) {
                        this._finalStatistic[name].holiday3Overtime++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】三薪加班 → ${collectDatumItem.value}`);
                    } else {
                        if (!this._everydayData[name].find(dateTitle.value)) {
                            this._finalStatistic[name].overtimeClockOut++;
                            this._finalStatistic[name].log.push(`【${dateTitle.value}】加班缺卡 → ${clockInDatumItem.desc}`);
                            return;
                        }

                        this._finalStatistic[name].holiday3Overtime++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】三薪加班 → ${collectDatumItem.value}`);
                    }
                }

                // 额外假期：用年假补充
                if (dateTitle.dateType === DateType.exHoliday && collectDatumItem.statistic.normal) {
                    if (clockInDatumItem.status) {
                        this._finalStatistic[name].annualLeave++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】年假调休 → ${collectDatumItem.value}`);
                    } else {
                        this._finalStatistic[name].absenteeism++;
                        this._finalStatistic[name].log.push(`【${dateTitle.value}】旷工 → ${clockInDatumItem.desc}`);
                    }
                }

                // 年假
                if (collectDatumItem.statistic.annualLeave) {
                    this._finalStatistic[name].annualLeave++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】年假     → ${collectDatumItem.value}`);
                }

                // 陪产假
                if (collectDatumItem.statistic.paternityLeave) {
                    this._finalStatistic[name].paternityLeave++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】陪产假   → ${collectDatumItem.value}`);
                }

                // 调休
                if (collectDatumItem.statistic.compensatoryLeave) {
                    this._finalStatistic[name].compensatoryLeave++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】调休     → ${collectDatumItem.value}`);
                }

                // 事假
                if (collectDatumItem.statistic.personalLeave) {
                    this._finalStatistic[name].personalLeave++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】事假     → ${collectDatumItem.value}`);
                }

                // 病假
                if (collectDatumItem.statistic.sickLeave) {
                    this._finalStatistic[name].sickLeave++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】病假     → ${collectDatumItem.value}`);
                }

                // 旷工
                if ([DateType.workday, DateType.exWorkday].includes(dateTitle.dateType) && collectDatumItem.statistic.absenteeism) {
                    this._finalStatistic[name].absenteeism++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】旷工     → ${collectDatumItem.value}`);
                }

                // 上班缺卡
                if (collectDatumItem.statistic.missingClockIn) {
                    this._finalStatistic[name].missingClockIn++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】上班缺卡 → ${collectDatumItem.value}`);
                }

                // 下班缺卡
                if (collectDatumItem.statistic.missingClockOut) {
                    this._finalStatistic[name].missingClockOut++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】下班缺卡 → ${collectDatumItem.value}`);
                }

                // 上班迟到
                if (collectDatumItem.statistic.lateClockIn) {
                    this._finalStatistic[name].lateClockIn++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】上班迟到 → ${collectDatumItem.value}`);
                }

                // 下班早退
                if (collectDatumItem.statistic.earlyClockOut) {
                    this._finalStatistic[name].earlyClockOut++;
                    this._finalStatistic[name].log.push(`【${dateTitle.value}】下班早退 → ${collectDatumItem.value}`);
                }
            });
        });

        return this;
    }
}
