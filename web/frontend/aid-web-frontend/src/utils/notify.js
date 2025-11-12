import { Notify } from "quasar"

const defaultOptions = { position: "center", timeout: 1500 }

export default {
  ok: (msg = "") => {
    Notify.create({ type: "positive", message: msg, ...defaultOptions })
  },
  error: (msg = "") => {
    Notify.create({ type: "negative", message: msg, ...defaultOptions })
  },
}
