<template>
    <div class="row">
        <div class="col">
            <q-page class="q-pa-md">
                <div class="row q-gutter-md">
                    <div class="col">
                        <q-card flat square bordered>
                            <q-card-section>
                                <div class="text-h4 text-deep-orange">文件管理</div>
                            </q-card-section>

                            <q-card-section class="q-pt-none">
                                <q-uploader :url="uploadUrl" label="上传文件" @uploaded="handleUploaded" class="max-wight" @failed="handleFailed" flat bordered field-name="file" style="width: 100%" />
                            </q-card-section>
                        </q-card>
                    </div>
                </div>
            </q-page>
        </div>
    </div>

    <div class="row">
        <div class="col">
            <q-page class="q-pa-md">
                <div class="row q-gutter-md">
                    <div class="col">
                        <q-table flat bordered separator title="文件管理" :rows="rows" dark color="amber" :pagination="{ rowsPerPage: 0 }">
                            <template v-slot:header="props">
                                <q-tr :props="props">
                                    <q-th align="left" key="uuid" name="uuid">uuid</q-th>
                                </q-tr>
                            </template>
                            <template v-slot:body="props">
                                <q-tr :props="props">
                                    <q-td key="uuid" :props="props">{{ props.row.uuid || "??" }}</q-td>
                                </q-tr>
                            </template>
                        </q-table>
                    </div>
                </div>
            </q-page>
        </div>
    </div>
</template>

<script setup>
import { API_BASE_URL } from "src/utils/fetch";
import notify from "src/utils/notify";
import { ref } from "vue";

const uploadUrl = ref(`${API_BASE_URL}/fileManager/upload`); // 后端上传接口
// const headers = ref({
//     // 如果需要认证，可以在这里设置请求头，例如Token
//     // 'Authorization': `Bearer ${yourToken}`
// })

const handleUploaded = info => {
    console.log("文件上传成功", info);
    // 可以在这里处理上传成功后的逻辑，如提示用户、更新文件列表等
    // this.$q.notify({ type: 'positive', message: '上传成功！' })
    notify.ok("上传成功");
};
</script>
