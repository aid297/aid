import { DateTitleService } from "src/services/dateTitleService";
import { Collect, CollectType } from "src/modules/collectModule";
import { Str } from "src/utils/str";

/**
 * 服务：汇总数据
 */
export class CollectService {
    _original = [];
    _dateTitles = DateTitleService.make();
    _collectData = [];

    constructor(original = null, dateTitles = null) {
        if (original !== undefined) this._original = original;
        if (dateTitles !== undefined) this._dateTitles = dateTitles;
    }

    static new = (original, dateTitles) => new CollectService(original, dateTitles);

    static make = () => new CollectService();

    get original() {
        return this._original;
    }
    get dateTitles() {
        return this._dateTitles;
    }
    get data() {
        return this._collectData;
    }

    set original(value) {
        this._original = value;
    }
    set dateTitles(value) {
        this._dateTitles = value;
    }
    set data(value) {
        this._collectData = value;
    }

    parse() {
        this._original.forEach(row => {
            let tmp = [];
            const name = Str.new(row[1]).replace();
            tmp.push(Collect.new(row[0], CollectType.rowNumber));
            tmp.push(Collect.new(name, CollectType.name));
            this._dateTitles.data.forEach(dateTitle => tmp.push(Collect.new(row[dateTitle.idx], CollectType.collect, dateTitle)));
            this._collectData.push(tmp);
        });

        return this;
    }
}
