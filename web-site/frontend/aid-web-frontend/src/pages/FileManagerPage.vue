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
                                <q-uploader :url="uploadUrl" label="上传文件" @uploaded="handleUploaded" class="max-wight"
                                    @failed="handleFailed" flat bordered field-name="file" style="width: 100%" />
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
                        <q-table flat bordered separator title="文件管理" :rows="rows" dark color="amber"
                            :pagination="{ rowsPerPage: 0 }">
                            <template v-slot:header="props">
                                <q-tr :props="props">
                                    <q-th align="left" key="name" name="name">名称</q-th>
                                    <q-th align="left" key="size" name="size">大小</q-th>
                                    <q-th align="left" key="kind" name="kind">类型</q-th>
                                    <q-th align="left" key="option" name="option">操作</q-th>
                                </q-tr>
                            </template>
                            <template v-slot:body="props">
                                <q-tr :props="props">
                                    <q-td align="left" key="name" :props="props">{{ props.row.name || '' }}</q-td>
                                    <q-td align="left" key="size" :props="props">{{ props.row.size || '' }}</q-td>
                                    <q-td align="left" key="kind" :props="props">{{ props.row.kind || '' }}</q-td>
                                    <q-td align="left" key="fullPath" :props="props">
                                        <q-btn-group flat>
                                            <q-btn size="sm" color="primary" icon="download" label="下载"
                                                @click="handleDownload(props.row)" :disable="props.row.kind === 'dir'" />
                                            <q-btn size="sm" color="negative" icon="delete" label="删除"
                                                @click="handleDelete(props.row)" />
                                            <q-btn size="sm" color="info" icon="edit" label="重命名"
                                                @click="handleRename(props.row)" />
                                        </q-btn-group>
                                    </q-td>
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
import { API_BASE_URL, axios } from 'src/utils/fetch';
import notify from 'src/utils/notify';
import { onMounted, ref } from 'vue';

const rows = ref([]);
const currentDir = ref('/');

/**
 * 加载文件列表
 * @param dir 所需目录
 */
const loadFileList = async () => {
    const { dirs, files } = (await axios.post('/fileManager/list', { dir: currentDir.value })).data.data;
    rows.value = [...dirs, ...files];
    console.log('文件列表已加载', rows.value);
};

onMounted(loadFileList);

const uploadUrl = ref(`${API_BASE_URL}/fileManager/upload`); // 后端上传接口
const handleUploaded = info => {
    console.log('文件上传成功', info);
    notify.ok('上传成功');
    loadFileList(); // 重新加载文件列表
};

const handleFailed = () => {
    notify.error('上传失败');
};

const handleDownload = row => {
    if (row.kind === 'dir') return;
    window.open(`${API_BASE_URL}/fileManager/download?path=${encodeURIComponent(row.fullPath)}`, '_blank');
};

const handleDelete = async row => {
    try {
        await axios.post('/fileManager/delete', { path: row.fullPath });
        notify.ok('删除成功');
        loadFileList();
    } catch (error) {
        notify.error('删除失败', error);
    }
};

const handleRename = row => {
    // TODO: 实现重命名功能
    notify.info(`重命名功能待实现: ${row.name}`);
};
</script>
