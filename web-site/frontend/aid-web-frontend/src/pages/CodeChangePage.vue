<template>
    <q-page class="q-pa-md">
        <q-card flat square bordered class="q-mb-md">
            <q-card-section>
                <div class="text-h4 text-deep-orange">代码格式转换</div>
            </q-card-section>
            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col q-gutter-md">
                        <q-btn square label="压缩json" color="purple" @click="compressJSON" />
                        <q-btn square label="格式化json" color="purple" @click="uncompressJSON" />
                        <q-btn square label="正序json" color="purple" @click="sortJSON(false)" />
                        <q-btn square label="倒序json" color="purple" @click="sortJSON(true)" />
                    </div>
                    <div class="col q-gutter-md">
                        <q-btn square label="json → yaml" color="amber" @click="jsonToYAML" />
                        <q-btn square label="yaml → json" color="amber" @click="yamlToJSON" />
                        <q-btn square label="yaml → toml" color="amber" @click="yamlToTOML" />
                    </div>
                    <div class="col q-gutter-md">
                        <q-btn square label="json → toml" color="blue" @click="jsonToTOML"></q-btn>
                        <q-btn square label="toml → json" color="blue" @click="tomlToJSON"></q-btn>
                        <q-btn square label="toml → yaml" color="blue" @click="tomlToYAML" />
                    </div>
                </div>
                <div class="row q-gutter-md">
                    <div class="col"><span>&nbsp;</span></div>
                </div>
                <div class="row q-gutter-md">
                    <div class="col">
                        <div ref="containerJSON" style="height: 70vh; border: 1px solid #e0e0e0; margin-top: 12px"></div>
                    </div>
                    <div class="col">
                        <div ref="containerYAML" style="height: 70vh; border: 1px solid #e0e0e0; margin-top: 12px"></div>
                    </div>
                    <div class="col">
                        <div ref="containerTOML" style="height: 70vh; border: 1px solid #e0e0e0; margin-top: 12px"></div>
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-page>
</template>

<script setup>
import { dump, load } from "js-yaml";
import JSON5 from "json5";
import { Loading } from "quasar";
import notify from "src/utils/notify";
import stripJsonComments from "strip-json-comments";
import { onBeforeUnmount, onMounted, ref } from "vue";

const containerJSON = ref(null);
const containerYAML = ref(null);
const containerTOML = ref(null);
let editorJSON = null;
let editorYAML = null;
let editorTOML = null;
let monacoJSON = null;
let monacoYAML = null;
let monacoTOML = null;

const codeJSON = `{"hello": "aid"}`;
const codeYAML = "";
const codeTOML = "";
let currentTheme = "vs-dark";

const initEditorJSON = async () => {
    // 动态导入 editor API 与样式
    monacoJSON = await import("monaco-editor/esm/vs/editor/editor.api.js");
    await import("monaco-editor/min/vs/editor/editor.main.css");

    // 为 Monaco
    self.MonacoEnvironment = {
        getWorkerUrl(_, label) {
            let workerPath = "";
            switch (label) {
                case "json":
                    workerPath = new URL("monaco-editor/esm/vs/language/json/json.worker.js", import.meta.url).toString();
                    break;
                case "css":
                case "scss":
                case "less":
                    workerPath = new URL("monaco-editor/esm/vs/language/css/css.worker.js", import.meta.url).toString();
                    break;
                case "html":
                case "handlebars":
                case "razor":
                    workerPath = new URL("monaco-editor/esm/vs/language/html/html.worker.js", import.meta.url).toString();
                    break;
                case "typescript":
                case "javascript":
                    workerPath = new URL("monaco-editor/esm/vs/language/typescript/ts.worker.js", import.meta.url).toString();
                    break;
                default:
                    workerPath = new URL("monaco-editor/esm/vs/editor/editor.worker.js", import.meta.url).toString();
            }
            return URL.createObjectURL(new Blob([`importScripts("${workerPath}");`], { type: "text/json" }));
        },
    };

    editorJSON = monacoJSON.editor.create(containerJSON.value, {
        value: codeJSON,
        language: "json",
        theme: currentTheme,
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 13,
    });
};

const initEditorYAML = async () => {
    // 动态导入 editor API 与样式
    monacoYAML = await import("monaco-editor/esm/vs/editor/editor.api.js");
    await import("monaco-editor/min/vs/editor/editor.main.css");

    // 为 Monaco
    self.MonacoEnvironment = {
        getWorkerUrl(_, label) {
            let workerPath = "";
            switch (label) {
                case "json":
                    workerPath = new URL("monaco-editor/esm/vs/language/json/json.worker.js", import.meta.url).toString();
                    break;
                case "css":
                case "scss":
                case "less":
                    workerPath = new URL("monaco-editor/esm/vs/language/css/css.worker.js", import.meta.url).toString();
                    break;
                case "html":
                case "handlebars":
                case "razor":
                    workerPath = new URL("monaco-editor/esm/vs/language/html/html.worker.js", import.meta.url).toString();
                    break;
                case "typescript":
                case "javascript":
                    workerPath = new URL("monaco-editor/esm/vs/language/typescript/ts.worker.js", import.meta.url).toString();
                    break;
                default:
                    workerPath = new URL("monaco-editor/esm/vs/editor/editor.worker.js", import.meta.url).toString();
            }
            return URL.createObjectURL(new Blob([`importScripts("${workerPath}");`], { type: "text/yaml" }));
        },
    };

    editorYAML = monacoYAML.editor.create(containerYAML.value, {
        value: codeYAML,
        language: "yaml",
        theme: currentTheme,
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 13,
    });
};

const initEditorTOML = async () => {
    // 动态导入 editor API 与样式
    monacoTOML = await import("monaco-editor/esm/vs/editor/editor.api.js");
    await import("monaco-editor/min/vs/editor/editor.main.css");

    // 为 Monaco
    self.MonacoEnvironment = {
        getWorkerUrl(_, label) {
            let workerPath = "";
            switch (label) {
                case "json":
                    workerPath = new URL("monaco-editor/esm/vs/language/json/json.worker.js", import.meta.url).toString();
                    break;
                case "css":
                case "scss":
                case "less":
                    workerPath = new URL("monaco-editor/esm/vs/language/css/css.worker.js", import.meta.url).toString();
                    break;
                case "html":
                case "handlebars":
                case "razor":
                    workerPath = new URL("monaco-editor/esm/vs/language/html/html.worker.js", import.meta.url).toString();
                    break;
                case "typescript":
                case "javascript":
                    workerPath = new URL("monaco-editor/esm/vs/language/typescript/ts.worker.js", import.meta.url).toString();
                    break;
                default:
                    workerPath = new URL("monaco-editor/esm/vs/editor/editor.worker.js", import.meta.url).toString();
            }
            return URL.createObjectURL(new Blob([`importScripts("${workerPath}");`], { type: "text/toml" }));
        },
    };

    editorTOML = monacoTOML.editor.create(containerTOML.value, {
        value: codeTOML,
        language: "toml",
        theme: currentTheme,
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 13,
    });
};

onMounted(() => {
    initEditorJSON().catch(err => notify.error(`编辑器初始化失败(JSON)：${err}`));

    initEditorYAML().catch(err => notify.error(`编辑器初始化失败(YAML)：${err}`));

    initEditorTOML().catch(err => notify.error(`编辑器初始化失败(TOML)：${err}`));
});

onBeforeUnmount(() => {
    Loading.show({ message: "加载中…", spinner: "QSpinnerDots" });
    if (editorJSON) {
        editorJSON.dispose();
        editorJSON = null;
    }

    if (editorYAML) {
        editorYAML.dispose();
        editorYAML = null;
    }

    if (editorTOML) {
        editorTOML.dispose();
        editorTOML = null;
    }

    Loading.hide();
});

const parseLenientJSON = text => {
    if (typeof text !== "string") throw new Error("输入不是字符串");

    // 1) 严格 JSON
    try {
        return JSON.parse(text);
    } catch (e) {
        notify.error(`解析json错误(严格)：${e}`);
    }

    // 2) 去注释后再 parse（支持 JSONC 风格注释）
    try {
        return JSON.parse(stripJsonComments(text));
    } catch (e) {
        notify.error(`解析json错误(JSONC风格)：${e}`);
    }

    // 3) JSON5（支持单引号、末尾逗号、不带引号的键等）
    try {
        return JSON5.parse(text);
    } catch (e) {
        throw new Error("无法解析为 JSON/JSONC/JSON5：" + (e.message || e));
    }
};

const compressJSON = () => {
    try {
        const srcText = editorJSON ? editorJSON.getValue() : "";
        if (!srcText) {
            notify.error("JSON内容为空");
            return;
        }

        // 解析并压缩：先尝试严格 JSON，再降级提示错误
        let obj;
        try {
            obj = parseLenientJSON(srcText);
        } catch (e) {
            notify.error(e.message || `解析失败：不是合法 JSON`);
            return;
        }

        const minified = JSON.stringify(obj); // 最紧凑形式

        if (editorJSON) {
            const model = editorJSON.getModel();
            try {
                if (model && monacoYAML && monacoYAML.editor && monacoYAML.editor.setModelLanguage) monacoYAML.editor.setModelLanguage(model, "json");
            } catch (e) {
                console.warn("setModelLanguage failed:", e);
            }
            if (model) model.setValue(minified);
            else editorJSON.setValue(minified);
        }

        const before = srcText.length;
        const after = minified.length;
        const saved = before - after;
        const percent = before ? Math.round((saved / before) * 100) : 0;

        notify.ok(`压缩完成：${before} → ${after} 字符，节省 ${saved} (${percent}%)`);
    } catch (e) {
        notify.error(`压缩失败：${e.message || e}`);
    }
};

const uncompressJSON = () => {
    try {
        const srcText = editorJSON ? editorJSON.getValue() : "";
        if (!srcText) {
            notify.error("JSON内容为空");
            return;
        }

        let obj;
        try {
            obj = parseLenientJSON(srcText);
        } catch (e) {
            notify.error(e.message || `解析失败：不是合法 JSON`);
            return;
        }

        // 使用 4 个空格作为缩进（等同于 tab=4 空格）
        const pretty = JSON.stringify(obj, null, 4);

        if (editorJSON) {
            const model = editorJSON.getModel();
            try {
                if (model && monacoJSON && monacoJSON.editor && monacoJSON.editor.setModelLanguage) monacoJSON.editor.setModelLanguage(model, "json");
            } catch (e) {
                console.warn("setModelLanguage failed:", e);
            }
            if (model) model.setValue(pretty);
            else editorJSON.setValue(pretty);
        }

        notify.ok("格式化完成（4 空格缩进）");
    } catch (e) {
        notify.error(`格式化失败：${e.message || e}`);
    }
};

const sortJSON = (reverse = false) => {
    try {
        const srcText = editorJSON ? editorJSON.getValue() : "";
        if (!srcText) {
            notify.warning("JSON 编辑器内容为空");
            return;
        }

        let obj;
        try {
            obj = parseLenientJSON(srcText);
        } catch (e) {
            notify.error(`解析 JSON 失败：${e.message || e}`);
            return;
        }

        // 递归排序对象的所有键
        const sortObjectKeys = input => {
            if (Array.isArray(input)) {
                // 如果是数组,递归处理数组中的每个元素
                return input.map(item => sortObjectKeys(item));
            } else if (input !== null && typeof input === "object") {
                // 如果是对象,按键排序并递归处理值
                const keys = Object.keys(input);
                const sortedKeys = reverse
                    ? keys.sort((a, b) => b.localeCompare(a)) // 倒序
                    : keys.sort(); // 正序

                return sortedKeys.reduce((sorted, key) => {
                    sorted[key] = sortObjectKeys(input[key]);
                    return sorted;
                }, {});
            }
            // 基本类型直接返回
            return input;
        };

        const sortedObj = sortObjectKeys(obj);
        const sortedJSON = JSON.stringify(sortedObj, null, 4); // 使用 4 空格缩进

        if (editorJSON) {
            editorJSON.setValue(sortedJSON);
        }

        notify.ok("JSON 键已按字母顺序排序");
    } catch (e) {
        notify.error(`排序失败：${e.message || e}`);
    }
};

const jsonToYAML = () => {
    try {
        const srcText = editorJSON ? editorJSON.getValue() : "";
        if (!srcText) {
            notify.error("JSON内容为空");
            return;
        }

        const obj = JSON.parse(srcText);
        const yamlText = dump(obj, { noRefs: true, lineWidth: -1 });

        if (editorYAML) {
            const model = editorYAML.getModel();
            // 尝试设置语言为 yaml（若 monaco 已加载且支持）
            try {
                if (model && monacoYAML && monacoYAML.editor && monacoYAML.editor.setModelLanguage) monacoYAML.editor.setModelLanguage(model, "yaml");
            } catch (e) {
                console.warn("setModelLanguage failed:", e);
            }
            if (model) model.setValue(yamlText);
            else editorYAML.setValue(yamlText);
        }

        notify.ok("转换成功");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const jsonToTOML = async () => {
    try {
        const srcText = editorJSON ? editorJSON.getValue() : "";
        if (!srcText) {
            notify.error("JSON内容为空");
            return;
        }

        let obj;
        try {
            obj = JSON.parse(srcText);
        } catch (e) {
            notify.error(`解析失败：不是合法 JSON → ${e.message || e}`);
            return;
        }

        // 在浏览器中 shim 全局并按需加载 @iarna/toml，避免 'global is not defined'
        if (typeof global === "undefined") window.global = window;

        let tomlText;
        try {
            const mod = await import("@iarna/toml");
            // 兼容不同导出形式
            const stringify = mod.stringify || mod.default?.stringify || mod.default;
            tomlText = stringify(obj);
        } catch (e) {
            notify.error(`TOML 序列化失败：${e.message || e}`);
            return;
        }

        if (editorTOML) {
            const model = editorTOML.getModel();
            try {
                if (model && monacoTOML && monacoTOML.editor && monacoTOML.editor.setModelLanguage) monacoTOML.editor.setModelLanguage(model, "toml");
            } catch (e) {
                console.warn("setModelLanguage failed:", e);
            }

            if (model) model.setValue(tomlText);
            else editorTOML.setValue(tomlText);
        }

        notify.ok("转换成功");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const yamlToJSON = () => {
    try {
        const dstText = editorYAML ? editorYAML.getValue() : "";
        if (!dstText) {
            notify.error({ type: "negative", message: "yaml内容为空" });
            return;
        }

        const obj = load(dstText);
        const jsonText = JSON.stringify(obj, null, 4);

        if (editorJSON) {
            const model = editorJSON.getModel();
            try {
                if (model && monacoJSON && monacoJSON.editor && monacoJSON.editor.setModelLanguage) monacoJSON.editor.setModelLanguage(model, "json");
            } catch (e) {
                console.warn("setModelLanguage failed:", e);
            }
            if (model) model.setValue(jsonText);
            else editorJSON.setValue(jsonText);
        }

        notify.ok("转换成功");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const yamlToTOML = () => {
    if (editorYAML.getValue() === "") {
        notify.error("YAML内容为空");
        return;
    }
    if (editorJSON.getValue() === "") {
        yamlToJSON();
    }
    jsonToTOML();
};

const tomlToJSON = async () => {
    try {
        const tomlText = editorTOML ? editorTOML.getValue() : "";
        if (!tomlText) {
            notify.error("TOML内容为空");
            return;
        }

        // 在浏览器环境中为 @iarna/toml 做 global shim（如果之前未做）
        if (typeof global === "undefined") window.global = window;

        let obj;
        try {
            const mod = await import("@iarna/toml");
            const parse = mod.parse || mod.default?.parse || mod.default;
            obj = parse(tomlText);
        } catch (e) {
            notify.error(`TOML 解析失败：${e.message || e}`);
            return;
        }

        const jsonText = JSON.stringify(obj, null, 4);

        if (editorJSON) {
            const model = editorJSON.getModel();
            try {
                if (model && monacoJSON && monacoJSON.editor && monacoJSON.editor.setModelLanguage) monacoJSON.editor.setModelLanguage(model, "json");
            } catch (e) {
                console.warn("setModelLanguage failed:", e);
            }
            if (model) model.setValue(jsonText);
            else editorJSON.setValue(jsonText);
        }

        notify.ok("转换成功");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const tomlToYAML = () => {
    if (editorTOML.getValue() === "") {
        notify.error("TOML内容为空");
        return;
    }
    if (editorJSON.getValue() === "") tomlToJSON();
    jsonToYAML();
};
</script>

<style scoped></style>
