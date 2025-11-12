<template>
  <div class="q-pa-md q-gutter-md">
    <q-card flat square bordered class="q-mb-md">
      <q-card-section>
        <div class="text-h4 text-deep-orange">随机字符串</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <div class="row q-gutter-md">
          <div class="col-4"></div>
          <div class="col-4">
            <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
              <div>
                <q-checkbox v-model="digit" checked-icon="star" unchecked-icon="star_border" indeterminate-icon="help" />
                数字：【0-9】
              </div>
              <div>
                <q-checkbox v-model="lowercase" checked-icon="star" unchecked-icon="star_border" indeterminate-icon="help" />
                小写字母：【a-z】
              </div>
              <div>
                <q-checkbox v-model="uppercase" checked-icon="star" unchecked-icon="star_border" indeterminate-icon="help" />
                大写字母：【A-Z】
              </div>
              <div>
                <q-checkbox v-model="specialCharacter" checked-icon="star" unchecked-icon="star_border" indeterminate-icon="help" />
                特殊字符：【~!@#$%^&*()[]{}-_=+|;:'",./?`】
              </div>
              <div class="q-gutter-sm">
                <q-radio v-model="shape" val="normal" label="普通" />
                <q-radio v-model="shape" val="web-crypto" label="WebCrypto" />
                <q-radio v-model="shape" val="base64url" label="Base64URL" />
              </div>
              <div>
                <q-input v-model="number" label="数量" hint="必填" lazy-rules type="number" max="100" min="1" autofocus outlined square class="full-width" step="1" ref="txtNumber" />
              </div>
              <div>
                <q-btn square label="提交" type="submit" color="primary" class="full-width" />
              </div>
            </q-form>
          </div>
          <div class="col-4"></div>
        </div>
        <div class="row q-gutter-md q-mt-md">
          <div class="col-4"></div>
          <div class="col-4">
            <span @click="onCopy">结果：{{ result }}</span>
          </div>
          <div class="col-4"></div>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref, watch } from "vue"
import { copyToClipboard } from "quasar"
import notify from "src/utils/notify"

const number = ref(32)
const digit = ref(true)
const lowercase = ref(true)
const uppercase = ref(true)
const specialCharacter = ref(false)
const shape = ref("normal") // normal | web-crypto | base64url
const result = ref("")
const txtNumber = ref(null)

const digitCharts = "1234567890"
const lowercaseCharts = "abcdefghijklmnopqrstuvwxyz"
const uppercaseCharts = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const specialCharacterCharts = `~!@#$%^&*()[]{}-_=+|;:'",./?\\\``

watch(digit, () => {
  onSubmit()
})
watch(lowercase, () => {
  onSubmit()
})
watch(uppercase, () => {
  onSubmit()
})
watch(specialCharacter, () => {
  onSubmit()
})
watch(shape, () => {
  onSubmit()
})
watch(number, () => {
  onSubmit()
})

const getChars = () => {
  let chars = ""
  if (digit.value) chars += digitCharts
  if (lowercase.value) chars += lowercaseCharts
  if (uppercase.value) chars += uppercaseCharts
  if (specialCharacter.value) chars += specialCharacterCharts
  if (!chars) chars = lowercaseCharts + uppercaseCharts + digitCharts
  return chars
}

// 1) 简单版（非加密安全）
const randomString = (len, chars) => {
  let s = ""
  for (let i = 0; i < len; i++) s += chars[Math.floor(Math.random() * chars.length)]
  return s
}

// 2) 加密安全（浏览器：Web Crypto）
const randomStringSecure = (len = 16, chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_") => {
  const charLen = chars.length
  if (charLen > 255) throw new Error("chars length must be <= 255")
  const out = []
  const buf = new Uint8Array(len * 2)
  let i = 0
  while (i < len) {
    crypto.getRandomValues(buf)
    for (let b of buf) {
      if (b >= Math.floor(256 / charLen) * charLen) continue // 拒绝采样避免偏差
      out.push(chars[b % charLen])
      if (++i === len) break
    }
  }
  return out.join("")
}

// 3) 加密安全：随机字节转 Base64URL
const randomBase64Url = (bytes = 16) => {
  const arr = new Uint8Array(bytes)
  crypto.getRandomValues(arr)
  const b64 = btoa(String.fromCharCode(...arr))
  return b64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/g, "")
}

const onSubmit = () => {
  const funcs = {
    normal: randomString,
    "web-crypto": randomStringSecure,
    base64url: randomBase64Url,
  }

  result.value = funcs[shape.value] ? funcs[shape.value](number.value, getChars()) : ""

  txtNumber.value?.focus()
}

const onCopy = async () => {
  await copyToClipboard(result.value)
  notify.ok("已复制到剪贴板")
}
</script>
