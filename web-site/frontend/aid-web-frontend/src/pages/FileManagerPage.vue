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
						<q-table flat bordered separator title="文件管理" :rows="rows" dark color="amber" :pagination="{ rowsPerPage: 1000 }">
							<template v-slot:header="props">
								<q-tr :props="props">
									<q-th align="left" key="name" name="name">名称</q-th>
									<q-th align="left" key="size" name="size">大小</q-th>
									<q-th align="left" key="type" name="type">类型</q-th>
									<q-th align="left" key="modified" name="modified">修改时间</q-th>
								</q-tr>
							</template>
							<template v-slot:body="props">
								<q-tr :props="props">
									<q-td key="name" :props="props">{{ props.row.name || '??' }}</q-td>
									<q-td key="size" :props="props">{{ props.row.size || '??' }}</q-td>
									<q-td key="type" :props="props">{{ props.row.kind || '??' }}</q-td>
									<q-td key="modified" :props="props">{{ props.row.modified || '??' }}</q-td>
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
import { API_BASE_URL, fetcher } from 'src/utils/fetch';
import notify from 'src/utils/notify';
import { onMounted, ref } from 'vue';

const rows = ref([]);
const currentDir = ref('/');

/**
 * 加载文件列表
 * @param dir 所需目录
 */
const loadFileList = () => {
	fetcher.post('/fileManager/list', JSON.stringify({ dir: currentDir })).then(resp => {
        console.log("resp",resp);
		rows.value = resp.data.dirs || [];
		resp.data.dirs.forEach(dir => rows.value.push(dir));
		resp.data.files.forEach(file => rows.value.push(file));
		console.log(rows.value);
	});
};

onMounted(loadFileList);

const uploadUrl = ref(`${API_BASE_URL}/fileManager/upload`); // 后端上传接口
const handleUploaded = info => {
	console.log('文件上传成功', info);
	// 可以在这里处理上传成功后的逻辑，如提示用户、更新文件列表等
	// this.$q.notify({ type: 'positive', message: '上传成功！' })
	notify.ok('上传成功');
};

onMounted;
</script>
