<template>
  <q-page class="q-pa-md">
    <div class="row q-gutter-md">
      <!-- 生成 QRCode -->
      <div class="col">
        <q-card flat square bordered class="my-card">
          <q-card-section>
            <div class="text-h4 text-deep-orange">生成 QRCode</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <div class="row q-gutter-md">
              <div class="col">
                <q-input v-model="qrCodeText" type="textarea" autogrow outlined square label="文本 / URL（要生成二维码的内容）" placeholder="在此输入文本或 URL" ref="txtQRCodeText" autofocus />
              </div>

              <div class="col">
                <div class="row q-gutter-sm">
                  <div class="col">
                    <q-input v-model.number="size" type="number" outlined square label="像素宽度" :min="64" :max="2048" />
                  </div>
                </div>
                <div class="row q-gutter-sm q-mt-sm">
                  <div class="col">
                    <q-select v-model="ecLevel" :options="ecOptions" outlined square label="容错等级" />
                  </div>
                </div>
                <div class="row q-gutter-sm q-mt-sm">
                  <div class="col">
                    <q-input v-model.number="margin" type="number" outlined square label="留白" :min="1" :max="20" />
                  </div>
                </div>
                <div class="row q-gutter-sm q-mt-sm">
                  <div class="col">
                    <q-btn square label="生成" color="primary" @click="generateQRCode" class="full-width" />
                  </div>
                  <div class="col">
                    <q-btn square label="清除" color="negative" @click="clearAll" class="full-width" />
                  </div>
                </div>
              </div>
            </div>

            <div class="row q-gutter-md q-mt-md">
              <div class="col">
                <q-card flat square bordered>
                  <q-card-section>
                    <div class="text-subtitle2 q-mb-sm">预览</div>
                    <div v-if="imgUrl" class="flex flex-center">
                      <img :src="imgUrl" alt="qrcode" style="max-width: 100%; height: auto; display: block" />
                    </div>
                    <div v-else class="text-caption">尚未生成二维码</div>
                  </q-card-section>
                  <q-separator />
                  <q-card-actions align="right">
                    <q-btn square flat label="下载 PNG" @click="downloadQRCode" :disable="!imgUrl" />
                    <q-btn square flat label="复制图片" @click="copyImageToClipboard" :disable="!imgUrl" />
                    <q-btn square flat label="复制 DataURL" @click="copyDataURL" :disable="!imgUrl" />
                  </q-card-actions>
                </q-card>
              </div>

              <div class="col">
                <q-card flat square bordered>
                  <q-card-section>
                    <div class="text-subtitle2 q-mb-sm">高级（DataURL / Canvas）</div>
                    <q-input v-model="imgUrl" type="textarea" outlined square label="Data URL" autogrow readonly />
                    <div class="q-mt-sm">
                      <canvas ref="canvasRef" :width="size" :height="size" style="display: none"></canvas>
                    </div>
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- QRCode 解析 -->
      <div class="col">
        <q-card flat square bordered class="my-card">
          <q-card-section>
            <div class="text-h6">解析 QRCode（从图片）</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <div class="row q-gutter-md">
              <div class="col">
                <input ref="fileInput" type="file" accept="image/*" @change="onParseFileChange" style="display: none" />
                <div class="row q-gutter-sm">
                  <div class="col">
                    <q-btn square label="选择图片解析" color="primary" @click="triggerParseFile" class="full-width" />
                  </div>
                  <div class="col">
                    <q-btn square label="从左侧 DataURL 解析" color="secondary" @click="parseFromDataURL" class="full-width" :disable="!imgUrl" />
                  </div>
                </div>

                <div class="q-mt-md">
                  <q-input v-model="parseSource" outlined square label="图片 URL 或 DataURL（可粘贴）" placeholder="可输入图片 URL 或 data:*;base64,..." ref="txtParseSource" />
                  <div class="row q-gutter-sm q-mt-sm">
                    <div class="col">
                      <q-btn square label="解析 URL/DataURL" color="primary" @click="parseFromURLinInput" class="full-width" />
                    </div>
                    <div class="col">
                      <q-btn square label="清除解析结果" color="negative" @click="clearParse" class="full-width" />
                    </div>
                  </div>
                </div>
              </div>

              <div class="col">
                <q-card flat square bordered>
                  <q-card-section>
                    <div class="text-subtitle2 q-mb-sm">解析结果</div>
                    <q-input v-model="parsedText" type="textarea" outlined square label="Decoded Text" readonly autogrow />
                  </q-card-section>
                  <q-separator />
                  <q-card-actions align="right">
                    <q-btn square flat label="复制结果" @click="copyParsed" :disable="!parsedText" />
                  </q-card-actions>
                </q-card>
                <canvas ref="parseCanvasRef" style="display: none"></canvas>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, watch } from "vue"
import QRCode from "qrcode"
import jsQR from "jsqr"
import notify from "src/utils/notify"

const qrCodeText = ref("")
const size = ref(256)
const margin = ref(1)
const ecLevel = ref("M")
const imgUrl = ref("")
const canvasRef = ref(null)
const txtQRCodeText = ref(null)
const txtParseSource = ref(null)

watch(margin, (newVal) => {
  if (newVal < 1) margin.value = 1
  if (newVal > 20) margin.value = 20
})

const ecOptions = [
  { label: "L (7%)", value: "L" },
  { label: "M (15%)", value: "M" },
  { label: "Q (25%)", value: "Q" },
  { label: "H (30%)", value: "H" },
]

/* 生成相关 */
const generateQRCode = async () => {
  if (!qrCodeText.value || !qrCodeText.value.trim()) {
    notify.error("请输入要生成的内容")
    return
  }
  try {
    const opts = {
      errorCorrectionLevel: ecLevel.value,
      width: size.value,
      margin: margin.value,
      color: {
        dark: "#000000",
        light: "#FFFFFF",
      },
    }
    // 生成 Data URL 并同时渲染到 canvas（可作为备用）
    imgUrl.value = await QRCode.toDataURL(qrCodeText.value, opts)
    try {
      const canvas = canvasRef.value
      if (canvas) await QRCode.toCanvas(canvas, qrCodeText.value, opts)
    } catch (e) {
      // canvas 渲染失败不影响 DataURL
      console.warn("Canvas render failed", e)
    }
    notify.ok("二维码生成成功")
  } catch (e) {
    console.error(e)
    notify.error("二维码生成失败：" + (e.message || e))
  }
}

const downloadQRCode = () => {
  if (!imgUrl.value) {
    notify.error("尚未生成二维码")
    return
  }
  const a = document.createElement("a")
  a.href = imgUrl.value
  a.download = "qrcode.png"
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const copyDataURL = async () => {
  if (!imgUrl.value) {
    notify.error("尚未生成二维码")
    return
  }

  const text = imgUrl.value

  // 优先使用现代 clipboard API（有时 navigator.clipboard 为 undefined）
  if (typeof navigator !== "undefined" && navigator.clipboard && typeof navigator.clipboard.writeText === "function") {
    try {
      await navigator.clipboard.writeText(text)
      notify.ok("已复制解析结果")
      return
    } catch (e) {
      console.warn("navigator.clipboard.writeText failed:", e)
    }
  }

  // 回退：使用不可见 textarea + document.execCommand('copy')
  try {
    const ta = document.createElement("textarea")
    ta.value = text
    ta.setAttribute("readonly", "")
    ta.style.position = "fixed"
    ta.style.left = "-9999px"
    document.body.appendChild(ta)
    ta.select()

    const ok = document.execCommand("copy")
    document.body.removeChild(ta)

    if (ok) {
      notify.ok("已复制解析结果")
    } else {
      notify.error("复制失败")
    }
  } catch (e) {
    console.error("copy fallback failed", e)
    notify.error("复制失败")
  }
}

const copyImageToClipboard = async () => {
  if (!imgUrl.value) {
    notify.error("尚未生成二维码")
    return
  }
  try {
    // fetch dataURL -> blob
    const res = await fetch(imgUrl.value)
    const blob = await res.blob()
    if (navigator.clipboard && window.ClipboardItem) {
      await navigator.clipboard.write([new ClipboardItem({ [blob.type]: blob })])
      notify.ok("图片已复制到剪贴板")
    } else {
      // fallback: 复制 dataURL 文本
      await navigator.clipboard.writeText(imgUrl.value)
      notify.ok("浏览器不支持直接复制图片，已复制 DataURL 文本")
    }
  } catch (e) {
    console.error(e)
    notify.error("复制图片失败：" + (e.message || e))
  }
}

const clearAll = () => {
  qrCodeText.value = ""
  imgUrl.value = ""
  const canvas = canvasRef.value
  if (canvas) {
    const ctx = canvas.getContext("2d")
    ctx && ctx.clearRect(0, 0, canvas.width, canvas.height)
  }
  txtQRCodeText.value?.focus()
}

/* 解析相关 */
const fileInput = ref(null)
const parseCanvasRef = ref(null)
const parseSource = ref("") // URL or DataURL input
const parsedText = ref("")

const triggerParseFile = () => {
  fileInput.value?.click()
}

const onParseFileChange = (e) => {
  const f = e.target.files?.[0]
  if (!f) return
  if (!f.type.startsWith("image/")) {
    notify.error("请选择图片文件")
    return
  }
  const reader = new FileReader()
  reader.onload = () => {
    parseImageFromDataURL(reader.result)
  }
  reader.onerror = () => notify.error("读取文件失败")
  reader.readAsDataURL(f)
  // reset input
  if (fileInput.value) fileInput.value.value = ""
}

const parseFromDataURL = async () => {
  if (!imgUrl.value) {
    notify.error("上方没有 DataURL")
    return
  }
  parseImageFromDataURL(imgUrl.value)
}

const parseFromURLinInput = async () => {
  if (!parseSource.value || !parseSource.value.trim()) {
    notify.error("请输入图片 URL 或 DataURL")
    return
  }
  // 如果是普通 URL，尝试加载
  const s = parseSource.value.trim()
  if (s.startsWith("data:")) {
    parseImageFromDataURL(s)
  } else {
    // load remote image (may be CORS blocked)
    try {
      const res = await fetch(s)
      if (!res.ok) throw new Error("网络错误：" + res.status)
      const blob = await res.blob()
      const reader = new FileReader()
      reader.onload = () => parseImageFromDataURL(reader.result)
      reader.readAsDataURL(blob)
    } catch (e) {
      notify.error("加载图片失败：" + (e.message || e))
    }
  }
}

function parseImageFromDataURL(dataUrl) {
  parsedText.value = ""
  const img = new Image()
  img.crossOrigin = "anonymous"
  img.onload = () => {
    try {
      const canvas = parseCanvasRef.value
      canvas.width = img.width
      canvas.height = img.height
      const ctx = canvas.getContext("2d")
      ctx.drawImage(img, 0, 0)
      const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
      const code = jsQR(imageData.data, canvas.width, canvas.height)
      if (code && code.data) {
        parsedText.value = code.data
        notify.ok("解析成功")
      } else {
        notify.error("未识别到 QRCode")
      }
    } catch (e) {
      console.error(e)
      notify.error("解析失败：" + (e.message || e))
    }
  }
  img.onerror = (e) => {
    notify.error(`图片加载失败：${e.message || e}`)
  }
  img.src = dataUrl
}

const copyParsed = async () => {
  if (!parsedText.value) return notify.error("没有解析结果可复制")
  try {
    await navigator.clipboard.writeText(parsedText.value)
    notify.ok("已复制解析结果")
  } catch {
    const ta = document.createElement("textarea")
    ta.value = parsedText.value
    ta.style.position = "fixed"
    ta.style.opacity = "0"
    document.body.appendChild(ta)
    ta.select()
    try {
      document.execCommand("copy")
      notify.ok("已复制解析结果")
    } catch {
      notify.error("复制失败")
    } finally {
      document.body.removeChild(ta)
    }
  }
}

const clearParse = () => {
  parseSource.value = ""
  parsedText.value = ""
  const canvas = parseCanvasRef.value
  if (canvas) {
    const ctx = canvas.getContext("2d")
    ctx && ctx.clearRect(0, 0, canvas.width, canvas.height)
  }
  txtParseSource.value?.focus()
}
</script>

<style scoped></style>
