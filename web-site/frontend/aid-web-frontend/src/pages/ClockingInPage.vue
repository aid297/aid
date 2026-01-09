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
                                                        <q-date v-model="holiday3DatesProxy" multiple flat bordered
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

                            <q-card-section><date-list :dateList="holiday3Dates"></date-list></q-card-section>
                        </q-card>
                    </div>

                    <div class="col">
                        <q-card flat square bordered class="my-card">
                            <q-card-section>
                                <q-list>
                                    <q-item clickable>
                                        <q-item-section>
                                            <span style="font-size: 20px" class="text-primary">2. 调休</span>
                                            <span class="text-secondary">PS: 调休按照工作日计算</span>
                                        </q-item-section>

                                        <q-item-section avatar>
                                            <q-btn-group square>
                                                <q-btn square @click="clearExWorkdayDates" round color="red"
                                                    icon="fa fa-times"></q-btn>
                                                <q-btn square icon="event" round color="primary">
                                                    <q-popup-proxy @before-show="exWorkday_updateProxy" cover
                                                        transition-show="scale" transition-hide="scale">
                                                        <q-date v-model="exWorkdayDatesProxy" multiple flat bordered
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

                            <q-card-section><date-list :dateList="exWorkdayDates"></date-list></q-card-section>
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
                                                        <q-date v-model="exHolidayDatesProxy" multiple flat bordered
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

                            <q-card-section><date-list :dateList="exHolidayDates"></date-list></q-card-section>
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
                                                        <q-date v-model="holidayDatesProxy" multiple flat bordered
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

                            <q-card-section><date-list :dateList="holidayDates"></date-list></q-card-section>
                        </q-card>
                    </div>
                </div>

                <div class="row q-gutter-md q-mt-sm">
                    <div class="col">
                        <q-toggle v-model="isAutoDownloadStatisticsExcel" label="是否自动下载" />&emsp;&emsp;
                        <q-btn square label="手动下载" type="submit" color="primary" @click="generateExcel" />
                        <q-uploader url="#" color="purple" label="5. 选择考勤表" flat bordered style="width: 100%"
                            @added="onUploadedClockInExcel" />
                    </div>
                    <div class="col">
                        <p>&nbsp;</p>
                        <q-uploader url="#" color="green-10" label="6. 选择统计表" flat bordered style="width: 100%"
                            @added="onUploadedStatisticExcel" />
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-page>
</template>

<script setup>
import { ref, watch } from "vue";

import ExcelJS from "exceljs";
import { saveAs } from "file-saver";
import moment from "moment";

import DateList from "src/components/DateListCom.vue";

import { ClockInService } from "src/services/clockInService";
import { CollectService } from "src/services/collectService";
import { DateTitleService } from "src/services/dateTitleService";
import { StatisticService } from "src/services/statisticService";

import { EveryDaysModule } from "src/modules/everydayModule";
import { Str } from "src/utils/str";

import notify from "src/utils/notify";

// 起始日期
const now = moment();
let startYear = ref(now.format("YYYY"));
let startMonth = ref(now.format("MM"));
let startDay = ref(26);

const isAutoDownloadStatisticsExcel = ref(false);

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

const onUploadedClockInExcel = async files => {
    const file = files[0];
    const originalClockIn = [];
    const originalCollect = [];
    const everydayData = {};
    let dateTitle = null; // 日期表头
    let clockInData = {}; // 打卡记录数据
    let collectData = {}; // 汇总数据

    const fileReader = new FileReader();
    fileReader.onload = async event => {
        try {
            // 获取文件的ArrayBuffer数据
            const arrayBuffer = event.target?.result;

            // 使用ExcelJS加载文件
            const workbook = new ExcelJS.Workbook();
            await workbook.xlsx.load(arrayBuffer);

            const clockInSheet = workbook.getWorksheet("打卡时间");
            const collectSheet = workbook.getWorksheet("月度汇总");
            const everydaySheet = workbook.getWorksheet("每日统计");

            if (clockInSheet === undefined) {
                notify.error("工作表 “打卡时间” 不存在");
                return;
            }

            if (collectSheet === undefined) {
                notify.error("工作表“月度汇总”不存在");
                return;
            }

            if (everydaySheet === undefined) {
                notify.error("工作表 “每日统计” 不存在");
                console.error("工作表 “每日统计” 不存在");
                return;
            }

            // 获取每日统计数据 -> 每日打卡
            everydaySheet.eachRow(row => {
                const values = row.values;
                const name = Str.new(values[1]).replace();
                if (!(name in everydayData)) everydayData[name] = EveryDaysModule.new();

                everydayData[name].push(values);
            });

            // 获取初始数据 -> 打卡
            clockInSheet.eachRow(row => {
                const values = row.values;
                values[0] = row.number;
                originalClockIn.push(values);
            });

            // 获取初始数据 -> 汇总
            collectSheet.eachRow(row => {
                const values = row.values;
                values[0] = row.number;
                originalCollect.push(values);
            });

            // 解析标题 -> 日期
            dateTitle = DateTitleService.new(originalClockIn[0]).parse(holiday3Dates.value, exWorkdayDates.value, holidayDates.value, exHolidayDates.value, `${startYear.value}-${startMonth.value}-${startDay.value}`);

            // 解析数据 -> 打卡
            ClockInService.new(originalClockIn.slice(1), dateTitle)
                .parse()
                .data.forEach(item => { clockInData[item[1].value] = item; });

            // 解析数据 -> 汇总
            CollectService.new(originalCollect.slice(1), dateTitle)
                .parse()
                .data.forEach(item => { collectData[item[1].value] = item; });

            // 分析数据
            finalStatistic = StatisticService.new(dateTitle.data, clockInData, collectData, everydayData).parse().data;

            console.log("# 获取考勤数据（分析后）", finalStatistic);

            if (isAutoDownloadStatisticsExcel.value) await generateExcel()  // 保存统计文件到本地
        } catch (error) {
            console.error("解析Excel失败：", error);
        }
    };

    fileReader.readAsArrayBuffer(file);
};

const onUploadedStatisticExcel = async files => {
    if (!finalStatistic) {
        notify.error("请先上传考勤表进行统计");
        return;
    }

    const columnMap = {
        1: {
            overtime: "G",
            compensatoryLeave: "H",
            note: "J",
        }, 2: {
            overtime: "K",
            compensatoryLeave: "L",
            note: "N",
        },
        3: {
            overtime: "O",
            compensatoryLeave: "P",
            note: "U",
        },
        4: {
            overtime: "V",
            compensatoryLeave: "W",
            note: "Y",
        },
        5: {
            overtime: "Z",
            compensatoryLeave: "AA",
            note: "AC",
        },
        6: {
            overtime: "AD",
            compensatoryLeave: "AE",
            note: "AG",
        },
        7: {
            overtime: "AH",
            compensatoryLeave: "AI",
            note: "AK",
        },
        8: {
            overtime: "AL",
            compensatoryLeave: "AM",
            note: "AO",
        },
        9: {
            overtime: "AP",
            compensatoryLeave: "AQ",
            note: "AS",
        },
        10: {
            overtime: "AT",
            compensatoryLeave: "AU",
            note: "AW",
        },
        11: {
            overtime: "AX",
            compensatoryLeave: "AY",
            note: "BA",
        },
        12: {
            overtime: "BB",
            compensatoryLeave: "BC",
            note: "BE",
        },
    }

    const nameToLineNum = {};

    const file = files[0];
    const fileReader = new FileReader();

    fileReader.onload = async event => {
        // 获取文件的ArrayBuffer数据
        const arrayBuffer = event.target?.result;

        // 使用ExcelJS加载文件
        const workbook = new ExcelJS.Workbook();
        await workbook.xlsx.load(arrayBuffer);

        const originalSheet = workbook.getWorksheet(startYear.value);

        originalSheet.eachRow(row => { if (row.number < 4) return; nameToLineNum[Str.new(row.values[3]).replace()] = row.number; });

        // 更新数据
        for (const [name, value] of Object.entries(finalStatistic)) {
            let overtime = 0;
            let compensatoryLeave = 0;
            let log = "";
            const lineNum = nameToLineNum[name]; if (!lineNum) continue;
            const col = columnMap[parseInt(startMonth.value) + 1]; if (!col) continue;

            // 计算加班
            overtime = value.weekendOvertime + value.holidayOvertime + (value.holiday3Overtime) * 3;
            // 计算调休：年假、陪产假、调休、事假、旷工、病假
            compensatoryLeave = value.annualLeave + value.paternityLeave + value.compensatoryLeave + value.absenteeism + value.personalLeave;
            // 日志
            log = value.log.map((item, idx) => `${(idx + 1).toString().padStart(3, "0")}、${item}`).join("\r\n");

            const logLines = value.log.length;
            const currentRow = originalSheet.getRow(lineNum);
            currentRow.height = Math.max(25, logLines * 16); // 每行约16像素，最小25

            const colNameOvertime = `${col.overtime}${lineNum}`;
            const colNameCompensatoryLeave = `${col.compensatoryLeave}${lineNum}`;
            const colNameNote = `${col.note}${lineNum}`;

            originalSheet.getColumn(col.note).width = 140;
            originalSheet.getCell(colNameOvertime).value = overtime;
            originalSheet.getCell(colNameCompensatoryLeave).value = compensatoryLeave;
            const cell = originalSheet.getCell(colNameNote);
            cell.value = log;
            cell.alignment = { wrapText: true, vertical: "middle", horizontal: "left" };
        }

        const buffer = await workbook.xlsx.writeBuffer();
        const blob = new Blob([buffer], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });
        const filename = `统计2：${startYear.value}-${parseInt(startMonth.value) + 1}.xlsx`
        saveAs(blob, filename);
        notify.ok(`文件已保存：${filename}`);
    };

    fileReader.readAsArrayBuffer(file);
};

/**
 * 生成 Excel 文件
 */
const generateExcel = async () => {
    const workbook = new ExcelJS.Workbook()
    const worksheet = workbook.addWorksheet("Sheet1")
    const defaultBorder = { top: { style: "thin" }, left: { style: "thin" }, bottom: { style: "thin" }, right: { style: "thin" } }

    // 设置表头
    worksheet.columns = [
        { header: "姓名", key: "name", width: 10, font: { color: { argb: "000000" } } },
        {
            header: "周末\r\n加班",
            key: "weekendOvertime",
            width: 10,
            font: { color: { argb: "00FF00" } },
        },
        {
            header: "假日\r\n加班",
            key: "holidayOvertime",
            width: 10,
            font: { color: { argb: "00FF00" } },
        },
        {
            header: "三薪日\r\n加班",
            key: "holiday3Overtime",
            width: 10,
            font: { color: { argb: "00FF00" } },
        },
        { header: "年假", key: "annualLeave", width: 10, font: { color: { argb: "0000FF" } } },
        { header: "陪产假", key: "paternityLeave", width: 10, font: { color: { argb: "0000FF" } } },
        {
            header: "调休",
            key: "compensatoryLeave",
            width: 10,
            font: { color: { argb: "0000FF" } },
        },
        { header: "事假", key: "personalLeave", width: 10, font: { color: { argb: "0000FF" } } },
        { header: "病假", key: "sickLeave", width: 10, font: { color: { argb: "0000FF" } } },
        { header: "旷工", key: "absenteeism", width: 10, font: { color: { argb: "FF0000" } } },
        {
            header: "上班\r\n缺卡",
            key: "missingClockIn",
            width: 10,
            font: { color: { argb: "FF0000" } },
        },
        {
            header: "下班\r\n缺卡",
            key: "missingClockOut",
            width: 10,
            font: { color: { argb: "FF0000" } },
        },
        {
            header: "上班\r\n迟到",
            key: "lateClockIn",
            width: 10,
            font: { color: { argb: "FF0000" } },
        },
        {
            header: "下班\r\n早退",
            key: "earlyClockOut",
            width: 10,
            font: { color: { argb: "FF0000" } },
        },
        {
            header: "加班\r\n缺卡",
            key: "overtimeClockOut",
            width: 10,
            font: { color: { argb: "FF0000" } },
        },
        // { header: '休息', key: 'reset', width: 16, font: { color: { argb: 'FF0000' } }, },
        { header: "日志", key: "log", width: 135, font: { color: { argb: "000000" } } },
    ]

    const headerRow = worksheet.getRow(1) // 获取第一行（表头）
    headerRow.font = { name: "仿宋", bold: true, size: 14 }
    headerRow.border = defaultBorder
    headerRow.height = 42
    headerRow.eachCell(cell => (cell.alignment = { wrapText: true, vertical: "middle", horizontal: "center" }))

    // 添加数据
    Object.entries(finalStatistic).
        forEach(([name, row]) => {
            worksheet
                .addRow({
                    name: name,
                    weekendOvertime: row.weekendOvertime || "",  // 周末加班
                    holidayOvertime: row.holidayOvertime || "",  // 假日加班
                    holiday3Overtime: row.holiday3Overtime || "",  // 三薪加班
                    annualLeave: row.annualLeave || "",  // 年假
                    paternityLeave: row.paternityLeave || "",  // 陪产假
                    compensatoryLeave: row.compensatoryLeave || "",  // 调休
                    personalLeave: row.personalLeave || "",// 事假
                    sickLeave: row.sickLeave || "",  // 病假
                    absenteeism: row.absenteeism || "", // 旷工
                    missingClockIn: row.missingClockIn || "", // 上班缺卡
                    missingClockOut: row.missingClockOut || "",  // 下班缺卡
                    lateClockIn: row.lateClockIn || "",  // 上班迟到
                    earlyClockOut: row.earlyClockOut || "",  // 下班早退
                    overtimeClockOut: row.overtimeClockOut || "", // 加班缺卡
                    // reset: row.reset,
                    log: row.log.map((item, idx) => `${(idx + 1).toString().padStart(3, "0")}、${item}`).join("\r\n"),
                })
                .eachCell((cell, idx) => {
                    let style = {
                        font: { name: "仿宋", size: 12, bold: false, color: { argb: ["000000"] } },
                        border: defaultBorder,
                        alignment: { wrapText: true, vertical: "middle", horizontal: idx !== 16 ? "center" : "left" },
                    }

                    for (const item of [
                        { target: [2, 3, 4], color: "3D9C6A" }, // 绿色
                        { target: [5, 6, 7, 8, 9], color: "0000FF" }, // 蓝色
                        { target: [10, 11, 12, 13, 14, 15], color: "FF0000" }, // 红色
                    ])
                        if (item.target.includes(idx)) {
                            style.font.color.argb = [item.color]
                            break
                        }

                    if ([2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15].includes(idx)) style.font.bold = true

                    cell.font = style.font
                    cell.border = style.border
                    cell.alignment = style.alignment
                })

            // 设置冻结第一行和第一列
            worksheet.views = [
                {
                    state: "frozen",
                    xSplit: 1, // 冻结第一列
                    ySplit: 1, // 冻结第一行
                    topLeftCell: "B2", // 可滚动的区域从B2开始
                },
            ]
        })

    worksheet.eachRow((row, rowNumber) => {
        if (rowNumber % 2 === 0) row.fill = { type: "pattern", pattern: "solid", fgColor: { argb: "D2D2D2" } }
    })

    // 生成 Blob 并下载
    const buffer = await workbook.xlsx.writeBuffer()
    const blob = new Blob([buffer], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" })
    const filename = `统计：${startYear.value}-${parseInt(startMonth.value) + 1}.xlsx`;
    saveAs(blob, filename)
    notify.ok(`文件已保存：${filename}`)

    txtStartMonth.value?.focus()
}

// 三薪日
const holiday3Dates = ref([]);
const holiday3DatesProxy = ref([]);
const holiday3Dates_updateProxy = () => (holiday3DatesProxy.value = holiday3Dates.value);
const holiday3Dates_save = () => (holiday3Dates.value = holiday3DatesProxy.value);
const clearHoliday3Dates = () => (holiday3Dates.value = []);
holiday3DatesProxy.value = [];

// 额外工作日：不算加班（国假倒休）
const exWorkdayDates = ref([]);
const exWorkdayDatesProxy = ref([]);
const exWorkday_updateProxy = () => (exWorkdayDatesProxy.value = exWorkdayDates.value);
const exWorkday_save = () => (exWorkdayDates.value = exWorkdayDatesProxy.value);
const clearExWorkdayDates = () => (exWorkdayDates.value = []);
exWorkdayDatesProxy.value = [];

// 额外假期：需要用年假补充
const exHolidayDates = ref([]);
const exHolidayDatesProxy = ref([]);
const exHolidayDates_updateProxy = () => (exHolidayDatesProxy.value = exHolidayDates.value);
const exHolidayDates_save = () => (exHolidayDates.value = exHolidayDatesProxy.value);
const clearExHolidayDates = () => (exHolidayDates.value = []);
exHolidayDatesProxy.value = [];

// 国假
const holidayDates = ref([]);
const holidayDatesProxy = ref([]);
const holidayDates_updateProxy = () => (holidayDatesProxy.value = holidayDates.value);

const holidayDates_save = () => (holidayDates.value = holidayDatesProxy.value);
const clearHolidayDates = () => (holidayDates.value = []);
holidayDatesProxy.value = [];
</script>
