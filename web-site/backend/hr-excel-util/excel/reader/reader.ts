import * as XLSX from 'xlsx'

export class Reader {
	// _originalBuffer: Promise<Buffer>
	_workbook?: XLSX.WorkBook
	_sheetName?: string
	_err?: Error

	constructor(file: Buffer) {
		this._workbook = XLSX.read(file, { type: 'buffer' })
	}

	static new = (file: Buffer) => new Reader(file)

	error(): Error | undefined { return this._err }

	setSheetName(sheetName: string): Reader {
		if (!sheetName) {
			this._err = new Error('sheetName不能为空')
			return this
		}

		if (!this._workbook) {
			this._err = new Error('请先调用load方法加载工作簿')
			return this
		}

		if (!this._workbook.SheetNames.includes(sheetName)) {
			this._err = new Error(`工作簿中没有名为${sheetName}的工作表`)
			return this
		}

		return this
	}

	loadSheet(sheetName: string = ""): string[][] {
		try {
			if (!this._sheetName && !sheetName) {
				this._err = new Error('需要先调用setSheetName设置sheet名称')
				return []
			}

			if (this._workbook === null || this._workbook === undefined) {
				this._err = new Error('请先调用load方法加载工作簿')
				return []
			}

			const sheet = this._workbook.Sheets[this._sheetName || sheetName]
			if (!sheet) {
				this._err = new Error(`工作簿中没有名为${this._sheetName}的工作表`)
				return []
			}

			return XLSX.utils.sheet_to_json(sheet, { header: 1 })
		} catch (error) {
			this._err = error as Error
			return []
		}
	}
}