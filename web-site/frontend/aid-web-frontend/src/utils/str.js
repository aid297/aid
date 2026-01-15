/**
 * 字符串处理
 */
export class Str {
	_original = '';

	constructor(str) {
		this._original = str;
	}

	static new = str => new Str(str);

	static make = () => new Str('');

	/**
	 * 寻找字符出现的数量
	 * @param {string} target 寻找目标
	 * @returns {number}
	 */
	findChar(target) {
		const matches = String(this._original).match(new RegExp(target, 'g'));
		return matches ? matches.length : 0;
	}

	/**
	 * 正则替换
	 * @param {string} src 正则
	 * @param {string} dst 替换目标
	 * @returns {string}
	 */
	replace(src = /\s+/g, dst = '') {
		return this._original.replace(src, dst);
	}

	get original() {
		return this._original;
	}

	set original(value) {
		this._original = value;
	}
}
