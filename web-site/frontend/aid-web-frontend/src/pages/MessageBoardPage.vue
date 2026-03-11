<template>
    <q-page class="q-pa-md">
        <div class="row q-gutter-md">
            <!-- 生成 QRCode -->
            <div class="col">
                <q-card flat square bordered class="my-card">
                    <q-card-section>
                        <div class="text-h4 text-deep-orange">留言板</div>
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-input v-model="newMessageBoard" type="textarea" autogrow outlined square label="留言内容"
                                    placeholder="在此输入留言内容" ref="txtMessageBoard" autofocus>
                                    <template v-slot:append>
                                        <q-icon name="send" color="orange" @click="storeMessageBoard" />
                                    </template>
                                </q-input>
                            </div>
                        </div>
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                        <q-card class="my-card text-white q-mb-md"
                            :class="colors[Math.floor(Math.random() * colors.length)]"
                            v-for="messageBoard in messageBoards" :key="messageBoard.id">
                            <q-card-section>
                                <div class="text-h6">{{ messageBoard.id }}</div>
                                <div class="text-subtitle2"></div>
                            </q-card-section>

                            <q-card-section class="message-content" style="white-space: pre-wrap;">
                                {{ messageBoard.content }}
                            </q-card-section>

                            <q-separator dark />

                            <q-card-actions>
                                <q-btn flat @click="destroyMessageBoard(messageBoard.id)">删除</q-btn>
                                <q-btn flat @click="copy(messageBoard.content)">复制</q-btn>
                            </q-card-actions>
                        </q-card>
                    </q-card-section>
                </q-card>
            </div>
        </div>
    </q-page>
</template>
<script lang="js" setup>

import { axios } from 'src/utils/fetch';
import notify from 'src/utils/notify';
import { onMounted, ref } from 'vue';

const messageBoards = ref([]);
const newMessageBoard = ref('');
const txtMessageBoard = ref(null);

onMounted(() => loadMessageBoards());

const colors = [
    'bg-primary',
    'bg-secondary',
    'bg-accent',
    'bg-dark',
    'bg-positive',
    'bg-negative',
    'bg-warning',
    'bg-info',
    'bg-orange',
    'bg-teal',
    'bg-cyan',
    'bg-indigo',
    'bg-pink',
    'bg-purple'
];

const loadMessageBoards = () => {
    axios.post('/messageBoard/list').then(res => { messageBoards.value = res.data.content.messageBoards.reverse(); });
};

const copy = (content = '') => {
    if (!content) {
        notify.error('没有内容可复制');
        return;
    }
    navigator.clipboard
        ?.writeText(content)
        .then(() => {
            notify.ok('已复制到剪贴板');
        })
        .catch(() => {
            // 回退方案
            const ta = document.createElement('textarea');
            ta.value = content;
            ta.style.position = 'fixed';
            ta.style.opacity = '0';
            document.body.appendChild(ta);
            ta.select();
            try {
                document.execCommand('copy');
                notify.ok('已复制到剪贴板');
            } catch {
                notify.error('复制失败');
            } finally {
                document.body.removeChild(ta);
            }
        });
}

/**
 * 保存新消息
 */
const storeMessageBoard = async () => {
    if (!newMessageBoard.value) {
        return;
    }

    await axios.post('/messageBoard/store', { body: { content: newMessageBoard.value } }).then(res => { notify.ok(res.data.msg) });
    loadMessageBoards();
    newMessageBoard.value = '';
    txtMessageBoard.value.focus();
};

/**
 * 删除消息
 * @param id {string} 消息 ID
 */
const destroyMessageBoard = async (id) => {
    await axios.post('/messageBoard/destroy', { body: { id } }).then(() => { notify.ok('删除成功') });
    loadMessageBoards();
}

</script>

<style scoped>
/* 也可以写在一个类里，更整洁 */
.message-content {
    white-space: pre-wrap;
    /* 关键属性：保留空格和换行，并自动换行 */
    word-wrap: break-word;
    /* 防止长单词溢出 */
}
</style>
