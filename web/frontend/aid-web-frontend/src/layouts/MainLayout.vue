<template>
  <q-layout view="hHh lpR fFf">
    <q-header bordered class="bg-orange text-white" height-hint="98">
      <q-toolbar>
        <q-btn square dense flat round icon="menu" @click="leftDrawerOpen = !leftDrawerOpen" />

        <q-toolbar-title>
          <q-avatar>
            <img src="/page/favicon.ico" alt="" srcset="" />
            <!-- <img src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg" /> -->
          </q-avatar>
          🤠🤠🤠
        </q-toolbar-title>

        <!-- <q-btn square dense flat round icon="menu" @click="rightDrawerOpen = !rightDrawerOpen" /> -->
      </q-toolbar>

      <q-tabs align="center"> </q-tabs>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" side="left" bordered>
      <q-list link dense separator style="min-width: 200px">
        <q-item dense clickable v-ripple @click="hrefTo(router, '/')" :class="{ 'bg-orange': isActive('/'), 'text-white': isActive('/') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-calendar-week" :style="{ fontSize: '18px' }">&nbsp;考勤预处理</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/timeCalculator')" :class="{ 'bg-orange': isActive('/timeCalculator'), 'text-white': isActive('/timeCalculator') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-clock" :style="{ fontSize: '18px' }">&nbsp;时间计算器</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/rezip')" :class="{ 'bg-orange': isActive('/rezip'), 'text-white': isActive('/rezip') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-briefcase" :style="{ fontSize: '18px' }">&nbsp;重新压缩</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/uuid')" :class="{ 'bg-orange': isActive('/uuid'), 'text-white': isActive('/uuid') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-pen-fancy" :style="{ fontSize: '18px' }">&nbsp;生成 UUID</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/randString')" :class="{ 'bg-orange': isActive('/randString'), 'text-white': isActive('/randString') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-feather" :style="{ fontSize: '18px' }">&nbsp;生成随机字符串</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/codeChange')" :class="{ 'bg-orange': isActive('/codeChange'), 'text-white': isActive('/codeChange') }">
          <q-item-section avatar class="q-mt-md q-mb-md" :active="isActive('/codeChange')">
            <i class="fa-solid fa-code-compare" :style="{ fontSize: '18px' }">&nbsp;代码格式转换</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/base64ToImage')" :class="{ 'bg-orange': isActive('/base64ToImage'), 'text-white': isActive('/base64ToImage') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-image" :style="{ fontSize: '18px' }">&nbsp;Base64转换图片</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/encoder')" :class="{ 'bg-orange': isActive('/encoder'), 'text-white': isActive('/encoder') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-shuffle" :style="{ fontSize: '18px' }">&nbsp;编码转换</i>
          </q-item-section>
        </q-item>

        <q-item dense clickable v-ripple @click="hrefTo(router, '/qrCode')" :class="{ 'bg-orange': isActive('/qrCode'), 'text-white': isActive('/qrCode') }">
          <q-item-section avatar class="q-mt-md q-mb-md">
            <i class="fa-solid fa-qrcode" :style="{ fontSize: '18px' }">&nbsp;QRCode</i>
          </q-item-section>
        </q-item>
      </q-list>
      <!-- drawer content -->
    </q-drawer>

    <!-- <q-drawer v-model="rightDrawerOpen" side="right" overlay bordered> -->
    <!-- drawer content -->
    <!-- </q-drawer> -->

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-footer bordered class="bg-grey-8 text-white">
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="/page/favicon.ico" />
          </q-avatar>
          <div></div>
        </q-toolbar-title>
      </q-toolbar>
    </q-footer>
  </q-layout>
</template>

<script setup>
import { ref } from "vue"
import { useRouter, useRoute } from "vue-router"

const router = useRouter()
const route = useRoute()
const leftDrawerOpen = ref(true)

const isActive = (path) => route.path === path
const hrefTo = (router, url) => router.push(url)

// const rightDrawerOpen = ref(false)
</script>

<style scoped></style>
