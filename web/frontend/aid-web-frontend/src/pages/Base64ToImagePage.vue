<template>
  <q-page class="q-pa-md">
    <q-card flat square bordered class="q-mb-md">
      <q-card-section>
        <div class="text-h4 text-deep-orange">Base64 → 图片</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <div class="row q-gutter-md">
          <div class="col q-gutter-md">
            <q-input v-model="baseText" type="textarea" label="输入 base64 文本（可带 data:*;base64, 前缀或仅 base64）" autogrow outlined square class="full-width" autofocus ref="txtBaseText" />
          </div>

          <div class="col flex items-center justify-center">
            <div v-if="imgUrl" class="q-pa-sm">
              <img :src="imgUrl" :alt="altText" style="max-width: 100%; width: 30vw; height: auto; border: 1px solid #e0e0e0; display: block" />
            </div>
          </div>
        </div>

        <div class="row q-gutter-md q-mt-sm">
          <div class="col">
            <q-btn square label="下载图片" color="secondary" @click="downloadImage" class="full-width" :disable="!imgUrl" />
          </div>
          <div class="col">
            <q-btn square label="清除" color="negative" @click="clear" class="full-width" :disable="!imgUrl" />
          </div>
          <div class="col">
            <q-btn square label="在新标签打开" color="blue" @click="openInNewTab" class="full-width" :disable="!imgUrl" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card flat square bordered class="q-mb-md">
      <q-card-section>
        <div class="text-h4 text-deep-orange">上传图片 → Base64</div>
      </q-card-section>
      <q-card-section class="q-pt-none">
        <input ref="fileInput" type="file" accept="image/*" @change="onFileChange" style="display: none" />
        <div class="row q-gutter-sm">
          <div class="col">
            <q-btn square label="选择图片" color="primary" @click="triggerFileUpload" class="full-width" icon="upload" />
          </div>
          <div class="col">
            <q-btn square label="复制 Base64" color="secondary" @click="copyBase64" :disable="!base64Result" class="full-width" />
          </div>
          <div class="col">
            <q-btn square label="填入上方" color="purple" @click="fillToTextarea" :disable="!base64Result" class="full-width" />
          </div>
        </div>
        <div v-if="base64Result" class="q-mt-sm">
          <q-input v-model="base64Result" type="textarea" label="生成的 Base64（Data URI）" readonly autogrow outlined square />
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onBeforeUnmount, watch } from "vue"
import notify from "src/utils/notify"

const baseText = ref("")
const imgUrl = ref("")
const mime = ref("image/png") // fallback mime
const altText = ref("decoded-image")
const base64Result = ref("") // 上传图片生成的 base64
const fileInput = ref(null)
const txtBaseText = ref(null)

// 自动监听 baseText，带防抖
let _parseTimer = null
const PARSE_DELAY = 400 // ms

watch(baseText, (val) => {
  if (_parseTimer) clearTimeout(_parseTimer)
  _parseTimer = setTimeout(() => {
    if (!val) {
      revokeCurrentUrl()
      return
    }

    const parsed = parseBase64Input(val)
    if (!parsed) {
      revokeCurrentUrl()
      return
    }

    try {
      const usedMime = parsed.mime || mime.value || "image/png"
      const blob = base64ToBlob(parsed.base64, usedMime)
      // 先释放旧 URL
      revokeCurrentUrl()
      currentObjectUrl = URL.createObjectURL(blob)
      imgUrl.value = currentObjectUrl
      mime.value = usedMime
    } catch (e) {
      console.error("auto-parse base64 failed", e)
      revokeCurrentUrl()
    }
  }, PARSE_DELAY)
})

onBeforeUnmount(() => {
  revokeCurrentUrl()
  if (_parseTimer) {
    clearTimeout(_parseTimer)
    _parseTimer = null
  }
})

// 将 base64 字符串（可能带 data:*;base64, 前缀）解析为 { mime, base64 }
function parseBase64Input(input) {
  if (!input) return null
  const trimmed = input.trim()
  // data:[<mediatype>][;base64],<data>
  const match = trimmed.match(/^data:([^;]+);base64,(.+)$/i)
  if (match) {
    return { mime: match[1], base64: match[2] }
  }
  // 纯 base64（无前缀）
  // 尝试自动检测是否为 base64 by checking valid chars (rough)
  const pure = trimmed.replace(/\s+/g, "")
  if (/^[A-Za-z0-9+/=]+$/.test(pure)) {
    return { mime: null, base64: pure }
  }
  return null
}

function base64ToBlob(base64, mimeType = "application/octet-stream") {
  const binaryString = atob(base64)
  const len = binaryString.length
  const bytes = new Uint8Array(len)
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }
  return new Blob([bytes], { type: mimeType })
}

let currentObjectUrl = null

function revokeCurrentUrl() {
  if (currentObjectUrl) {
    URL.revokeObjectURL(currentObjectUrl)
    currentObjectUrl = null
    imgUrl.value = ""
  }
}

const downloadImage = () => {
  if (!imgUrl.value) {
    notify.error("当前没有可下载的图片，请先粘贴base64编码到文本框")
    return
  }
  const a = document.createElement("a")
  a.href = imgUrl.value
  // 根据 mime 尝试设置扩展名
  const extMap = {
    "image/png": "png",
    "image/jpeg": "jpg",
    "image/jpg": "jpg",
    "image/gif": "gif",
    "image/svg+xml": "svg",
    "image/webp": "webp",
  }
  const ext = extMap[mime.value] || "bin"
  a.download = `image.${ext}`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const openInNewTab = () => {
  if (!imgUrl.value) return
  window.open(imgUrl.value, "_blank")
}

const clear = () => {
  baseText.value = ""
  revokeCurrentUrl()
  if (txtBaseText.value && typeof txtBaseText.value.focus === "function") {
    txtBaseText.value.focus()
  } else if (txtBaseText.value && txtBaseText.value.$el) {
    // 回退：尝试聚焦原生 textarea 元素
    const ta = txtBaseText.value.$el.querySelector("textarea")
    ta && ta.focus()
  }
}

onBeforeUnmount(() => revokeCurrentUrl())

// 新增：触发文件选择
const triggerFileUpload = () => {
  fileInput.value?.click()
}

// 新增：处理文件上传
const onFileChange = (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  if (!file.type.startsWith("image/")) {
    notify.error("请选择图片文件")
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    base64Result.value = e.target.result // Data URI 格式
    notify.ok(`图片已转换为 Base64（${Math.round(base64Result.value.length / 1024)} KB）`)
    // 清空 file input 以便重复上传同一文件
    if (fileInput.value) fileInput.value.value = ""
  }
  reader.onerror = () => {
    notify.error("读取文件失败")
  }
  reader.readAsDataURL(file)
}

// 新增：复制 Base64 结果
const copyBase64 = () => {
  if (!base64Result.value) {
    notify.error("没有 Base64 内容可复制")
    return
  }
  navigator.clipboard
    ?.writeText(base64Result.value)
    .then(() => {
      notify.ok("已复制到剪贴板")
    })
    .catch(() => {
      // 回退方案
      const ta = document.createElement("textarea")
      ta.value = base64Result.value
      ta.style.position = "fixed"
      ta.style.opacity = "0"
      document.body.appendChild(ta)
      ta.select()
      try {
        document.execCommand("copy")
        notify.ok("已复制到剪贴板")
      } catch {
        notify.error("复制失败")
      } finally {
        document.body.removeChild(ta)
      }
    })
}

// 新增：填入到下方文本框
const fillToTextarea = () => {
  if (!base64Result.value) {
    notify.error("没有 Base64 内容可填入")
    return
  }
  baseText.value = base64Result.value
  notify.ok("已填入到文本框")
}
</script>

<style scoped>
/* 可根据需要微调 */
</style>
