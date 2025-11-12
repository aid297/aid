import { defineStore, acceptHMRUpdate } from "pinia"

export const useFinalStatisticStore = defineStore("finalStatistic", {
  state: () => ({ data: [] }),
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useFinalStatisticStore, import.meta.hot))
}
