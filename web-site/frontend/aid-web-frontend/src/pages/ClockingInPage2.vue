<template>
    <q-page class="q-pa-md q-gutter-md">
        <q-card flat square bordered class="q-mb-md">
            <q-card-section>
                <div class="text-h4 text-deep-orange">考勤统计</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col">
                        输入起始年份（格式：2025）
                        <q-input outlined square label="输入起始年份：" v-model="startYear" type="number" />
                    </div>
                    <div class="col">
                        输入起始月份（格式：01）
                        <q-input outlined square label="输入起始月份：" v-model="startMonth" autofocus ref="txtStartMonth"
                            type="number" min="1" max="12" />
                    </div>
                    <div class="col">
                        输入起始日期（格式：02）
                        <q-input outlined square label="输入起始日期：" v-model="startDay" type="number" min="1" max="31" />
                    </div>
                </div>

                <div class="row q-gutter-md q-mt-sm">
                    <div class="col">
                        <q-card flat square bordered class="my-card">
                            <q-card-section>
                                <q-list>
                                    <q-item clickable>
                                        <q-item-section>
                                            <span style="font-size: 20px" class="text-primary">1. 三薪日</span>
                                            <span class="text-secondary">&nbsp;</span>
                                        </q-item-section>

                                        <q-item-section avatar>
                                            <q-btn-group square>
                                                <q-btn square @click="clearHoliday3Dates" round color="red"
                                                    icon="fa fa-times"></q-btn>
                                                <q-btn square icon="event" round color="primary">
                                                    <q-popup-proxy @before-show="holiday3Dates_updateProxy" cover
                                                        transition-show="scale" transition-hide="scale">
                                                        <q-date v-model="holidays3Proxy" multiple flat bordered
                                                            today-btn landscape years-in-month-view mark="YYYY-MM-DD">
                                                            <div class="row items-center justify-end q-gutter-sm">
                                                                <q-btn-group square>
                                                                    <q-btn square flat label="取消" color="negative"
                                                                        icon="fa fa-times" v-close-popup />
                                                                    <q-btn square flat label="确定" color="green"
                                                                        icon="fa fa-check" @click="holiday3Dates_save"
                                                                        v-close-popup />
                                                                </q-btn-group>
                                                            </div>
                                                        </q-date>
                                                    </q-popup-proxy>
                                                </q-btn>
                                            </q-btn-group>
                                        </q-item-section>
                                    </q-item>
                                </q-list>
                            </q-card-section>

                            <q-separator inset />

                            <q-card-section><date-list :dateList="holidays3"></date-list></q-card-section>
                        </q-card>
                    </div>

                    <div class="col">
                        <q-card flat square bordered class="my-card">
                            <q-card-section>
                                <q-list>
                                    <q-item clickable>
                                        <q-item-section>
                                            <span style="font-size: 20px" class="text-primary">2. 调休</span>
                                            <span class="text-secondary">PS: 放假也得正常上班</span>
                                        </q-item-section>

                                        <q-item-section avatar>
                                            <q-btn-group square>
                                                <q-btn square @click="clearExWorkdayDates" round color="red"
                                                    icon="fa fa-times"></q-btn>
                                                <q-btn square icon="event" round color="primary">
                                                    <q-popup-proxy @before-show="exWorkday_updateProxy" cover
                                                        transition-show="scale" transition-hide="scale">
                                                        <q-date v-model="forceWorkdaysProxy" multiple flat bordered
                                                            today-btn landscape years-in-month-view mark="YYYY-MM-DD">
                                                            <div class="row items-center justify-end q-gutter-sm">
                                                                <q-btn-group square>
                                                                    <q-btn square flat label="取消" color="negative"
                                                                        icon="fa fa-times" v-close-popup />
                                                                    <q-btn square flat label="确定" color="green"
                                                                        icon="fa fa-check" @click="exWorkday_save"
                                                                        v-close-popup />
                                                                </q-btn-group>
                                                            </div>
                                                        </q-date>
                                                    </q-popup-proxy>
                                                </q-btn>
                                            </q-btn-group>
                                        </q-item-section>
                                    </q-item>
                                </q-list>
                            </q-card-section>

                            <q-separator inset />

                            <q-card-section><date-list :dateList="forceWorkdays"></date-list></q-card-section>
                        </q-card>
                    </div>

                    <div class="col">
                        <q-card flat square bordered class="my-card">
                            <q-card-section>
                                <q-list>
                                    <q-item clickable>
                                        <q-item-section>
                                            <span style="font-size: 20px" class="text-primary">3. 年假调休</span>
                                            <span class="text-secondary">PS: 春节公司使用年假全司休息</span>
                                        </q-item-section>

                                        <q-item-section avatar>
                                            <q-btn-group square>
                                                <q-btn square @click="clearExHolidayDates" round color="red"
                                                    icon="fa fa-times"></q-btn>
                                                <q-btn square icon="event" round color="primary">
                                                    <q-popup-proxy @before-show="exHolidayDates_updateProxy" cover
                                                        transition-show="scale" transition-hide="scale">
                                                        <q-date v-model="forceAnnualLeavesProxy" multiple flat bordered
                                                            today-btn landscape years-in-month-view mark="YYYY-MM-DD">
                                                            <div class="row items-center justify-end q-gutter-sm">
                                                                <q-btn-group square>
                                                                    <q-btn square flat label="取消" color="negative"
                                                                        icon="fa fa-times" v-close-popup />
                                                                    <q-btn square flat label="确定" color="green"
                                                                        icon="fa fa-check" @click="exHolidayDates_save"
                                                                        v-close-popup />
                                                                </q-btn-group>
                                                            </div>
                                                        </q-date>
                                                    </q-popup-proxy>
                                                </q-btn>
                                            </q-btn-group>
                                        </q-item-section>
                                    </q-item>
                                </q-list>
                            </q-card-section>

                            <q-separator inset />

                            <q-card-section><date-list :dateList="forceAnnualLeaves"></date-list></q-card-section>
                        </q-card>
                    </div>

                    <div class="col">
                        <q-card flat square bordered class="my-card">
                            <q-card-section>
                                <q-list>
                                    <q-item clickable>
                                        <q-item-section>
                                            <span style="font-size: 20px" class="text-primary">4. 国假</span>
                                            <span class="text-secondary">PS: 春节除三薪日外的假期</span>
                                        </q-item-section>

                                        <q-item-section avatar>
                                            <q-btn-group square>
                                                <q-btn square @click="clearHolidayDates" round color="red"
                                                    icon="fa fa-times"></q-btn>
                                                <q-btn square icon="event" round color="primary">
                                                    <q-popup-proxy @before-show="holidayDates_updateProxy" cover
                                                        transition-show="scale" transition-hide="scale">
                                                        <q-date v-model="holidays2Proxy" multiple flat bordered
                                                            today-btn landscape years-in-month-view mark="YYYY-MM-DD">
                                                            <div class="row items-center justify-end q-gutter-sm">
                                                                <q-btn-group square>
                                                                    <q-btn square flat label="取消" color="negative"
                                                                        icon="fa fa-times" v-close-popup />
                                                                    <q-btn square flat label="确定" color="green"
                                                                        icon="fa fa-check" @click="holidayDates_save"
                                                                        v-close-popup />
                                                                </q-btn-group>
                                                            </div>
                                                        </q-date>
                                                    </q-popup-proxy>
                                                </q-btn>
                                            </q-btn-group>
                                        </q-item-section>
                                    </q-item>
                                </q-list>
                            </q-card-section>

                            <q-separator inset />

                            <q-card-section><date-list :dateList="holidays2"></date-list></q-card-section>
                        </q-card>
                    </div>
                </div>

                <div class="row q-gutter-md q-mt-sm">
                    <div class="col">
                        <q-toggle v-model="isAutoDownload" label="是否自动下载" />&emsp;&emsp;
                        <q-btn square label="手动下载" type="submit" color="primary" @click="generateExcel" />
                        <q-uploader :url="`${API_URL}/v1/checkingIn/cal`" color="purple" label="5. 选择考勤表" flat bordered
                            auto-upload style="width: 100%" :factory="factoryFn" field-name="file"
                            @uploaded="onUploaded" @failed="onUploadFailed" />
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-page>
</template>

<script setup>
import { ref, watch } from 'vue';

import ExcelJS from 'exceljs';
import { saveAs } from 'file-saver';
import moment from 'moment';

import DateList from 'src/components/DateListCom.vue';



import notify from 'src/utils/notify';

import { API_URL, initConfig } from 'src/utils/fetch';


// 起始日期
const now = moment();
let startYear = ref(now.format('YYYY'));
let startMonth = ref(now.format('MM'));
let startDay = ref(26);

const isAutoDownload = ref(false);

watch(startMonth, newVal => {
    if (newVal < 1) startMonth.value = 1;
    if (newVal > 12) startMonth.value = 12;
});
watch(startDay, newVal => {
    if (newVal < 1) startDay.value = 1;
    if (newVal > 31) startDay.value = 31;
});
const txtStartMonth = ref(null);
let finalStatistic = {}; // 最终汇总统计

class statistic {
    _missClockIn = 0 // 上班缺卡
    _missClockOut = 0 // 下班缺卡
    _missClockInOvertime = 0 // 加班上班缺卡
    _missClockOutOvertime = 0 // 加班下班缺卡
    _miss = 0 // 旷工
    _travelRequest = 0 // 出差
    _late = 0 // 迟到
    _early = 0 // 早退
    _reset = 0 // 调休
    _annualLeave = 0 // 年假
    _maternityLeave = 0 // 产假
    _paternityLeave = 0 // 陪产假
    _sickLeave = 0 // 病假
    _personalLeave = 0 // 事假
    _note = [] // 事件记录

    static new = () => new statistic

    get missClockIn() { return this._missClockIn }
    set missClockIn(value) { this._missClockIn += value }
    get missClockOut() { return this._missClockOut }
    set missClockOut(value) { this._missClockOut += value }
    get missClockInOvertime() { return this._missClockInOvertime }
    set missClockInOvertime(value) { this._missClockInOvertime += value }
    get missClockOutOvertime() { return this._missClockOutOvertime }
    set missClockOutOvertime(value) { this._missClockOutOvertime += value }
    get miss() { return this._miss }
    set miss(value) { this._miss += value }
    get travelRequest() { return this._travelRequest }
    set travelRequest(value) { this._travelRequest += value }
    get late() { return this._late }
    set late(value) { this._late += value }
    get early() { return this._early }
    set early(value) { this._early += value }
    get rest() { return this._reset }
    set rest(value) { this._reset += value }
    get annualLeave() { return this._annualLeave }
    set annualLeave(value) { this._annualLeave += value }
    get maternityLeave() { return this._maternityLeave }
    set maternityLeave(value) { this._maternityLeave += value }
    get paternityLeave() { return this._paternityLeave }
    set paternityLeave(value) { this._paternityLeave += value }
    get sickLeave() { return this._sickLeave }
    set sickLeave(value) { this._sickLeave += value }
    get personalLeave() { return this._personalLeave }
    set personalLeave(value) { this._personalLeave += value }
    get note() { return this._note }
    set note(value) { this._note.push(value) }
}


/**
 * 计算结果
 */
const cal = async (originalDates = null, standarDates = null, monthly = null, everyday = null) => {
    var statistics = {}

    for (const date in originalDates) {
        if (!Object.hasOwn(monthly, date)) continue;
        if (statistics.hasOwn(date)) statistics[date] = {};

        const elementForMonthly = monthly[date]
        standarDates
        everyday

        for (const name in elementForMonthly) {
            if (!Object.hasOwn(elementForMonthly, name)) continue;
            if (statistics[date].hasOwn[name]) statistics[date][name] = "";

            const element = elementForMonthly[name];
            const s = statistic.new()

            if (element.includes("上班缺卡")) {
                console.log(`上班缺卡：${date}, ${name}, ${element}`)
                switch (standarDates.standardDateMap[date]["kind"]) {
                    case "WORKDAY":
                    case "FORCE-WORKDAY":
                    case "FORCE-ANNUAL-LEAVE":
                        // 工作日、调休、年假日按照工作日缺卡处理
                        s.missClockIn = 1
                        s.note = `【上班缺卡】${element}`
                        break;
                    case "WEEKEND":
                    case "HOLIDAY2":
                    case "HOLIDAY3":
                        // 周末、节日、三薪日按照加班缺卡处理
                        s.missClockInOvertime = 1
                        s.note = `【加班缺卡】${element}`
                        break;
                }
                continue;
            }

            if (element.includes("下班缺卡")) {
                console.log(`${date}, ${name}, ${element}`)
                switch (standarDates.standardDateMap[date]["kind"]) {
                    case "WORKDAY":
                    case "FORCE-WORKDAY":
                    case "FORCE-ANNUAL-LEAVE":
                        // 工作日、调休、年假日按照工作日缺卡处理
                        s.missClockOut = 1
                        s.note = `【下班缺卡】${element}`
                        break;
                    case "WEEKEND":
                    case "HOLIDAY2":
                    case "HOLIDAY3":
                        // 周末、节日、三薪日按照加班缺卡处理
                        s.missClockOutOvertime = 1
                        s.note = `【加班缺卡】${element}`
                        break;
                }
                continue;
            }

            if (element.includes("早退")) {
                console.log(`${date}, ${name}, ${element}`)
                s.note = `【  早退  】${element}`
                s.early = 1
                continue
            }


            if (element.includes("迟到")) {
                console.log(`${date}, ${name}, ${element}`)
                s.note = `【  迟到  】${element}`
                s.late = 1
                continue
            }

        }
    }
};

const factoryFn = async () => {
    await initConfig()
    // 你可以在这里根据文件信息动态决定参数
    return {
        url: `${API_URL}/v1/checkingIn/cal`, // 上传地址
        method: 'POST',     // 请求方法

        // 1. 设置请求头 (例如 Token)
        headers: [],

        // 2. 设置额外的表单字段 (Form Data)
        // 这些字段会和文件一起被发送到后端
        formFields: [
            { name: 'originalDate', value: `${startYear.value}-${startMonth.value}-${startDay.value}` },
            { name: 'forceWorkdays', value: forceWorkdays.value },
            { name: 'forceAnnualLeaves', value: forceAnnualLeaves.value },
            { name: 'holidays2', value: holidays2.value },
            { name: 'holidays3', value: holidays3.value },
        ],
    }
}

const onUploaded = async (info) => {
    // info.xhr 是原生的 XMLHttpRequest 对象
    // 后端返回的数据通常在 xhr.response 中
    try {
        const response = JSON.parse(info.xhr.response)
        const { everyday, monthly, standardDate } = response.content
        if (!standardDate.hasOwn("originalDate")) {
            notify.error("缺少原始日期信息")
            return
        }
        if (!standardDate.hasOwn("standardDateMap")) {
            notify.error("缺少标准日期信息")
            return
        }
        const { originalDate, standardDateMap } = standardDate

        if (response.code === 200) {
            cal(originalDate, standardDateMap, monthly, everyday)
        } else {
            // 业务逻辑错误（例如：文件类型不对，虽然HTTP状态是200）
            notify.error(response.msg)
        }
    } catch (e) {
        console.error('解析响应失败', e)
        notify.error(`解析响应失败: ${e}`)
    }
}

const onUploadFailed = async (info) => {
    console.error('上传失败:', info)
    var errMsg = ''

    if (info.xhr) {
        // 尝试获取后端返回的错误信息
        try {
            const errResp = JSON.parse(info.xhr.response)
            errMsg = errResp.msg || '服务器错误'
        } catch {
            errMsg = `HTTP ${info.xhr.status}: ${info.xhr.statusText}`
        } finally {
            notify.error(errMsg)
        }
    }
}


/**
 * 生成 Excel 文件
 */
const generateExcel = async () => {
    const workbook = new ExcelJS.Workbook();
    const worksheet = workbook.addWorksheet('Sheet1');
    const defaultBorder = { top: { style: 'thin' }, left: { style: 'thin' }, bottom: { style: 'thin' }, right: { style: 'thin' } };

    // 设置表头
    worksheet.columns = [
        { header: '姓名', key: 'name', width: 10, font: { color: { argb: '000000' } } },
        {
            header: '周末\r\n加班',
            key: 'weekendOvertime',
            width: 10,
            font: { color: { argb: '00FF00' } },
        },
        {
            header: '假日\r\n加班',
            key: 'holidayOvertime',
            width: 10,
            font: { color: { argb: '00FF00' } },
        },
        {
            header: '三薪日\r\n加班',
            key: 'holiday3Overtime',
            width: 10,
            font: { color: { argb: '00FF00' } },
        },
        { header: '年假', key: 'annualLeave', width: 10, font: { color: { argb: '0000FF' } } },
        { header: '陪产假', key: 'paternityLeave', width: 10, font: { color: { argb: '0000FF' } } },
        {
            header: '调休',
            key: 'compensatoryLeave',
            width: 10,
            font: { color: { argb: '0000FF' } },
        },
        { header: '事假', key: 'personalLeave', width: 10, font: { color: { argb: '0000FF' } } },
        { header: '病假', key: 'sickLeave', width: 10, font: { color: { argb: '0000FF' } } },
        { header: '旷工', key: 'absenteeism', width: 10, font: { color: { argb: 'FF0000' } } },
        {
            header: '上班\r\n缺卡',
            key: 'missingClockIn',
            width: 10,
            font: { color: { argb: 'FF0000' } },
        },
        {
            header: '下班\r\n缺卡',
            key: 'missingClockOut',
            width: 10,
            font: { color: { argb: 'FF0000' } },
        },
        {
            header: '上班\r\n迟到',
            key: 'lateClockIn',
            width: 10,
            font: { color: { argb: 'FF0000' } },
        },
        {
            header: '下班\r\n早退',
            key: 'earlyClockOut',
            width: 10,
            font: { color: { argb: 'FF0000' } },
        },
        {
            header: '加班\r\n缺卡',
            key: 'overtimeClockOut',
            width: 10,
            font: { color: { argb: 'FF0000' } },
        },
        // { header: '休息', key: 'reset', width: 16, font: { color: { argb: 'FF0000' } }, },
        { header: '日志', key: 'log', width: 135, font: { color: { argb: '000000' } } },
    ];

    const headerRow = worksheet.getRow(1); // 获取第一行（表头）
    headerRow.font = { name: '仿宋', bold: true, size: 14 };
    headerRow.border = defaultBorder;
    headerRow.height = 42;
    headerRow.eachCell(cell => (cell.alignment = { wrapText: true, vertical: 'middle', horizontal: 'center' }));

    // 添加数据
    Object.entries(finalStatistic).forEach(([name, row]) => {
        worksheet
            .addRow({
                name: name,
                weekendOvertime: row.weekendOvertime || '', // 周末加班
                holidayOvertime: row.holidayOvertime || '', // 假日加班
                holiday3Overtime: row.holiday3Overtime || '', // 三薪加班
                annualLeave: row.annualLeave || '', // 年假
                paternityLeave: row.paternityLeave || '', // 陪产假
                compensatoryLeave: row.compensatoryLeave || '', // 调休
                personalLeave: row.personalLeave || '', // 事假
                sickLeave: row.sickLeave || '', // 病假
                absenteeism: row.absenteeism || '', // 旷工
                missingClockIn: row.missingClockIn || '', // 上班缺卡
                missingClockOut: row.missingClockOut || '', // 下班缺卡
                lateClockIn: row.lateClockIn || '', // 上班迟到
                earlyClockOut: row.earlyClockOut || '', // 下班早退
                overtimeClockOut: row.overtimeClockOut || '', // 加班缺卡
                // reset: row.reset,
                log: row.log.map((item, idx) => `${(idx + 1).toString().padStart(3, '0')}、${item}`).join('\r\n'),
            })
            .eachCell((cell, idx) => {
                let style = {
                    font: { name: '仿宋', size: 12, bold: false, color: { argb: ['000000'] } },
                    border: defaultBorder,
                    alignment: { wrapText: true, vertical: 'middle', horizontal: idx !== 16 ? 'center' : 'left' },
                };

                for (const item of [
                    { target: [2, 3, 4], color: '3D9C6A' }, // 绿色
                    { target: [5, 6, 7, 8, 9], color: '0000FF' }, // 蓝色
                    { target: [10, 11, 12, 13, 14, 15], color: 'FF0000' }, // 红色
                ])
                    if (item.target.includes(idx)) {
                        style.font.color.argb = [item.color];
                        break;
                    }

                if ([2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15].includes(idx)) style.font.bold = true;

                cell.font = style.font;
                cell.border = style.border;
                cell.alignment = style.alignment;
            });

        // 设置冻结第一行和第一列
        worksheet.views = [
            {
                state: 'frozen',
                xSplit: 1, // 冻结第一列
                ySplit: 1, // 冻结第一行
                topLeftCell: 'B2', // 可滚动的区域从B2开始
            },
        ];
    });

    worksheet.eachRow((row, rowNumber) => {
        if (rowNumber % 2 === 0) row.fill = { type: 'pattern', pattern: 'solid', fgColor: { argb: 'D2D2D2' } };
    });

    // 生成 Blob 并下载
    const buffer = await workbook.xlsx.writeBuffer();
    const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
    const filename = `统计：${startYear.value}-${parseInt(startMonth.value) + 1}.xlsx`;
    saveAs(blob, filename);
    notify.ok(`文件已保存：${filename}`);

    txtStartMonth.value?.focus();
};

// 三薪日
const holidays3 = ref([]);
const holidays3Proxy = ref([]);
const holiday3Dates_updateProxy = () => (holidays3Proxy.value = holidays3.value);
const holiday3Dates_save = () => (holidays3.value = holidays3Proxy.value);
const clearHoliday3Dates = () => (holidays3.value = []);
holidays3Proxy.value = [];

// 额外工作日：不算加班（国假倒休）
const forceWorkdays = ref([]);
const forceWorkdaysProxy = ref([]);
const exWorkday_updateProxy = () => (forceWorkdaysProxy.value = forceWorkdays.value);
const exWorkday_save = () => (forceWorkdays.value = forceWorkdaysProxy.value);
const clearExWorkdayDates = () => (forceWorkdays.value = []);
forceWorkdaysProxy.value = [];

// 额外假期：需要用年假补充
const forceAnnualLeaves = ref([]);
const forceAnnualLeavesProxy = ref([]);
const exHolidayDates_updateProxy = () => (forceAnnualLeavesProxy.value = forceAnnualLeaves.value);
const exHolidayDates_save = () => (forceAnnualLeaves.value = forceAnnualLeavesProxy.value);
const clearExHolidayDates = () => (forceAnnualLeaves.value = []);
forceAnnualLeavesProxy.value = [];

// 国假
const holidays2 = ref([]);
const holidays2Proxy = ref([]);
const holidayDates_updateProxy = () => (holidays2Proxy.value = holidays2.value);

const holidayDates_save = () => (holidays2.value = holidays2Proxy.value);
const clearHolidayDates = () => (holidays2.value = []);
holidays2Proxy.value = [];
</script>
