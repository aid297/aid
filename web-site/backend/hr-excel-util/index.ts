import * as fs from 'fs'
import moment from 'moment'
import { Reader } from "./excel/reader/reader"
import { StatisticModule } from './modules/statistic.module'

// 读取本地Excel文件的函数
// async function readLocalExcel(filePath: string): Promise<any[]> {
// 	try {
// 		// 1. 读取文件内容（二进制格式）
// 		const fileBuffer = await fs.readFile(filePath);

// 		// 2. 解析为工作簿（Workbook）
// 		const workbook = XLSX.read(fileBuffer, { type: 'buffer' });

// 		// 3. 获取第一个工作表（SheetNames是工作表名称数组）
// 		const firstSheetName = workbook.SheetNames[0];
// 		const worksheet = firstSheetName ? workbook.Sheets[firstSheetName] : undefined;
// 		if (!worksheet) {
// 			throw new Error('工作簿中没有有效的工作表');
// 		}

// 		// 4. 转换为JSON格式（header: 1表示按行输出二维数组）
// 		const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

// 		console.log('解析后的数据：', jsonData);
// 		return jsonData;
// 	} catch (error) {
// 		console.error('读取Excel文件失败：', error);
// 		throw error; // 抛出错误供上层处理
// 	}
// }

// 使用示例（替换为你的本地文件路径）
const [excelFilePath, startMonth] = process.argv.slice(2)
if (!excelFilePath) {
	console.error("请提供Excel文件路径")
	process.exit(1)
}

if (!startMonth) {
	console.error("请提供起始月份")
	process.exit(1)
}

const workdays = JSON.parse(fs.readFileSync(`${excelFilePath}.workday.json`, "utf-8"))
if (!workdays) {
	console.error("读取工作日失败")
	process.exit(1)
}
const annualVacations = JSON.parse(fs.readFileSync(`${excelFilePath}.annualVacation.json`, "utf-8"))
if (!annualVacations) {
	console.error("读取年假失败")
	process.exit(1)
}

const buffer = fs.readFileSync(`${excelFilePath}.xlsx`)
const reader = Reader.new(buffer)
let clockInStatisticOriginalData: string[][] = reader.loadSheet("每日统计")
let monthlyCollectionOriginalData: string[][] = reader.loadSheet("月度汇总")

clockInStatisticOriginalData.slice(4).map(row => { row[0] ? row[0] = row[0].replaceAll(/ /g, "") : row[0] = ""; return row })
monthlyCollectionOriginalData.slice(4).map(row => { row[0] ? row[0] = row[0].replaceAll(/ /g, "") : row[0] = ""; return row })

fs.writeFile(`${excelFilePath}.每日统计.json`, JSON.stringify(clockInStatisticOriginalData, null, 4), "utf-8", () => { })
fs.writeFile(`${excelFilePath}.月度汇总.json`, JSON.stringify(monthlyCollectionOriginalData, null, 4), "utf-8", () => { })

/**
 * 处理标准时间表头
 * @param clockInOriginalData 打卡时间元数据
 * @returns {StatisticModule[]} 时间统计模块数组
 */
const parseStatisticTimeHeader = (clockInOriginalData: string[][]): StatisticModule[] => {
	const header: string[] = clockInOriginalData[3] ?? []
	const exec: StatisticModule[] = []

	if (!header) throw new Error("处理【打卡时间】的时间表头错误：没有表头")

	let now: moment.Moment

	header.slice(6).forEach((datum, idx) => {
		if (idx === 1) {
			if (["六", "日"].includes(datum)) {
				now = moment(`${startMonth}-26`)
			} else {
				now = moment(`${startMonth}-${datum}`)
			}
		} else {
			now.add(1, "day")
		}
		exec.push(new StatisticModule(now.format("YYYY-MM-DD"), now.format("DD"), idx - 1, workdays, annualVacations))
	})

	return exec
}
const statistics: StatisticModule[] = parseStatisticTimeHeader(monthlyCollectionOriginalData)

/**
 * 写入原始月度汇总数据 → 时间统计模块
 * @param data 月度汇总元数据
 * @param statistics 时间统计模块
 */
const parseStatisticData = (data: string[][] = [], statistics: StatisticModule[] = []): void => {
	if (!data || !statistics) throw new Error("处理【月度汇总】的数据错误：没有数据或时间表头")

	data.forEach(datum => {
		if (!datum[0]) return
		statistics.forEach(statistic => {
			if (!datum[0] || datum.length < 8) return
			statistic.setOriginalMonthStatistic(datum[0], datum[statistic.idx + 7] ?? "")
		})
	})
}
parseStatisticData(monthlyCollectionOriginalData?.slice(4), statistics)

/**
 * 写入每日打卡数据 → 时间统计模块
 * @param data 打卡时间元数据
 * @param statistics 时间统计模块
 */
const parseClockInData = (data: string[][] = [], statistics: StatisticModule[] = []): void => {
	if (!data || !statistics) throw new Error("处理【每日统计】的数据错误：没有数据或时间表头")

	// 格式化时间格式
	let formattedClockInData: Map<string, Map<string, string[]>> = new Map()

	data.forEach(datum => {
		if (!datum[0] || !datum[7]) return
		if (!formattedClockInData.get(datum[0])) { formattedClockInData.set(datum[0], new Map()) }

		if (!formattedClockInData.get(datum[0])?.get(datum[7])) {
			const [date, week] = datum[7].split(" ")
			formattedClockInData.get(datum[0])?.set(moment(date, "YY-MM-DD").format("YYYY-MM-DD"), [...datum.slice(10, 14), week ?? ""])
		}
	})

	statistics.forEach(statistic => {
		formattedClockInData.forEach((datum, name) => {
			if (!name) return
			if (!statistic.toString()) return
			if (!datum.get(statistic.toString())) return

			statistic.setOriginalClockIn(name, datum.get(statistic.toString()) ?? [])
		})
	})
}
parseClockInData(clockInStatisticOriginalData?.slice(4), statistics)

statistics.forEach(statistic => statistic.analysis())
fs.writeFile(`${excelFilePath}.时间统计模块.json`, JSON.stringify(statistics.map(statistic => statistic.toDict()), null, 4), "utf-8", () => { })

