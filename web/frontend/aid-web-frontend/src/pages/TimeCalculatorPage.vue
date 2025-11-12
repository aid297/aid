<template>
  <div class="q-pa-md q-gutter-md">
    <q-card flat square bordered class="q-mb-md">
      <q-card-section>
        <div class="text-h4 text-deep-orange">时间计算器</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <div class="row q-gutter-md">
          <div class="col-4"></div>
          <div class="col-4">
            <h3>相差：{{ diffDays }}天</h3>
          </div>
          <div class="col-4"></div>
        </div>

        <div class="row q-gutter-md">
          <div class="col">
            <q-card flat square bordered class="my-card">
              <q-card-section>
                <q-list>
                  <q-item clickable>
                    <q-item-section>
                      <span style="font-size: 20px" class="text-primary">1. 起始日期</span>
                      <span class="text-secondary">&nbsp;</span>
                    </q-item-section>

                    <q-item-section avatar>
                      <q-btn square @click="clearOriginalDates" round color="red" icon="fa fa-times"></q-btn>
                    </q-item-section>

                    <q-item-section avatar>
                      <q-btn square icon="event" round color="primary">
                        <q-popup-proxy @before-show="originalDates_updateProxy" cover transition-show="scale"
                          transition-hide="scale">
                          <q-date v-model="originalDatesProxy" flat bordered today-btn landscape years-in-month-view
                            mark="YYYY-MM-DD">
                            <div class="row items-center justify-end q-gutter-sm">
                              <q-btn square label="Cancel" color="primary" flat v-close-popup />
                              <q-btn square label="OK" color="primary" flat @click="originalDates_save" v-close-popup />
                            </div>
                          </q-date>
                        </q-popup-proxy>
                      </q-btn>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-card-section>

              <q-separator inset />

              <q-card-section>{{ originalDate }}</q-card-section>
            </q-card>
          </div>
          <div class="col">
            <q-card flat square bordered class="my-card">
              <q-card-section>
                <q-list>
                  <q-item clickable>
                    <q-item-section>
                      <span style="font-size: 20px" class="text-primary">2. 终止日期</span>
                      <span class="text-secondary">&nbsp;</span>
                    </q-item-section>

                    <q-item-section avatar>
                      <q-btn square @click="clearFinishedDates" round color="red" icon="fa fa-times"></q-btn>
                    </q-item-section>

                    <q-item-section avatar>
                      <q-btn square icon="event" round color="primary">
                        <q-popup-proxy @before-show="finishedDates_updateProxy" cover transition-show="scale"
                          transition-hide="scale">
                          <q-date v-model="finishedDatesProxy" flat bordered today-btn landscape years-in-month-view
                            mark="YYYY-MM-DD">
                            <div class="row items-center justify-end q-gutter-sm">
                              <q-btn square label="Cancel" color="primary" flat v-close-popup />
                              <q-btn square label="OK" color="primary" flat @click="finishedDates_save" v-close-popup />
                            </div>
                          </q-date>
                        </q-popup-proxy>
                      </q-btn>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-card-section>

              <q-separator inset />

              <q-card-section>{{ finishedDate }}</q-card-section>
            </q-card>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref, watch } from "vue"

import moment from "moment"
import notify from "src/utils/notify"

// 起始日期
const originalDate = ref("")
const originalDatesProxy = ref("")
const originalDates_updateProxy = () => (originalDatesProxy.value = originalDate.value)
const originalDates_save = () => (originalDate.value = originalDatesProxy.value)
const clearOriginalDates = () => (originalDate.value = "")
originalDatesProxy.value = ""

// 终止日期日期
const finishedDate = ref("")
const finishedDatesProxy = ref("")
const finishedDates_updateProxy = () => (finishedDatesProxy.value = finishedDate.value)
const finishedDates_save = () => (finishedDate.value = finishedDatesProxy.value)
const clearFinishedDates = () => (finishedDate.value = "")
finishedDatesProxy.value = ""

const diffDays = ref(0)

watch(originalDate, () => calculator(originalDate.value, finishedDate.value))
watch(finishedDate, () => calculator(originalDate.value, finishedDate.value))

const calculator = function (originalDate, finishedDate) {
  if (originalDate && finishedDate) {
    diffDays.value = moment(finishedDate, "YYYY-MM-DD").diff(moment(originalDate, "YYYY-MM-DD"), "days")

    notify.ok(`相差：${diffDays.value}天`)
  }
}
</script>
