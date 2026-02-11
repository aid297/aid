import { Notify } from 'quasar';

const defaultOptions = { position: 'center' };
const notAskOptions = { timeout: 1500, ...defaultOptions };

export default {
    ok: (msg = '') => {
        Notify.create({ type: 'positive', message: msg, ...notAskOptions });
    },
    error: (msg = '') => {
        Notify.create({ type: 'negative', message: msg, ...notAskOptions });
    },
    ask: (msg = '', onOK, onCancel) => {
        Notify.create({
            type: 'warning',
            message: msg,
            ...defaultOptions,
            actions: [
                { label: '取消', color: 'positive', handler: onCancel },
                { label: '确定', color: 'negative', handler: onOK },
            ]
        });
    }
};
