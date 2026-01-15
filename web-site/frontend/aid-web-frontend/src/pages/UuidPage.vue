<template>
    <q-page class="q-pa-md">
        <div class="row q-gutter-md">
            <!-- 生成 UUID -->
            <div class="col">
                <q-card flat square bordered>
                    <q-card-section>
                        <div class="text-h4 text-deep-orange">生成 UUID</div>
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                        <div class="row">
                            <div class="col">
                                <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
                                    <div class="row q-gutter-md">
                                        <div class="col">
                                            <q-input v-model="number" label="生成uuid的数量" hint="必填" lazy-rules
                                                type="number" max="100" min="1" outlined square autofocus
                                                ref="txtNumber" step="1" />
                                        </div>
                                        <div class="col">
                                            <q-select v-model="uuidVersion" :options="versions" outlined square
                                                label="uuid版本" />
                                        </div>
                                    </div>
                                    <div class="row q-gutter-md q-mt-sm">
                                        <div class="col">
                                            <q-toggle v-model="noSubsTractKey" label="是否去掉“-”" class="full-width" />
                                        </div>
                                        <div class="col">
                                            <q-toggle v-model="isUpper" label="是否大写" class="full-width" />
                                        </div>
                                        <div class="col">
                                            <q-btn square label="提交" type="submit" color="primary" class="full-width" />
                                        </div>
                                    </div>
                                </q-form>
                            </div>
                        </div>
                        <div class="row">
                            <div class="col">
                                <hr />
                            </div>
                        </div>
                        <div class="row q-mt-sm">
                            <div class="col">
                                <q-table flat bordered separator title="UUID 生成结果" :rows="rows" dark color="amber"
                                    :pagination="{ rowsPerPage: 0 }">
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
                    </q-card-section>
                </q-card>
            </div>
        </div>
    </q-page>
</template>

<script setup>
import { fetcher } from "src/utils/fetch";
import { onMounted, ref } from "vue";

const number = ref(10);
const noSubsTractKey = ref(false);
const isUpper = ref(false);
const uuidVersion = ref("v6");
const txtNumber = ref(null);

const rows = ref([]);
const versions = ref([]);

onMounted(() => {
    fetcher.post("/uuid/versions").then(res => {
        for (const [k, v] of Object.entries(res.data.versions)) versions.value.push({ label: k, value: v });
    });
});

const onSubmit = async () => {
    rows.value = (await fetcher
        .post("/uuid/generate", {
            body: JSON.stringify({
                number: parseInt(number.value),
                noSubsTractKey: noSubsTractKey.value,
                isUpper: isUpper.value,
                version: uuidVersion.value.value,
            }),
        })).data.uuids;
    // fetcher
    //     .post("/uuid/generate", {
    //         body: JSON.stringify({
    //             number: parseInt(number.value),
    //             noSubsTractKey: noSubsTractKey.value,
    //             isUpper: isUpper.value,
    //             version: uuidVersion.value.value,
    //         }),
    //     })
    //     .then(res => (rows.value = res.data.uuids));
    txtNumber.value?.focus();
};

const onReset = () => (number.value = 10);
</script>
