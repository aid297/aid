<template>
  <div class="row">
    <div class="col">
      <q-page class="q-pa-md">
        <div class="row q-gutter-md">
          <div class="col">
            <q-card flat square bordered>
              <q-card-section>
                <div class="text-h4 text-deep-orange">重新压缩</div>
                <div class="text-h6 text-brown">解决mac压缩包在Windows中解压缩中文无法正确显示的问题</div>
              </q-card-section>

              <q-card-section class="q-pt-none">
                <q-uploader url="http://aid.yujizhou.com:9900/api/v1/rezip/upload" color="purple" label="选择压缩包"
                  field-name="f" flat bordered style="width: 100%" @uploaded="onUploaded" />
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-page>
    </div>
  </div>
</template>

<script setup>
import notify from "src/utils/notify"
const onUploaded = function (file) {
  const res = JSON.parse(file.xhr.response)

  notify.ok(res.msg)
  if (res.code === 200) location.href = `http://${res.data.to}`
}
</script>
