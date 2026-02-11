<template>
    <div class="row">
        <div class="col">
            <div class="row q-gutter-md">
                <div class="col">
                    <q-card flat square bordered>
                        <q-card-section>
                            <div class="text-h4 text-deep-orange">文件管理</div>
                        </q-card-section>

                        <q-card-section class="q-pt-none">
                            <div class="row">
                                <div class="col">
                                    <q-input outlined bottom-slots v-model="newFolderName" label="新建目录名称" counter
                                        maxlength="64" :dense="dense" @keydown.enter="handleStoreFolder"
                                        ref="inputStoreFolder">
                                        <template v-slot:append>
                                            <q-icon v-if="newFolderName !== ''" name="close" @click="newFolderName = ''"
                                                class="cursor-pointer" />
                                        </template>

                                        <template v-slot:after>
                                            <q-btn round dense flat icon="send" @click="handleStoreFolder" />
                                        </template>
                                    </q-input>
                                </div>
                                <div class="col">
                                    <q-uploader :url="uploadUrl" label="上传文件" @uploaded="handleUploaded"
                                        class="max-wight" @failed="handleFailed" flat bordered field-name="file"
                                        style="width: 100%" />
                                </div>
                            </div>
                        </q-card-section>
                    </q-card>
                </div>
            </div>
        </div>
    </div>

    <div class="row">
        <div class="col">
            <div class="row q-gutter-md">
                <div class="col">
                    <q-card flat square bordered>
                        <q-card-section>
                            <q-table flat bordered separator title="文件管理" :rows="rows" dark color="amber"
                                :pagination="{ rowsPerPage: 0 }">
                                <template v-slot:header="props">
                                    <q-tr :props="props">
                                        <q-th align="left" key="name" name="name">名称</q-th>
                                        <q-th align="left" key="option" name="option">操作</q-th>
                                    </q-tr>
                                </template>
                                <template v-slot:body="props">
                                    <q-tr :props="props">
                                        <q-td align="left" key="name" :props="props">
                                            <a class="text-white text-decoration-none" v-if="props.row.kind === 'DIR'"
                                                href="#" @click.prevent="loadFileList(props.row.name);">
                                                <i class="fa fa-folder">&nbsp;</i>{{ props.row.name || '' }}
                                            </a>
                                            <span v-else>
                                                <a class="text-white text-decoration-none" href="#"
                                                    @click.prevent="handleDownload(props.row)">
                                                    <i class="fa-regular fa-file">&nbsp;</i>{{ props.row.name || '' }}
                                                </a>
                                            </span>
                                        </q-td>
                                        <q-td align="left" key="path" :props="props">
                                            <q-btn-group flat>
                                                <q-btn size="sm" color="negative" @click="handleDestroy(props.row)"
                                                    v-if="props.row.name !== '..'">
                                                    <i class="fa fa-trash">&nbsp;</i>删除
                                                </q-btn>
                                                <q-btn size="sm" color="info" @click="handleZip(props.row)"
                                                    v-if="props.row.name !== '..'">
                                                    <i class="fa fa-box-archive">&nbsp;</i>压缩
                                                </q-btn>
                                            </q-btn-group>
                                        </q-td>
                                    </q-tr>
                                </template>
                            </q-table>
                        </q-card-section>
                    </q-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { API_BASE_URL, axios } from 'src/utils/fetch';
import notify from 'src/utils/notify';
import { computed, onMounted, ref } from 'vue';

const newFolderName = ref('');
const rows = ref([]);
const currentDir = ref('');
const inputStoreFolder = ref(null);

/**
 * 加载文件列表
 * @param dir 所需目录
 */
const loadFileList = async (name = '') => {
    const { currentPath:newCurrentDir, filesystemers } = (await axios.post('/fileManager/list', { body: { path: currentDir.value, name } })).data.content;
    currentDir.value = newCurrentDir;
    rows.value = [{ path: newCurrentDir, name: '..', kind: 'DIR' }, ...filesystemers]; // 在文件列表前添加返回上级目录的项;
    newFolderName.value = '';
    inputStoreFolder.value.focus();
};

onMounted(loadFileList);

// 根据当前目录动态生成上传URL
const uploadUrl = computed(() => `${API_BASE_URL}/fileManager/upload?path=${encodeURIComponent(currentDir.value)}`);

const handleUploaded = async info => {
    console.log('文件上传成功', info);
    await loadFileList(); // 重新加载文件列表
};

const handleStoreFolder = async () => {
    if (newFolderName.value.trim() !== '') {
        await axios.post('/fileManager/storeFolder', { body: { path: currentDir.value, name: newFolderName.value } });
        await loadFileList(); // 重新加载文件列表
    }
}

const handleDownload = async row => {
    if (row.kind === 'DIR') return;
    window.open(`${API_BASE_URL}/fileManager/download?path=${encodeURIComponent(currentDir.value)}&name=${encodeURIComponent(row.name)}`, '_blank');
};

const handleDestroy = async row => {
    try {
        notify.ask(`确定要删除 【${row.name}】 吗？`, async () => {
            await axios.post('/fileManager/destroy', { body: { path: currentDir.value, name: row.name } });
            await loadFileList(); // 重新加载文件列表
        });
    } catch (error) {
        console.error('删除文件失败', error);
        notify.error('删除失败', error);
    }
};

const handleZip = async row => {
    await axios.post('/fileManager/zip', { body: { path: currentDir.value, name: row.name } });
    await loadFileList(); // 重新加载文件列表
};
</script>
