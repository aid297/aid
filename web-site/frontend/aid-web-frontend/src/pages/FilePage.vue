<template>
	<div class="row">
		<div class="col">
			<q-page class="q-pa-md">
				<div class="row q-gutter-md">
					<div class="col">
						<q-card flat square bordered>
							<q-card-section>
								<div class="text-h4 text-deep-orange">文件管理</div>
								<div class="text-h6 text-brown">上传文件</div>
							</q-card-section>

							<q-card-section class="q-pt-none">
								<q-uploader :url="`${API_BASE_URL}/upload`" label="选择文件" multiple @uploaded="handleUploaded" class="max-wight" @failed="handleFailed" />
							</q-card-section>
						</q-card>
					</div>
				</div>
			</q-page>
		</div>
	</div>
</template>

<script setup>
import { API_BASE_URL } from 'src/utils/fetch';
import notify from 'src/utils/notify';

const handleUploaded = info => {
	console.log('文件上传成功', info);
	// 可以在这里处理上传成功后的逻辑，如提示用户、更新文件列表等
	// this.$q.notify({ type: 'positive', message: '上传成功！' })
	notify.ok('上传成功');
};
</script>
